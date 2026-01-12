package ui

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time" 
	"github.com/jung-kurt/gofpdf"
)

// ExportData is the structure for JSON export with additional metadata
type ExportData struct {
	ExportedAt      string              `json:"exported_at"`
	Repository      RepoExport          `json:"repository"`
	Metrics         MetricsExport       `json:"metrics"`
	Languages       map[string]int      `json:"languages"`
	TopContributors []ContributorExport `json:"top_contributors"`
	CommitCount     int                 `json:"commit_count_1y"`
}

type RepoExport struct {
	FullName      string `json:"full_name"`
	Description   string `json:"description"`
	Stars         int    `json:"stars"`
	Forks         int    `json:"forks"`
	OpenIssues    int    `json:"open_issues"`
	CreatedAt     string `json:"created_at"`
	LastPush      string `json:"last_push"`
	DefaultBranch string `json:"default_branch"`
	URL           string `json:"url"`
}

type MetricsExport struct {
	HealthScore   int    `json:"health_score"`
	BusFactor     int    `json:"bus_factor"`
	BusRisk       string `json:"bus_risk"`
	MaturityScore int    `json:"maturity_score"`
	MaturityLevel string `json:"maturity_level"`
}

type ContributorExport struct {
	Login   string `json:"login"`
	Commits int    `json:"commits"`
}

// getDownloadsDir returns the user's Downloads folder
func getDownloadsDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "Downloads"), nil
}

// openFileManager opens the file manager to show the exported file
func openFileManager(filePath string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		// Windows: use explorer with /select to highlight the file
		cmd = exec.Command("explorer", "/select,"+filePath)
	case "darwin":
		// macOS: use open with -R to reveal in Finder
		cmd = exec.Command("open", "-R", filePath)
	case "linux":
		// Linux: use xdg-open on the directory
		dir := filepath.Dir(filePath)
		cmd = exec.Command("xdg-open", dir)
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	return cmd.Start()
}

// generateFilename creates a filename with repo name and timestamp
func generateFilename(repoName, ext string) string {
	safeName := strings.ReplaceAll(repoName, "/", "_")
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	return fmt.Sprintf("%s_%s.%s", safeName, timestamp, ext)
}

func ExportJSON(data AnalysisResult, _ string) (string, error) {
	downloadsDir, err := getDownloadsDir()
	if err != nil {
		return "", err
	}

	filename := filepath.Join(downloadsDir, generateFilename(data.Repo.FullName, "json"))

	var topContribs []ContributorExport
	maxContribs := 10
	if len(data.Contributors) < maxContribs {
		maxContribs = len(data.Contributors)
	}
	for i := 0; i < maxContribs; i++ {
		topContribs = append(topContribs, ContributorExport{
			Login:   data.Contributors[i].Login,
			Commits: data.Contributors[i].Commits,
		})
	}

	export := ExportData{
		ExportedAt: time.Now().Format(time.RFC3339),
		Repository: RepoExport{
			FullName:      data.Repo.FullName,
			Description:   data.Repo.Description,
			Stars:         data.Repo.Stars,
			Forks:         data.Repo.Forks,
			OpenIssues:    data.Repo.OpenIssues,
			CreatedAt:     data.Repo.CreatedAt.Format("2006-01-02"),
			LastPush:      data.Repo.PushedAt.Format("2006-01-02"),
			DefaultBranch: data.Repo.DefaultBranch,
			URL:           data.Repo.HTMLURL,
		},
		Metrics: MetricsExport{
			HealthScore:   data.HealthScore,
			BusFactor:     data.BusFactor,
			BusRisk:       data.BusRisk,
			MaturityScore: data.MaturityScore,
			MaturityLevel: data.MaturityLevel,
		},
		Languages:       data.Languages,
		TopContributors: topContribs,
		CommitCount:     len(data.Commits),
	}

	file, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(export); err != nil {
		return "", err
	}

	_ = openFileManager(filename)

	return filename, nil
}

func ExportMarkdown(data AnalysisResult, _ string) (string, error) {
	downloadsDir, err := getDownloadsDir()
	if err != nil {
		return "", err
	}

	filename := filepath.Join(downloadsDir, generateFilename(data.Repo.FullName, "md"))

	file, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	md := fmt.Sprintf("# Analysis for %s\n\n", data.Repo.FullName)
	md += fmt.Sprintf("*Exported: %s*\n\n", time.Now().Format("2006-01-02 15:04"))

	md += "## Repository Info\n"
	md += fmt.Sprintf("- **Stars:** %d\n", data.Repo.Stars)
	md += fmt.Sprintf("- **Forks:** %d\n", data.Repo.Forks)
	md += fmt.Sprintf("- **Open Issues:** %d\n", data.Repo.OpenIssues)
	md += fmt.Sprintf("- **Created:** %s\n", data.Repo.CreatedAt.Format("2006-01-02"))
	md += fmt.Sprintf("- **URL:** %s\n\n", data.Repo.HTMLURL)

	md += "## Metrics\n"
	md += fmt.Sprintf("- **Health Score:** %d/100\n", data.HealthScore)
	md += fmt.Sprintf("- **Bus Factor:** %d (%s)\n", data.BusFactor, data.BusRisk)
	md += fmt.Sprintf("- **Maturity:** %s (%d)\n", data.MaturityLevel, data.MaturityScore)
	md += fmt.Sprintf("- **Commits (1 year):** %d\n", len(data.Commits))
	md += fmt.Sprintf("- **Contributors:** %d\n\n", len(data.Contributors))

	md += "## Languages\n"
	total := 0
	for _, bytes := range data.Languages {
		total += bytes
	}
	for lang, bytes := range data.Languages {
		pct := float64(bytes) / float64(total) * 100
		md += fmt.Sprintf("- %s: %.1f%%\n", lang, pct)
	}
	md += "\n"

	md += "## Top Contributors\n"
	maxContribs := 10
	if len(data.Contributors) < maxContribs {
		maxContribs = len(data.Contributors)
	}
	for i := 0; i < maxContribs; i++ {
		c := data.Contributors[i]
		md += fmt.Sprintf("%d. %s (%d commits)\n", i+1, c.Login, c.Commits)
	}

	_, err = file.WriteString(md)
	if err != nil {
		return "", err
	}

	_ = openFileManager(filename)

	return filename, nil
}

func ExportPDF(data AnalysisResult, _ string) (string, error) {
	downloadsDir, err := getDownloadsDir()
	if err != nil {
		return "", err
	}

	filename := filepath.Join(downloadsDir, generateFilename(data.Repo.FullName, "pdf"))

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Analysis for "+data.Repo.FullName)
	pdf.Ln(12)

	pdf.SetFont("Arial", "I", 10)
	pdf.Cell(0, 10, fmt.Sprintf("Exported: %s", time.Now().Format("2006-01-02 15:04")))
	pdf.Ln(15)

	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "Repository Info")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 11)
	pdf.Cell(0, 8, fmt.Sprintf("Stars: %d", data.Repo.Stars))
	pdf.Ln(6)
	pdf.Cell(0, 8, fmt.Sprintf("Forks: %d", data.Repo.Forks))
	pdf.Ln(6)
	pdf.Cell(0, 8, fmt.Sprintf("Open Issues: %d", data.Repo.OpenIssues))
	pdf.Ln(6)
	pdf.Cell(0, 8, fmt.Sprintf("Created: %s", data.Repo.CreatedAt.Format("2006-01-02")))
	pdf.Ln(6)
	pdf.Cell(0, 8, fmt.Sprintf("URL: %s", data.Repo.HTMLURL))
	pdf.Ln(15)

	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "Metrics")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 11)
	pdf.Cell(0, 8, fmt.Sprintf("Health Score: %d/100", data.HealthScore))
	pdf.Ln(6)
	pdf.Cell(0, 8, fmt.Sprintf("Bus Factor: %d (%s)", data.BusFactor, data.BusRisk))
	pdf.Ln(6)
	pdf.Cell(0, 8, fmt.Sprintf("Maturity: %s (%d)", data.MaturityLevel, data.MaturityScore))
	pdf.Ln(6)
	pdf.Cell(0, 8, fmt.Sprintf("Commits (1 year): %d", len(data.Commits)))
	pdf.Ln(6)
	pdf.Cell(0, 8, fmt.Sprintf("Contributors: %d", len(data.Contributors)))
	pdf.Ln(15)

	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "Languages")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 11)
	total := 0
	for _, bytes := range data.Languages {
		total += bytes
	}
		if total == 0 {  
		pdf.Cell(0, 8, "No language data available")  
		pdf.Ln(6)  
	} else {  
		for lang, bytes := range data.Languages {  
			pct := float64(bytes) / float64(total) * 100  
			pdf.Cell(0, 8, fmt.Sprintf("- %s: %.1f%%", lang, pct))  
			pdf.Ln(6)  
		}  
	}  
	pdf.Ln(9)

	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "Top Contributors")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 11)
	maxContribs := 10
	if len(data.Contributors) < maxContribs {
		maxContribs = len(data.Contributors)
	}
	for i := 0; i < maxContribs; i++ {
		c := data.Contributors[i]
		pdf.Cell(0, 8, fmt.Sprintf("%d. %s (%d commits)", i+1, c.Login, c.Commits))
		pdf.Ln(6)
	}

	err = pdf.OutputFileAndClose(filename)
	if err != nil {
		return "", err
	}

	_ = openFileManager(filename)

	return filename, nil
}

// CompareExportData is the structure for comparison JSON export
type CompareExportData struct {
	ExportedAt string     `json:"exported_at"`
	Repo1      ExportData `json:"repo1"`
	Repo2      ExportData `json:"repo2"`
	Verdict    string     `json:"verdict"`
}

func buildExportData(data AnalysisResult) ExportData {
	var topContribs []ContributorExport
	maxContribs := 10
	if len(data.Contributors) < maxContribs {
		maxContribs = len(data.Contributors)
	}
	for i := 0; i < maxContribs; i++ {
		topContribs = append(topContribs, ContributorExport{
			Login:   data.Contributors[i].Login,
			Commits: data.Contributors[i].Commits,
		})
	}

	return ExportData{
		ExportedAt: time.Now().Format(time.RFC3339),
		Repository: RepoExport{
			FullName:      data.Repo.FullName,
			Description:   data.Repo.Description,
			Stars:         data.Repo.Stars,
			Forks:         data.Repo.Forks,
			OpenIssues:    data.Repo.OpenIssues,
			CreatedAt:     data.Repo.CreatedAt.Format("2006-01-02"),
			LastPush:      data.Repo.PushedAt.Format("2006-01-02"),
			DefaultBranch: data.Repo.DefaultBranch,
			URL:           data.Repo.HTMLURL,
		},
		Metrics: MetricsExport{
			HealthScore:   data.HealthScore,
			BusFactor:     data.BusFactor,
			BusRisk:       data.BusRisk,
			MaturityScore: data.MaturityScore,
			MaturityLevel: data.MaturityLevel,
		},
		Languages:       data.Languages,
		TopContributors: topContribs,
		CommitCount:     len(data.Commits),
	}
}

func ExportCompareJSON(data CompareResult) (string, error) {
	downloadsDir, err := getDownloadsDir()
	if err != nil {
		return "", err
	}

	safeName1 := strings.ReplaceAll(data.Repo1.Repo.FullName, "/", "_")
	safeName2 := strings.ReplaceAll(data.Repo2.Repo.FullName, "/", "_")
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := filepath.Join(downloadsDir, fmt.Sprintf("compare_%s_vs_%s_%s.json", safeName1, safeName2, timestamp))

	var verdict string
	if data.Repo1.MaturityScore > data.Repo2.MaturityScore {
		verdict = fmt.Sprintf("%s appears more mature and stable", data.Repo1.Repo.FullName)
	} else if data.Repo2.MaturityScore > data.Repo1.MaturityScore {
		verdict = fmt.Sprintf("%s appears more mature and stable", data.Repo2.Repo.FullName)
	} else {
		verdict = "Both repositories are similarly mature"
	}

	export := CompareExportData{
		ExportedAt: time.Now().Format(time.RFC3339),
		Repo1:      buildExportData(data.Repo1),
		Repo2:      buildExportData(data.Repo2),
		Verdict:    verdict,
	}

	file, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(export); err != nil {
		return "", err
	}

	_ = openFileManager(filename)

	return filename, nil
}

func ExportCompareMarkdown(data CompareResult) (string, error) {
	downloadsDir, err := getDownloadsDir()
	if err != nil {
		return "", err
	}

	safeName1 := strings.ReplaceAll(data.Repo1.Repo.FullName, "/", "_")
	safeName2 := strings.ReplaceAll(data.Repo2.Repo.FullName, "/", "_")
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := filepath.Join(downloadsDir, fmt.Sprintf("compare_%s_vs_%s_%s.md", safeName1, safeName2, timestamp))

	file, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	r1 := data.Repo1
	r2 := data.Repo2

	md := fmt.Sprintf("# Comparison: %s vs %s\n\n", r1.Repo.FullName, r2.Repo.FullName)
	md += fmt.Sprintf("*Exported: %s*\n\n", time.Now().Format("2006-01-02 15:04"))

	md += "## Summary\n\n"
	md += "| Metric | " + r1.Repo.FullName + " | " + r2.Repo.FullName + " |\n"
	md += "|--------|--------|--------|\n"
	md += fmt.Sprintf("| Stars | %d | %d |\n", r1.Repo.Stars, r2.Repo.Stars)
	md += fmt.Sprintf("| Forks | %d | %d |\n", r1.Repo.Forks, r2.Repo.Forks)
	md += fmt.Sprintf("| Commits (1y) | %d | %d |\n", len(r1.Commits), len(r2.Commits))
	md += fmt.Sprintf("| Contributors | %d | %d |\n", len(r1.Contributors), len(r2.Contributors))
	md += fmt.Sprintf("| Health Score | %d | %d |\n", r1.HealthScore, r2.HealthScore)
	md += fmt.Sprintf("| Bus Factor | %d (%s) | %d (%s) |\n", r1.BusFactor, r1.BusRisk, r2.BusFactor, r2.BusRisk)
	md += fmt.Sprintf("| Maturity | %s (%d) | %s (%d) |\n", r1.MaturityLevel, r1.MaturityScore, r2.MaturityLevel, r2.MaturityScore)

	md += "\n## Verdict\n\n"
	if r1.MaturityScore > r2.MaturityScore {
		md += fmt.Sprintf("**%s** appears more mature and stable.\n", r1.Repo.FullName)
	} else if r2.MaturityScore > r1.MaturityScore {
		md += fmt.Sprintf("**%s** appears more mature and stable.\n", r2.Repo.FullName)
	} else {
		md += "Both repositories are similarly mature.\n"
	}

	_, err = file.WriteString(md)
	if err != nil {
		return "", err
	}

	_ = openFileManager(filename)

	return filename, nil
}