package output

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"sort"
	"time"
)

//go:embed templates/compare.html
var compareTemplateFS embed.FS

var compareHTMLTemplate = template.Must(template.New("compare.html").ParseFS(compareTemplateFS, "templates/compare.html"))

type compareHTMLData struct {
	GeneratedAt string
	Verdict     string
	Repo1       compareHTMLRepo
	Repo2       compareHTMLRepo
	Metrics     []compareHTMLMetric
}

type compareHTMLRepo struct {
	FullName        string
	Description     string
	PrimaryLanguage string
	Stars           int
	Forks           int
	OpenIssues      int
	CommitsLastYear int
	Contributors    int
	HealthScore     int
	BusFactor       int
	BusRisk         string
	MaturityScore   int
	MaturityLevel   string
	HealthClass     string
	Languages       []compareLanguageShare
}

type compareHTMLMetric struct {
	Label      string
	Left       string
	Right      string
	LeftClass  string
	RightClass string
}

type compareLanguageShare struct {
	Name       string
	Percent    float64
	Width      int
	Percentage string
}

// RenderCompareHTML renders the report as a polished HTML document.
func RenderCompareHTML(report CompareReport) ([]byte, error) {
	data := compareHTMLData{
		GeneratedAt: report.GeneratedAt.Format(time.RFC3339),
		Verdict:     report.Verdict,
		Repo1:       buildCompareHTMLRepo(report.Repo1),
		Repo2:       buildCompareHTMLRepo(report.Repo2),
		Metrics:     buildCompareHTMLMetrics(report),
	}

	var buf bytes.Buffer
	if err := compareHTMLTemplate.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("failed to render comparison HTML: %w", err)
	}

	return buf.Bytes(), nil
}

func buildCompareHTMLRepo(repo CompareRepository) compareHTMLRepo {
	return compareHTMLRepo{
		FullName:        repo.FullName,
		Description:     repo.Description,
		PrimaryLanguage: displayOrUnknown(repo.PrimaryLanguage),
		Stars:           repo.Stars,
		Forks:           repo.Forks,
		OpenIssues:      repo.OpenIssues,
		CommitsLastYear: repo.CommitsLastYear,
		Contributors:    repo.Contributors,
		HealthScore:     repo.HealthScore,
		BusFactor:       repo.BusFactor,
		BusRisk:         repo.BusRisk,
		MaturityScore:   repo.MaturityScore,
		MaturityLevel:   repo.MaturityLevel,
		HealthClass:     scoreClass(repo.HealthScore),
		Languages:       buildCompareLanguageShares(repo.Languages, 5),
	}
}

func buildCompareHTMLMetrics(report CompareReport) []compareHTMLMetric {
	return []compareHTMLMetric{
		metricRow("Stars", report.Repo1.Stars, report.Repo2.Stars, true),
		metricRow("Forks", report.Repo1.Forks, report.Repo2.Forks, true),
		metricRow("Open Issues", report.Repo1.OpenIssues, report.Repo2.OpenIssues, false),
		metricRow("Commits (1y)", report.Repo1.CommitsLastYear, report.Repo2.CommitsLastYear, true),
		metricRow("Contributors", report.Repo1.Contributors, report.Repo2.Contributors, true),
		metricRow("Health Score", report.Repo1.HealthScore, report.Repo2.HealthScore, true),
		metricRow("Bus Factor", report.Repo1.BusFactor, report.Repo2.BusFactor, true),
		metricRow("Maturity Score", report.Repo1.MaturityScore, report.Repo2.MaturityScore, true),
		{
			Label:      "Primary Language",
			Left:       displayOrUnknown(report.Repo1.PrimaryLanguage),
			Right:      displayOrUnknown(report.Repo2.PrimaryLanguage),
			LeftClass:  "",
			RightClass: "",
		},
	}
}

func metricRow(label string, left, right int, higherIsBetter bool) compareHTMLMetric {
	leftClass, rightClass := "", ""
	if left > right {
		if higherIsBetter {
			leftClass = "metric-better"
			rightClass = "metric-worse"
		} else {
			leftClass = "metric-worse"
			rightClass = "metric-better"
		}
	} else if right > left {
		if higherIsBetter {
			leftClass = "metric-worse"
			rightClass = "metric-better"
		} else {
			leftClass = "metric-better"
			rightClass = "metric-worse"
		}
	}

	return compareHTMLMetric{
		Label:      label,
		Left:       fmt.Sprintf("%d", left),
		Right:      fmt.Sprintf("%d", right),
		LeftClass:  leftClass,
		RightClass: rightClass,
	}
}

func buildCompareLanguageShares(languages map[string]int, limit int) []compareLanguageShare {
	if len(languages) == 0 || limit <= 0 {
		return nil
	}

	type languageEntry struct {
		name string
		size int
	}

	entries := make([]languageEntry, 0, len(languages))
	total := 0
	for name, size := range languages {
		if size <= 0 {
			continue
		}
		total += size
		entries = append(entries, languageEntry{name: name, size: size})
	}

	if total == 0 || len(entries) == 0 {
		return nil
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].size > entries[j].size
	})

	if limit > len(entries) {
		limit = len(entries)
	}

	shares := make([]compareLanguageShare, 0, limit)
	for _, entry := range entries[:limit] {
		percent := float64(entry.size) / float64(total) * 100
		width := int(percent)
		if width < 4 {
			width = 4
		}
		if width > 100 {
			width = 100
		}

		shares = append(shares, compareLanguageShare{
			Name:       entry.name,
			Percent:    percent,
			Width:      width,
			Percentage: fmt.Sprintf("%.1f%%", percent),
		})
	}

	return shares
}

func scoreClass(score int) string {
	switch {
	case score >= 80:
		return "score-good"
	case score >= 60:
		return "score-warning"
	default:
		return "score-bad"
	}
}

// removed unused helpers: HasLanguages, joinClasses