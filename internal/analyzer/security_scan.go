// Package analyzer provides security vulnerability scanning for dependencies.
package analyzer

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"
)

// Vulnerability represents a security vulnerability
type Vulnerability struct {
	ID          string   `json:"id"`
	Summary     string   `json:"summary"`
	Severity    string   `json:"severity"`
	Package     string   `json:"package"`
	Version     string   `json:"version"`
	FixedIn     string   `json:"fixed_in"`
	References  []string `json:"references"`
	PublishedAt string   `json:"published_at"`
}

// SecurityScanResult holds security scan results
type SecurityScanResult struct {
	Vulnerabilities []Vulnerability `json:"vulnerabilities"`
	TotalCount      int             `json:"total_count"`
	CriticalCount   int             `json:"critical_count"`
	HighCount       int             `json:"high_count"`
	MediumCount     int             `json:"medium_count"`
	LowCount        int             `json:"low_count"`
	ScannedPackages int             `json:"scanned_packages"`
	ScanTime        time.Time       `json:"scan_time"`
	SecurityScore   int             `json:"security_score"`
}

type osvQuery struct {
	Package struct {
		Name      string `json:"name"`
		Ecosystem string `json:"ecosystem"`
	} `json:"package"`
	Version string `json:"version,omitempty"`
}

type osvResponse struct {
	Vulns []osvVuln `json:"vulns"`
}

type osvVuln struct {
	ID       string `json:"id"`
	Summary  string `json:"summary"`
	Severity []struct {
		Type  string `json:"type"`
		Score string `json:"score"`
	} `json:"severity"`
	Affected []struct {
		Ranges []struct {
			Events []struct {
				Fixed string `json:"fixed,omitempty"`
			} `json:"events"`
		} `json:"ranges"`
	} `json:"affected"`
	References []struct {
		URL string `json:"url"`
	} `json:"references"`
	Published string `json:"published"`
}

// ScanDependencies scans dependencies for vulnerabilities
func ScanDependencies(deps *DependencyAnalysis) (*SecurityScanResult, error) {
	if deps == nil || len(deps.Files) == 0 {
		return &SecurityScanResult{
			Vulnerabilities: []Vulnerability{},
			ScanTime:        time.Now(),
			SecurityScore:   100,
		}, nil
	}

	result := &SecurityScanResult{
		Vulnerabilities: []Vulnerability{},
		ScanTime:        time.Now(),
	}

	client := &http.Client{Timeout: 10 * time.Second}

	for _, file := range deps.Files {
		ecosystem := mapEcosystem(file.FileType)
		if ecosystem == "" {
			continue
		}

		for _, dep := range file.Dependencies {
			result.ScannedPackages++
			vulns, _ := queryOSV(client, dep.Name, dep.Version, ecosystem)

			for _, v := range vulns {
				vuln := convertVuln(v, dep.Name, dep.Version)
				result.Vulnerabilities = append(result.Vulnerabilities, vuln)

				switch vuln.Severity {
				case "CRITICAL":
					result.CriticalCount++
				case "HIGH":
					result.HighCount++
				case "MEDIUM":
					result.MediumCount++
				case "LOW":
					result.LowCount++
				}
			}
		}
	}

	result.TotalCount = len(result.Vulnerabilities)
	result.SecurityScore = calcSecurityScore(result)
	return result, nil
}

func mapEcosystem(fileType string) string {
	m := map[string]string{"npm": "npm", "go": "Go", "python": "PyPI", "rust": "crates.io", "ruby": "RubyGems"}
	return m[fileType]
}

func queryOSV(client *http.Client, pkg, ver, eco string) ([]osvVuln, error) {
	query := osvQuery{}
	query.Package.Name = pkg
	query.Package.Ecosystem = eco
	ver = strings.TrimPrefix(strings.TrimPrefix(ver, "^"), "~")
	if ver != "" && ver != "*" {
		query.Version = ver
	}

	data, _ := json.Marshal(query)
	resp, err := client.Post("https://api.osv.dev/v1/query", "application/json", strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var r osvResponse
	json.NewDecoder(resp.Body).Decode(&r)
	return r.Vulns, nil
}

func convertVuln(o osvVuln, pkg, ver string) Vulnerability {
	v := Vulnerability{ID: o.ID, Summary: o.Summary, Package: pkg, Version: ver, PublishedAt: o.Published}
	v.Severity = getSeverity(o)
	for _, a := range o.Affected {
		for _, r := range a.Ranges {
			for _, e := range r.Events {
				if e.Fixed != "" {
					v.FixedIn = e.Fixed
				}
			}
		}
	}
	for _, ref := range o.References {
		if ref.URL != "" && len(v.References) < 3 {
			v.References = append(v.References, ref.URL)
		}
	}
	return v
}

func getSeverity(o osvVuln) string {
	for _, s := range o.Severity {
		if s.Type == "CVSS_V3" {
			score := parseCvssScore(s.Score)
			if score >= 9.0 {
				return "CRITICAL"
			} else if score >= 7.0 {
				return "HIGH"
			} else if score >= 4.0 {
				return "MEDIUM"
			} else if score > 0.0 {
				return "LOW"
			}
			return "NONE"
		}
	}
	return "MEDIUM"
}

func parseCvssScore(scoreStr string) float64 {
	var score float64
	if _, err := fmt.Sscanf(scoreStr, "%f", &score); err == nil {
		return score
	}
	if !strings.HasPrefix(scoreStr, "CVSS:") {
		return 0
	}
	parts := strings.Split(scoreStr, "/")
	metrics := make(map[string]string)
	for _, part := range parts {
		kv := strings.SplitN(part, ":", 2)
		if len(kv) == 2 {
			metrics[kv[0]] = kv[1]
		}
	}
	av := map[string]float64{"N": 0.85, "A": 0.62, "L": 0.55, "P": 0.20}[metrics["AV"]]
	ac := map[string]float64{"L": 0.77, "H": 0.44}[metrics["AC"]]
	ui := map[string]float64{"N": 0.85, "R": 0.62}[metrics["UI"]]
	c := map[string]float64{"H": 0.56, "L": 0.22}[metrics["C"]]
	i := map[string]float64{"H": 0.56, "L": 0.22}[metrics["I"]]
	a := map[string]float64{"H": 0.56, "L": 0.22}[metrics["A"]]
	var pr float64
	s := metrics["S"]
	if s == "C" {
		pr = map[string]float64{"N": 0.85, "L": 0.68, "H": 0.50}[metrics["PR"]]
	} else {
		pr = map[string]float64{"N": 0.85, "L": 0.62, "H": 0.27}[metrics["PR"]]
	}
	iss := 1 - (1-c)*(1-i)*(1-a)
	var impact float64
	if s == "C" {
		impact = 7.52*(iss-0.029) - 3.25*math.Pow(iss-0.02, 15)
	} else {
		impact = 6.42 * iss
	}
	exploitability := 8.22 * av * ac * pr * ui
	var base float64
	if s == "C" {
		base = 1.08 * (impact + exploitability)
	} else {
		base = impact + exploitability
	}
	if base > 10 {
		base = 10
	}
	if base < 0 {
		base = 0
	}
	return math.Ceil(base*10) / 10
}

func calcSecurityScore(r *SecurityScanResult) int {
	score := 100 - r.CriticalCount*25 - r.HighCount*15 - r.MediumCount*5 - r.LowCount*2
	if score < 0 {
		score = 0
	}
	return score
}

// GetSeverityEmoji returns emoji for severity
func GetSeverityEmoji(sev string) string {
	m := map[string]string{"CRITICAL": "🔴", "HIGH": "🟠", "MEDIUM": "🟡", "LOW": "🟢"}
	if e, ok := m[sev]; ok {
		return e
	}
	return "⚪"
}

// GetSecurityGrade returns letter grade
func GetSecurityGrade(score int) string {
	if score >= 90 {
		return "A"
	} else if score >= 80 {
		return "B"
	} else if score >= 70 {
		return "C"
	} else if score >= 60 {
		return "D"
	}
	return "F"
}
