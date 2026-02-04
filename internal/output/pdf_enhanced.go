package output

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/analyzer"
	"github.com/agnivo988/Repo-lyzer/internal/github"
	"github.com/jung-kurt/gofpdf"
)

// EnhancedPDFGenerator generates professional PDF reports with charts and branding
type EnhancedPDFGenerator struct {
	pdf      *gofpdf.Fpdf
	config   *EnhancedPDFConfig
	data     *PDFData
	pageNum  int
	tocItems []TOCItem
}

// PDFData contains all data needed for PDF generation
type PDFData struct {
	Repo             *github.Repo
	Commits          []github.Commit
	Contributors     []github.Contributor
	Languages        map[string]int
	HealthScore      int
	BusFactor        int
	BusRisk          string
	MaturityScore    int
	MaturityLevel    string
	Security         *analyzer.SecurityScanResult
	QualityDashboard *analyzer.QualityDashboard
}

// TOCItem represents an item in the table of contents
type TOCItem struct {
	Title string
	Page  int
}

// NewEnhancedPDFGenerator creates a new enhanced PDF generator
func NewEnhancedPDFGenerator(data *PDFData, config *EnhancedPDFConfig) *EnhancedPDFGenerator {
	pdf := gofpdf.New("P", "mm", "A4", "")

	return &EnhancedPDFGenerator{
		pdf:      pdf,
		config:   config,
		data:     data,
		pageNum:  0,
		tocItems: []TOCItem{},
	}
}

// Generate creates the complete PDF report
func (g *EnhancedPDFGenerator) Generate(filename string) error {
	// Add cover page if enabled
	if g.config.Report.ShowCoverPage {
		g.addCoverPage()
	}

	// Add content sections
	if g.config.Sections.ExecutiveSummary {
		g.addExecutiveSummary()
	}

	if g.config.Sections.RepositoryOverview {
		g.addRepositoryOverview()
	}

	if g.config.Sections.CodeQuality {
		g.addCodeQualitySection()
	}

	if g.config.Sections.Security {
		g.addSecuritySection()
	}

	if g.config.Sections.Contributors {
		g.addContributorsSection()
	}

	if g.config.Sections.Recommendations {
		g.addRecommendations()
	}

	// Save the PDF
	return g.pdf.OutputFileAndClose(filename)
}

// addCoverPage adds a professional cover page
func (g *EnhancedPDFGenerator) addCoverPage() {
	g.pdf.AddPage()

	// Add logo if provided
	if g.config.Company.Logo != "" && fileExists(g.config.Company.Logo) {
		// Logo at top center
		g.pdf.ImageOptions(g.config.Company.Logo, 70, 30, 70, 0, false, gofpdf.ImageOptions{}, 0, "")
	}

	// Title
	g.pdf.SetY(100)
	g.pdf.SetFont("Arial", "B", 28)
	g.pdf.CellFormat(0, 20, "REPOSITORY ANALYSIS", "", 0, "C", false, 0, "")
	g.pdf.Ln(20)

	// Repository name
	g.pdf.SetFont("Arial", "B", 20)
	g.pdf.SetTextColor(30, 64, 175) // Blue
	g.pdf.CellFormat(0, 15, g.data.Repo.FullName, "", 0, "C", false, 0, "")
	g.pdf.SetTextColor(0, 0, 0) // Reset to black
	g.pdf.Ln(30)

	// Company name if provided
	if g.config.Company.Name != "" {
		g.pdf.SetFont("Arial", "", 14)
		g.pdf.CellFormat(0, 10, "Prepared by:", "", 0, "C", false, 0, "")
		g.pdf.Ln(8)
		g.pdf.SetFont("Arial", "B", 16)
		g.pdf.CellFormat(0, 10, g.config.Company.Name, "", 0, "C", false, 0, "")
		g.pdf.Ln(20)
	}

	// Date
	g.pdf.SetY(240)
	g.pdf.SetFont("Arial", "I", 12)
	dateStr := fmt.Sprintf("Generated: %s", time.Now().Format("January 2, 2006"))
	g.pdf.CellFormat(0, 10, dateStr, "", 0, "C", false, 0, "")
	g.pdf.Ln(10)

	// Footer
	g.pdf.SetFont("Arial", "", 10)
	g.pdf.SetTextColor(128, 128, 128)
	g.pdf.CellFormat(0, 10, "Powered by Repo-lyzer", "", 0, "C", false, 0, "")
	g.pdf.SetTextColor(0, 0, 0)
}

// addExecutiveSummary adds executive summary page
func (g *EnhancedPDFGenerator) addExecutiveSummary() {
	g.pdf.AddPage()
	g.addTOCItem("Executive Summary", g.pdf.PageNo())

	g.addSectionHeader("Executive Summary")

	// Overall score and grade
	g.pdf.SetFont("Arial", "B", 14)
	g.pdf.Cell(0, 10, fmt.Sprintf("Health Score: %d/100", g.data.HealthScore))
	g.pdf.Ln(8)

	// Determine grade
	grade := "C"
	if g.data.HealthScore >= 90 {
		grade = "A+"
	} else if g.data.HealthScore >= 80 {
		grade = "A"
	} else if g.data.HealthScore >= 70 {
		grade = "B"
	} else if g.data.HealthScore >= 60 {
		grade = "C"
	} else {
		grade = "D"
	}

	g.pdf.Cell(0, 10, fmt.Sprintf("Overall Grade: %s", grade))
	g.pdf.Ln(15)

	// Key findings
	g.pdf.SetFont("Arial", "B", 12)
	g.pdf.Cell(0, 10, "Key Findings:")
	g.pdf.Ln(8)

	g.pdf.SetFont("Arial", "", 11)
	findings := []string{
		fmt.Sprintf("Repository has %d stars and %d forks", g.data.Repo.Stars, g.data.Repo.Forks),
		fmt.Sprintf("Active contributor community (%d contributors)", len(g.data.Contributors)),
		fmt.Sprintf("Maturity level: %s (score: %d)", g.data.MaturityLevel, g.data.MaturityScore),
		fmt.Sprintf("Bus factor: %d (%s)", g.data.BusFactor, g.data.BusRisk),
	}

	for _, finding := range findings {
		g.pdf.Cell(5, 7, "")
		g.pdf.MultiCell(0, 7, "• "+finding, "", "L", false)
	}

	g.pdf.Ln(10)

	// Risk assessment
	g.pdf.SetFont("Arial", "B", 12)
	riskLevel := "LOW"
	if g.data.HealthScore < 60 {
		riskLevel = "HIGH"
	} else if g.data.HealthScore < 75 {
		riskLevel = "MEDIUM"
	}
	g.pdf.Cell(0, 10, fmt.Sprintf("Risk Assessment: %s", riskLevel))
	g.pdf.Ln(15)

	// Add health score chart if charts enabled
	if g.config.Report.ShowCharts {
		g.addHealthScoreChart()
	}
}

// addRepositoryOverview adds repository overview section
func (g *EnhancedPDFGenerator) addRepositoryOverview() {
	g.pdf.AddPage()
	g.addTOCItem("Repository Overview", g.pdf.PageNo())

	g.addSectionHeader("Repository Overview")

	// Repository details
	g.pdf.SetFont("Arial", "B", 11)
	g.pdf.Cell(50, 7, "Repository:")
	g.pdf.SetFont("Arial", "", 11)
	g.pdf.Cell(0, 7, g.data.Repo.FullName)
	g.pdf.Ln(6)

	if g.data.Repo.Description != "" {
		g.pdf.SetFont("Arial", "B", 11)
		g.pdf.Cell(50, 7, "Description:")
		g.pdf.SetFont("Arial", "", 11)
		g.pdf.MultiCell(0, 7, g.data.Repo.Description, "", "L", false)
		g.pdf.Ln(2)
	}

	g.pdf.SetFont("Arial", "B", 11)
	g.pdf.Cell(50, 7, "Stars:")
	g.pdf.SetFont("Arial", "", 11)
	g.pdf.Cell(0, 7, fmt.Sprintf("%d", g.data.Repo.Stars))
	g.pdf.Ln(6)

	g.pdf.SetFont("Arial", "B", 11)
	g.pdf.Cell(50, 7, "Forks:")
	g.pdf.SetFont("Arial", "", 11)
	g.pdf.Cell(0, 7, fmt.Sprintf("%d", g.data.Repo.Forks))
	g.pdf.Ln(6)

	g.pdf.SetFont("Arial", "B", 11)
	g.pdf.Cell(50, 7, "Open Issues:")
	g.pdf.SetFont("Arial", "", 11)
	g.pdf.Cell(0, 7, fmt.Sprintf("%d", g.data.Repo.OpenIssues))
	g.pdf.Ln(6)

	g.pdf.SetFont("Arial", "B", 11)
	g.pdf.Cell(50, 7, "Created:")
	g.pdf.SetFont("Arial", "", 11)
	g.pdf.Cell(0, 7, g.data.Repo.CreatedAt.Format("January 2, 2006"))
	g.pdf.Ln(6)

	g.pdf.SetFont("Arial", "B", 11)
	g.pdf.Cell(50, 7, "Last Updated:")
	g.pdf.SetFont("Arial", "", 11)
	g.pdf.Cell(0, 7, g.data.Repo.PushedAt.Format("January 2, 2006"))
	g.pdf.Ln(15)

	// Add language distribution
	g.addLanguageDistribution()

	// Add commit activity chart if charts enabled
	if g.config.Report.ShowCharts && len(g.data.Commits) > 0 {
		g.pdf.Ln(10)
		g.addCommitActivityChart()
	}
}

// addCodeQualitySection adds code quality metrics
func (g *EnhancedPDFGenerator) addCodeQualitySection() {
	g.pdf.AddPage()
	g.addTOCItem("Code Quality & Metrics", g.pdf.PageNo())

	g.addSectionHeader("Code Quality & Metrics")

	// Health metrics
	g.pdf.SetFont("Arial", "B", 11)
	g.pdf.Cell(60, 7, "Health Score:")
	g.pdf.SetFont("Arial", "", 11)
	g.pdf.Cell(0, 7, fmt.Sprintf("%d/100", g.data.HealthScore))
	g.pdf.Ln(6)

	g.pdf.SetFont("Arial", "B", 11)
	g.pdf.Cell(60, 7, "Maturity Level:")
	g.pdf.SetFont("Arial", "", 11)
	g.pdf.Cell(0, 7, fmt.Sprintf("%s (Score: %d)", g.data.MaturityLevel, g.data.MaturityScore))
	g.pdf.Ln(6)

	g.pdf.SetFont("Arial", "B", 11)
	g.pdf.Cell(60, 7, "Bus Factor:")
	g.pdf.SetFont("Arial", "", 11)
	g.pdf.Cell(0, 7, fmt.Sprintf("%d (%s)", g.data.BusFactor, g.data.BusRisk))
	g.pdf.Ln(6)

	g.pdf.SetFont("Arial", "B", 11)
	g.pdf.Cell(60, 7, "Total Commits (1 year):")
	g.pdf.SetFont("Arial", "", 11)
	g.pdf.Cell(0, 7, fmt.Sprintf("%d", len(g.data.Commits)))
	g.pdf.Ln(6)

	g.pdf.SetFont("Arial", "B", 11)
	g.pdf.Cell(60, 7, "Contributors:")
	g.pdf.SetFont("Arial", "", 11)
	g.pdf.Cell(0, 7, fmt.Sprintf("%d", len(g.data.Contributors)))
	g.pdf.Ln(15)

	// Quality dashboard data if available
	if g.data.QualityDashboard != nil {
		g.addQualityDashboardData()
	}
}

// addSecuritySection adds security information
func (g *EnhancedPDFGenerator) addSecuritySection() {
	g.pdf.AddPage()
	g.addTOCItem("Security Analysis", g.pdf.PageNo())

	g.addSectionHeader("Security Analysis")

	if g.data.Security != nil {
		g.pdf.SetFont("Arial", "B", 11)
		g.pdf.Cell(60, 7, "Critical Vulnerabilities:")
		g.pdf.SetFont("Arial", "", 11)
		g.pdf.Cell(0, 7, fmt.Sprintf("%d", g.data.Security.CriticalCount))
		g.pdf.Ln(6)

		g.pdf.SetFont("Arial", "B", 11)
		g.pdf.Cell(60, 7, "High Vulnerabilities:")
		g.pdf.SetFont("Arial", "", 11)
		g.pdf.Cell(0, 7, fmt.Sprintf("%d", g.data.Security.HighCount))
		g.pdf.Ln(6)

		g.pdf.SetFont("Arial", "B", 11)
		g.pdf.Cell(60, 7, "Medium Vulnerabilities:")
		g.pdf.SetFont("Arial", "", 11)
		g.pdf.Cell(0, 7, fmt.Sprintf("%d", g.data.Security.MediumCount))
		g.pdf.Ln(6)

		g.pdf.SetFont("Arial", "B", 11)
		g.pdf.Cell(60, 7, "Low Vulnerabilities:")
		g.pdf.SetFont("Arial", "", 11)
		g.pdf.Cell(0, 7, fmt.Sprintf("%d", g.data.Security.LowCount))
		g.pdf.Ln(15)

		// Security recommendations
		if len(g.data.Security.Vulnerabilities) > 0 {
			g.pdf.SetFont("Arial", "B", 12)
			g.pdf.Cell(0, 10, "Top Security Issues:")
			g.pdf.Ln(8)

			g.pdf.SetFont("Arial", "", 10)
			maxIssues := 5
			if len(g.data.Security.Vulnerabilities) < maxIssues {
				maxIssues = len(g.data.Security.Vulnerabilities)
			}

			for i := 0; i < maxIssues; i++ {
				vuln := g.data.Security.Vulnerabilities[i]
				g.pdf.Cell(5, 6, "")
				text := fmt.Sprintf("%d. %s - Severity: %s", i+1, vuln.Package, vuln.Severity)
				g.pdf.MultiCell(0, 6, text, "", "L", false)
			}
		}
	} else {
		g.pdf.SetFont("Arial", "", 11)
		g.pdf.Cell(0, 7, "No security scan data available")
		g.pdf.Ln(6)
	}
}

// addContributorsSection adds contributor information
func (g *EnhancedPDFGenerator) addContributorsSection() {
	g.pdf.AddPage()
	g.addTOCItem("Contributors", g.pdf.PageNo())

	g.addSectionHeader("Contributors")

	g.pdf.SetFont("Arial", "", 11)
	g.pdf.Cell(0, 7, fmt.Sprintf("Total Contributors: %d", len(g.data.Contributors)))
	g.pdf.Ln(12)

	// Top contributors table
	g.pdf.SetFont("Arial", "B", 11)
	g.pdf.Cell(0, 10, "Top Contributors:")
	g.pdf.Ln(8)

	// Table header
	g.pdf.SetFont("Arial", "B", 10)
	g.pdf.CellFormat(15, 7, "Rank", "1", 0, "C", false, 0, "")
	g.pdf.CellFormat(100, 7, "Username", "1", 0, "L", false, 0, "")
	g.pdf.CellFormat(40, 7, "Commits", "1", 0, "C", false, 0, "")
	g.pdf.Ln(-1)

	// Table rows
	g.pdf.SetFont("Arial", "", 10)
	maxContribs := 15
	if len(g.data.Contributors) < maxContribs {
		maxContribs = len(g.data.Contributors)
	}

	for i := 0; i < maxContribs; i++ {
		c := g.data.Contributors[i]
		g.pdf.CellFormat(15, 7, fmt.Sprintf("%d", i+1), "1", 0, "C", false, 0, "")
		g.pdf.CellFormat(100, 7, c.Login, "1", 0, "L", false, 0, "")
		g.pdf.CellFormat(40, 7, fmt.Sprintf("%d", c.Commits), "1", 0, "C", false, 0, "")
		g.pdf.Ln(-1)
	}

	// Add contributor bar chart if charts enabled
	if g.config.Report.ShowCharts && len(g.data.Contributors) > 0 {
		g.pdf.Ln(10)
		g.addContributorChart()
	}
}

// addRecommendations adds recommendations section
func (g *EnhancedPDFGenerator) addRecommendations() {
	g.pdf.AddPage()
	g.addTOCItem("Recommendations", g.pdf.PageNo())

	g.addSectionHeader("Recommendations")

	recommendations := []string{}

	// Health-based recommendations
	if g.data.HealthScore < 70 {
		recommendations = append(recommendations, "Improve repository health by addressing open issues and maintaining regular commits")
	}

	// Bus factor recommendations
	if g.data.BusFactor < 3 {
		recommendations = append(recommendations, "Increase bus factor by encouraging more contributors to participate")
	}

	// Security recommendations
	if g.data.Security != nil && g.data.Security.CriticalCount > 0 {
		recommendations = append(recommendations, fmt.Sprintf("Address %d critical security vulnerabilities immediately", g.data.Security.CriticalCount))
	}

	// Quality dashboard recommendations
	if g.data.QualityDashboard != nil && len(g.data.QualityDashboard.Recommendations) > 0 {
		// Recommendations is []string, not a struct
		for _, rec := range g.data.QualityDashboard.Recommendations {
			recommendations = append(recommendations, rec)
		}
	}

	// Default recommendations if none generated
	if len(recommendations) == 0 {
		recommendations = append(recommendations, "Continue current maintenance practices")
		recommendations = append(recommendations, "Monitor metrics regularly for any degradation")
		recommendations = append(recommendations, "Keep dependencies up to date")
	}

	g.pdf.SetFont("Arial", "", 11)
	for i, rec := range recommendations {
		g.pdf.Cell(5, 7, "")
		g.pdf.MultiCell(0, 7, fmt.Sprintf("%d. %s", i+1, rec), "", "L", false)
		g.pdf.Ln(3)
	}
}

// Helper methods

func (g *EnhancedPDFGenerator) addSectionHeader(title string) {
	g.pdf.SetFont("Arial", "B", 16)
	g.pdf.SetTextColor(30, 64, 175)
	g.pdf.Cell(0, 12, title)
	g.pdf.SetTextColor(0, 0, 0)
	g.pdf.Ln(15)
}

func (g *EnhancedPDFGenerator) addTOCItem(title string, page int) {
	g.tocItems = append(g.tocItems, TOCItem{Title: title, Page: page})
}

func (g *EnhancedPDFGenerator) addLanguageDistribution() {
	g.pdf.SetFont("Arial", "B", 12)
	g.pdf.Cell(0, 10, "Programming Languages:")
	g.pdf.Ln(8)

	if len(g.data.Languages) == 0 {
		g.pdf.SetFont("Arial", "", 11)
		g.pdf.Cell(0, 7, "No language data available")
		g.pdf.Ln(6)
		return
	}

	total := 0
	for _, bytes := range g.data.Languages {
		total += bytes
	}

	g.pdf.SetFont("Arial", "", 10)
	for lang, bytes := range g.data.Languages {
		pct := float64(bytes) / float64(total) * 100
		if pct >= 1.0 { // Only show >= 1%
			g.pdf.Cell(5, 6, "")
			g.pdf.Cell(0, 6, fmt.Sprintf("• %s: %.1f%%", lang, pct))
			g.pdf.Ln(6)
		}
	}
}

func (g *EnhancedPDFGenerator) addHealthScoreChart() {
	imageData, err := GenerateHealthScoreGauge(g.data.HealthScore)
	if err != nil {
		return // Skip chart if error
	}

	// Save temp image
	tmpFile := filepath.Join(os.TempDir(), "health_chart.png")
	os.WriteFile(tmpFile, imageData, 0644)
	defer os.Remove(tmpFile)

	// Add to PDF
	g.pdf.ImageOptions(tmpFile, 50, 0, 110, 0, false, gofpdf.ImageOptions{ImageType: "PNG"}, 0, "")
	g.pdf.Ln(60)
}

func (g *EnhancedPDFGenerator) addCommitActivityChart() {
	imageData, err := GenerateCommitActivityChart(g.data.Commits)
	if err != nil {
		return
	}

	tmpFile := filepath.Join(os.TempDir(), "commit_chart.png")
	os.WriteFile(tmpFile, imageData, 0644)
	defer os.Remove(tmpFile)

	g.pdf.SetFont("Arial", "B", 12)
	g.pdf.Cell(0, 10, "Commit Activity Over Time:")
	g.pdf.Ln(8)

	g.pdf.ImageOptions(tmpFile, 10, 0, 190, 0, false, gofpdf.ImageOptions{ImageType: "PNG"}, 0, "")
	g.pdf.Ln(80)
}

func (g *EnhancedPDFGenerator) addContributorChart() {
	imageData, err := GenerateContributorBarChart(g.data.Contributors)
	if err != nil {
		return
	}

	tmpFile := filepath.Join(os.TempDir(), "contributor_chart.png")
	os.WriteFile(tmpFile, imageData, 0644)
	defer os.Remove(tmpFile)

	g.pdf.SetFont("Arial", "B", 12)
	g.pdf.Cell(0, 10, "Top Contributors Chart:")
	g.pdf.Ln(8)

	g.pdf.ImageOptions(tmpFile, 10, 0, 190, 0, false, gofpdf.ImageOptions{ImageType: "PNG"}, 0, "")
	g.pdf.Ln(80)
}

func (g *EnhancedPDFGenerator) addQualityDashboardData() {
	if g.data.QualityDashboard == nil {
		return
	}

	dash := g.data.QualityDashboard

	g.pdf.SetFont("Arial", "B", 12)
	g.pdf.Cell(0, 10, "Quality Dashboard:")
	g.pdf.Ln(8)

	g.pdf.SetFont("Arial", "B", 11)
	g.pdf.Cell(60, 7, "Overall Score:")
	g.pdf.SetFont("Arial", "", 11)
	g.pdf.Cell(0, 7, fmt.Sprintf("%d", dash.OverallScore))
	g.pdf.Ln(6)

	g.pdf.SetFont("Arial", "B", 11)
	g.pdf.Cell(60, 7, "Risk Level:")
	g.pdf.SetFont("Arial", "", 11)
	g.pdf.Cell(0, 7, dash.RiskLevel)
	g.pdf.Ln(6)

	g.pdf.SetFont("Arial", "B", 11)
	g.pdf.Cell(60, 7, "Grade:")
	g.pdf.SetFont("Arial", "", 11)
	g.pdf.Cell(0, 7, dash.QualityGrade)
	g.pdf.Ln(12)

	// Problem hotspots
	if len(dash.ProblemHotspots) > 0 {
		g.pdf.SetFont("Arial", "B", 11)
		g.pdf.Cell(0, 8, "Problem Hotspots:")
		g.pdf.Ln(6)

		g.pdf.SetFont("Arial", "", 10)
		for _, hotspot := range dash.ProblemHotspots {
			g.pdf.Cell(5, 6, "")
			text := fmt.Sprintf("• [%s] %s", hotspot.Severity, hotspot.Area)
			g.pdf.MultiCell(0, 6, text, "", "L", false)
		}
	}
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
