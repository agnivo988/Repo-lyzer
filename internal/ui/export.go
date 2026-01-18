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

// ValidateExportFormat checks if the given format is supported and returns a descriptive error if not
func ValidateExportFormat(format string) error {
	supportedFormats := []string{"json", "markdown", "csv", "html", "pdf"}
	
	// Convert to lowercase for case-insensitive comparison
	formatLower := strings.ToLower(format)
	
	for _, supported := range supportedFormats {
		if formatLower == supported {
			return nil
		}
	}
	
	return fmt.Errorf("unsupported export format '%s'. Supported formats are: %s", 
		format, strings.Join(supportedFormats, ", "))
}

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
	Login     string `json:"login"`
	Commits   int    `json:"commits"`
	AvatarURL string `json:"avatar_url,omitempty"`
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
			Login:     data.Contributors[i].Login,
			Commits:   data.Contributors[i].Commits,
			AvatarURL: data.Contributors[i].AvatarURL,
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

func ExportCSV(data AnalysisResult, _ string) (string, error) {
	downloadsDir, err := getDownloadsDir()
	if err != nil {
		return "", err
	}

	filename := filepath.Join(downloadsDir, generateFilename(data.Repo.FullName, "csv"))

	file, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Write CSV header
	header := "Metric,Value\n"
	file.WriteString(header)

	// Repository info
	fmt.Fprintf(file, "Repository Name,%s\n", data.Repo.FullName)
	fmt.Fprintf(file, "Description,%s\n", strings.ReplaceAll(data.Repo.Description, "\n", " "))
	fmt.Fprintf(file, "Stars,%d\n", data.Repo.Stars)
	fmt.Fprintf(file, "Forks,%d\n", data.Repo.Forks)
	fmt.Fprintf(file, "Open Issues,%d\n", data.Repo.OpenIssues)
	fmt.Fprintf(file, "Created At,%s\n", data.Repo.CreatedAt.Format("2006-01-02"))
	fmt.Fprintf(file, "Last Push,%s\n", data.Repo.PushedAt.Format("2006-01-02"))
	fmt.Fprintf(file, "Default Branch,%s\n", data.Repo.DefaultBranch)
	fmt.Fprintf(file, "URL,%s\n", data.Repo.HTMLURL)

	// Metrics
	fmt.Fprintf(file, "Health Score,%d\n", data.HealthScore)
	fmt.Fprintf(file, "Bus Factor,%d\n", data.BusFactor)
	fmt.Fprintf(file, "Bus Risk,%s\n", data.BusRisk)
	fmt.Fprintf(file, "Maturity Score,%d\n", data.MaturityScore)
	fmt.Fprintf(file, "Maturity Level,%s\n", data.MaturityLevel)
	fmt.Fprintf(file, "Total Commits,%d\n", len(data.Commits))
	fmt.Fprintf(file, "Total Contributors,%d\n", len(data.Contributors))

	// Languages
	file.WriteString("\nLanguages\n")
	file.WriteString("Language,Lines of Code\n")
	total := 0
	for _, bytes := range data.Languages {
		total += bytes
	}
	for lang, bytes := range data.Languages {
		pct := float64(bytes) / float64(total) * 100
		fmt.Fprintf(file, "%s,%.1f%%\n", lang, pct)
	}

	// Top Contributors
	file.WriteString("\nTop Contributors\n")
	file.WriteString("Login,Commits\n")
	maxContribs := 10
	if len(data.Contributors) < maxContribs {
		maxContribs = len(data.Contributors)
	}
	for i := 0; i < maxContribs; i++ {
		c := data.Contributors[i]
		fmt.Fprintf(file, "%s,%d\n", c.Login, c.Commits)
	}

	_ = openFileManager(filename)

	return filename, nil
}

func ExportHTML(data AnalysisResult, _ string) (string, error) {
	downloadsDir, err := getDownloadsDir()
	if err != nil {
		return "", err
	}

	filename := filepath.Join(downloadsDir, generateFilename(data.Repo.FullName, "html"))

	file, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Analysis for %s</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .header { background: #f0f0f0; padding: 20px; border-radius: 5px; }
        .section { margin: 20px 0; }
        table { border-collapse: collapse; width: 100%%; }
        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
        th { background-color: #f2f2f2; }
        .metric { font-weight: bold; }
    </style>
</head>
<body>
    <div class="header">
        <h1>Analysis for %s</h1>
        <p><em>Exported: %s</em></p>
    </div>

    <div class="section">
        <h2>Repository Info</h2>
        <table>
            <tr><td class="metric">Stars</td><td>%d</td></tr>
            <tr><td class="metric">Forks</td><td>%d</td></tr>
            <tr><td class="metric">Open Issues</td><td>%d</td></tr>
            <tr><td class="metric">Created</td><td>%s</td></tr>
            <tr><td class="metric">URL</td><td><a href="%s">%s</a></td></tr>
        </table>
    </div>

    <div class="section">
        <h2>Metrics</h2>
        <table>
            <tr><td class="metric">Health Score</td><td>%d/100</td></tr>
            <tr><td class="metric">Bus Factor</td><td>%d (%s)</td></tr>
            <tr><td class="metric">Maturity</td><td>%s (%d)</td></tr>
            <tr><td class="metric">Commits (1 year)</td><td>%d</td></tr>
            <tr><td class="metric">Contributors</td><td>%d</td></tr>
        </table>
    </div>

    <div class="section">
        <h2>Languages</h2>
        <table>
            <tr><th>Language</th><th>Percentage</th></tr>`,
		data.Repo.FullName, data.Repo.FullName, time.Now().Format("2006-01-02 15:04"),
		data.Repo.Stars, data.Repo.Forks, data.Repo.OpenIssues,
		data.Repo.CreatedAt.Format("2006-01-02"), data.Repo.HTMLURL, data.Repo.HTMLURL,
		data.HealthScore, data.BusFactor, data.BusRisk, data.MaturityLevel, data.MaturityScore,
		len(data.Commits), len(data.Contributors))

	// Languages
	total := 0
	for _, bytes := range data.Languages {
		total += bytes
	}
	for lang, bytes := range data.Languages {
		pct := float64(bytes) / float64(total) * 100
		html += fmt.Sprintf("<tr><td>%s</td><td>%.1f%%</td></tr>", lang, pct)
	}
	html += `        </table>
    </div>

    <div class="section">
        <h2>Top Contributors</h2>
        <table>
            <tr><th>Contributor</th><th>Commits</th></tr>`

	// Top Contributors
	maxContribs := 10
	if len(data.Contributors) < maxContribs {
		maxContribs = len(data.Contributors)
	}
	for i := 0; i < maxContribs; i++ {
		c := data.Contributors[i]
		html += fmt.Sprintf("<tr><td>%s</td><td>%d</td></tr>", c.Login, c.Commits)
	}

	html += `        </table>
    </div>
</body>
</html>`

	_, err = file.WriteString(html)
	if err != nil {
		return "", err
	}

	_ = openFileManager(filename)

	return filename, nil
}

// ExportAnalysis exports analysis data in the specified format with validation
func ExportAnalysis(data AnalysisResult, format string) (string, error) {
	// Validate the export format
	if err := ValidateExportFormat(format); err != nil {
		return "", err
	}
	
	// Convert format to lowercase for consistency
	format = strings.ToLower(format)
	
	// Call the appropriate export function based on format
	switch format {
	case "json":
		return ExportJSON(data, "")
	case "markdown":
		return ExportMarkdown(data, "")
	case "csv":
		return ExportCSV(data, "")
	case "html":
		return ExportHTML(data, "")
	case "pdf":
		return ExportPDF(data, "")
	default:
		// This should not happen due to validation, but just in case
		return "", fmt.Errorf("unexpected format after validation: %s", format)
	}
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
			Login:     data.Contributors[i].Login,
			Commits:   data.Contributors[i].Commits,
			AvatarURL: data.Contributors[i].AvatarURL,
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