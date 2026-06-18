package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/analyzer"
	"github.com/agnivo988/Repo-lyzer/internal/github"
	"github.com/olekukonko/tablewriter"
)

// CompareInput contains the data required to build a repository comparison report.
type CompareInput struct {
	Repo         *github.Repo
	Commits      []github.Commit
	Contributors []github.Contributor
	Languages    map[string]int
	HasReleases  bool
}

// CompareRepository contains the computed metrics for one repository.
type CompareRepository struct {
	FullName        string         `json:"full_name"`
	Description     string         `json:"description,omitempty"`
	Stars           int            `json:"stars"`
	Forks           int            `json:"forks"`
	OpenIssues      int            `json:"open_issues"`
	CommitsLastYear int            `json:"commit_count_1y"`
	Contributors    int            `json:"contributors"`
	HealthScore     int            `json:"health_score"`
	BusFactor       int            `json:"bus_factor"`
	BusRisk         string         `json:"bus_risk"`
	MaturityScore   int            `json:"maturity_score"`
	MaturityLevel   string         `json:"maturity_level"`
	PrimaryLanguage string         `json:"primary_language,omitempty"`
	Languages       map[string]int `json:"languages,omitempty"`
}

// CompareReport is the complete comparison payload used by every output format.
type CompareReport struct {
	GeneratedAt time.Time        `json:"generated_at"`
	Repo1       CompareRepository `json:"repo1"`
	Repo2       CompareRepository `json:"repo2"`
	Verdict     string           `json:"verdict"`
}

// IsIdentical reports whether the key terminal metrics are equivalent.
func (r CompareReport) IsIdentical() bool {
	return r.Repo1.Stars == r.Repo2.Stars &&
		r.Repo1.Forks == r.Repo2.Forks &&
		r.Repo1.CommitsLastYear == r.Repo2.CommitsLastYear &&
		r.Repo1.Contributors == r.Repo2.Contributors &&
		r.Repo1.BusFactor == r.Repo2.BusFactor &&
		r.Repo1.MaturityScore == r.Repo2.MaturityScore
}

// BuildCompareReport calculates the comparison metrics for two repositories.
func BuildCompareReport(repo1, repo2 CompareInput) CompareReport {
	r1 := buildCompareRepository(repo1)
	r2 := buildCompareRepository(repo2)

	return CompareReport{
		GeneratedAt: time.Now(),
		Repo1:       r1,
		Repo2:       r2,
		Verdict:     BuildCompareVerdict(r1, r2),
	}
}

// BuildCompareVerdict summarizes which repository appears more mature.
func BuildCompareVerdict(repo1, repo2 CompareRepository) string {
	if repo1.MaturityScore > repo2.MaturityScore {
		return fmt.Sprintf("➡️ %s appears more mature and stable.", repo1.FullName)
	}

	if repo2.MaturityScore > repo1.MaturityScore {
		return fmt.Sprintf("➡️ %s appears more mature and stable.", repo2.FullName)
	}

	return "➡️ Both repositories are similarly mature."
}

// RenderCompareTerminal renders the legacy terminal comparison output.
func RenderCompareTerminal(report CompareReport) string {
	var buf bytes.Buffer
	buf.WriteString("\n📊 Repository Comparison\n")

	if report.IsIdentical() {
		buf.WriteString("\n✅ No differences found between the two repositories.\n")
		buf.WriteString("Both repositories have identical metrics.\n")
		return buf.String()
	}

	table := tablewriter.NewWriter(&buf)
	table.Header("Metric", report.Repo1.FullName, report.Repo2.FullName)

	table.Append([]string{"⭐ Stars",
		fmt.Sprintf("%d", report.Repo1.Stars),
		fmt.Sprintf("%d", report.Repo2.Stars),
	})

	table.Append([]string{"🍴 Forks",
		fmt.Sprintf("%d", report.Repo1.Forks),
		fmt.Sprintf("%d", report.Repo2.Forks),
	})

	table.Append([]string{"📦 Commits (1y)",
		fmt.Sprintf("%d", report.Repo1.CommitsLastYear),
		fmt.Sprintf("%d", report.Repo2.CommitsLastYear),
	})

	table.Append([]string{"👥 Contributors",
		fmt.Sprintf("%d", report.Repo1.Contributors),
		fmt.Sprintf("%d", report.Repo2.Contributors),
	})

	table.Append([]string{"⚠️ Bus Factor",
		fmt.Sprintf("%d (%s)", report.Repo1.BusFactor, report.Repo1.BusRisk),
		fmt.Sprintf("%d (%s)", report.Repo2.BusFactor, report.Repo2.BusRisk),
	})

	table.Append([]string{"🏗️ Maturity",
		fmt.Sprintf("%s (%d)", report.Repo1.MaturityLevel, report.Repo1.MaturityScore),
		fmt.Sprintf("%s (%d)", report.Repo2.MaturityLevel, report.Repo2.MaturityScore),
	})

	table.Render()

	buf.WriteString("\n Verdict\n")
	buf.WriteString(report.Verdict)
	buf.WriteString("\n")

	return buf.String()
}

// RenderCompareJSON renders the full comparison report as indented JSON.
func RenderCompareJSON(report CompareReport) ([]byte, error) {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal comparison report to JSON: %w", err)
	}

	return append(data, '\n'), nil
}

// RenderCompareMarkdown renders the comparison report as markdown.
func RenderCompareMarkdown(report CompareReport) ([]byte, error) {
	var buf strings.Builder

	buf.WriteString(fmt.Sprintf("# Repository Comparison: %s vs %s\n\n", report.Repo1.FullName, report.Repo2.FullName))
	buf.WriteString(fmt.Sprintf("*Generated: %s*\n\n", report.GeneratedAt.Format(time.RFC3339)))

	buf.WriteString("## Summary\n\n")
	buf.WriteString(fmt.Sprintf("| Metric | %s | %s |\n", report.Repo1.FullName, report.Repo2.FullName))
	buf.WriteString("| --- | --- | --- |\n")
	buf.WriteString(fmt.Sprintf("| Stars | %d | %d |\n", report.Repo1.Stars, report.Repo2.Stars))
	buf.WriteString(fmt.Sprintf("| Forks | %d | %d |\n", report.Repo1.Forks, report.Repo2.Forks))
	buf.WriteString(fmt.Sprintf("| Open Issues | %d | %d |\n", report.Repo1.OpenIssues, report.Repo2.OpenIssues))
	buf.WriteString(fmt.Sprintf("| Commits (1y) | %d | %d |\n", report.Repo1.CommitsLastYear, report.Repo2.CommitsLastYear))
	buf.WriteString(fmt.Sprintf("| Contributors | %d | %d |\n", report.Repo1.Contributors, report.Repo2.Contributors))
	buf.WriteString(fmt.Sprintf("| Health Score | %d | %d |\n", report.Repo1.HealthScore, report.Repo2.HealthScore))
	buf.WriteString(fmt.Sprintf("| Bus Factor | %d (%s) | %d (%s) |\n", report.Repo1.BusFactor, report.Repo1.BusRisk, report.Repo2.BusFactor, report.Repo2.BusRisk))
	buf.WriteString(fmt.Sprintf("| Maturity | %s (%d) | %s (%d) |\n", report.Repo1.MaturityLevel, report.Repo1.MaturityScore, report.Repo2.MaturityLevel, report.Repo2.MaturityScore))
	buf.WriteString(fmt.Sprintf("| Primary Language | %s | %s |\n", displayOrUnknown(report.Repo1.PrimaryLanguage), displayOrUnknown(report.Repo2.PrimaryLanguage)))

	buf.WriteString("\n## Verdict\n\n")
	buf.WriteString(report.Verdict)
	buf.WriteString("\n")

	return []byte(buf.String()), nil
}

func buildCompareRepository(input CompareInput) CompareRepository {
	if input.Repo == nil {
		return CompareRepository{}
	}

	busFactor, busRisk := analyzer.BusFactor(input.Contributors)
	maturityScore, maturityLevel := analyzer.RepoMaturityScore(input.Repo, len(input.Commits), len(input.Contributors), input.HasReleases)
	healthScore := analyzer.CalculateHealth(input.Repo, input.Commits)

	primaryLanguage := input.Repo.Language
	if primaryLanguage == "" {
		primaryLanguage = topLanguageName(input.Languages)
	}

	return CompareRepository{
		FullName:        input.Repo.FullName,
		Description:     input.Repo.Description,
		Stars:           input.Repo.Stars,
		Forks:           input.Repo.Forks,
		OpenIssues:      input.Repo.OpenIssues,
		CommitsLastYear: len(input.Commits),
		Contributors:    len(input.Contributors),
		HealthScore:     healthScore,
		BusFactor:       busFactor,
		BusRisk:         busRisk,
		MaturityScore:   maturityScore,
		MaturityLevel:   maturityLevel,
		PrimaryLanguage: primaryLanguage,
		Languages:       copyLanguageMap(input.Languages),
	}
}

func topLanguageName(languages map[string]int) string {
	if len(languages) == 0 {
		return ""
	}

	type languageEntry struct {
		name string
		size int
	}

	entries := make([]languageEntry, 0, len(languages))
	for name, size := range languages {
		if size <= 0 {
			continue
		}
		entries = append(entries, languageEntry{name: name, size: size})
	}

	if len(entries) == 0 {
		return ""
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].size > entries[j].size
	})

	return entries[0].name
}

func copyLanguageMap(languages map[string]int) map[string]int {
	if len(languages) == 0 {
		return nil
	}

	copyMap := make(map[string]int, len(languages))
	for name, size := range languages {
		copyMap[name] = size
	}

	return copyMap
}

func displayOrUnknown(value string) string {
	if value == "" {
		return "Unknown"
	}

	return value
}
