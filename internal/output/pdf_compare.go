package output

import (
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf"
)

// ComparisonPDFGenerator generates PDF reports comparing two repositories
type ComparisonPDFGenerator struct {
	pdf    *gofpdf.Fpdf
	config *EnhancedPDFConfig
	report *CompareReport
}

// NewComparisonPDFGenerator creates a new comparison PDF generator
func NewComparisonPDFGenerator(report *CompareReport, config *EnhancedPDFConfig) *ComparisonPDFGenerator {
	pdf := gofpdf.New("P", "mm", "A4", "")

	return &ComparisonPDFGenerator{
		pdf:    pdf,
		config: config,
		report: report,
	}
}

// Generate creates the complete comparison PDF report
func (g *ComparisonPDFGenerator) Generate(filename string) error {
	// Add cover page
	g.addCompareCoverPage()

	// Add comparison overview
	g.addComparisonOverview()

	// Add detailed metrics
	g.addDetailedComparison()

	// Add verdict
	g.addComparisonVerdict()

	// Save the PDF
	return g.pdf.OutputFileAndClose(filename)
}

// addCompareCoverPage adds a cover page for the comparison report
func (g *ComparisonPDFGenerator) addCompareCoverPage() {
	g.pdf.AddPage()

	// Title
	g.pdf.SetFont("Arial", "B", 28)
	g.pdf.CellFormat(0, 20, "REPOSITORY COMPARISON", "", 0, "C", false, 0, "")
	g.pdf.Ln(30)

	// Repository names with vs
	g.pdf.SetFont("Arial", "B", 18)
	g.pdf.SetTextColor(30, 64, 175)

	// Left repo name
	g.pdf.CellFormat(80, 15, g.report.Repo1.FullName, "", 0, "C", false, 0, "")

	// VS
	g.pdf.SetFont("Arial", "B", 16)
	g.pdf.SetTextColor(128, 128, 128)
	g.pdf.CellFormat(35, 15, "vs", "", 0, "C", false, 0, "")

	// Right repo name
	g.pdf.SetFont("Arial", "B", 18)
	g.pdf.SetTextColor(30, 64, 175)
	g.pdf.CellFormat(0, 15, g.report.Repo2.FullName, "", 0, "C", false, 0, "")

	g.pdf.SetTextColor(0, 0, 0)
	g.pdf.Ln(30)

	// Verdict summary at bottom
	g.pdf.SetY(240)
	g.pdf.SetFont("Arial", "B", 14)
	g.pdf.SetTextColor(30, 64, 175)
	g.pdf.CellFormat(0, 10, fmt.Sprintf("Verdict: %s", g.report.Verdict), "", 0, "C", false, 0, "")
	g.pdf.SetTextColor(0, 0, 0)
	g.pdf.Ln(15)

	// Generation date
	g.pdf.SetFont("Arial", "I", 10)
	g.pdf.SetTextColor(128, 128, 128)
	g.pdf.CellFormat(0, 10, fmt.Sprintf("Generated: %s", time.Now().Format("January 2, 2006")), "", 0, "C", false, 0, "")
	g.pdf.SetTextColor(0, 0, 0)
}

// addComparisonOverview adds an overview of the comparison
func (g *ComparisonPDFGenerator) addComparisonOverview() {
	g.pdf.AddPage()

	g.addSectionHeader("Comparison Overview")

	// Side-by-side headers
	g.pdf.SetFont("Arial", "B", 12)
	g.pdf.CellFormat(90, 10, g.report.Repo1.FullName, "", 0, "C", false, 0, "")
	g.pdf.CellFormat(0, 10, g.report.Repo2.FullName, "", 0, "C", false, 0, "")
	g.pdf.Ln(10)

	// Stars
	g.pdf.SetFont("Arial", "", 11)
	g.pdf.CellFormat(90, 8, fmt.Sprintf("Stars: %d", g.report.Repo1.Stars), "", 0, "L", false, 0, "")
	g.pdf.CellFormat(0, 8, fmt.Sprintf("Stars: %d", g.report.Repo2.Stars), "", 0, "L", false, 0, "")
	g.pdf.Ln(8)

	// Forks
	g.pdf.CellFormat(90, 8, fmt.Sprintf("Forks: %d", g.report.Repo1.Forks), "", 0, "L", false, 0, "")
	g.pdf.CellFormat(0, 8, fmt.Sprintf("Forks: %d", g.report.Repo2.Forks), "", 0, "L", false, 0, "")
	g.pdf.Ln(8)

	// Open Issues
	g.pdf.CellFormat(90, 8, fmt.Sprintf("Open Issues: %d", g.report.Repo1.OpenIssues), "", 0, "L", false, 0, "")
	g.pdf.CellFormat(0, 8, fmt.Sprintf("Open Issues: %d", g.report.Repo2.OpenIssues), "", 0, "L", false, 0, "")
	g.pdf.Ln(8)

	// Commits
	g.pdf.CellFormat(90, 8, fmt.Sprintf("Commits (1y): %d", g.report.Repo1.CommitsLastYear), "", 0, "L", false, 0, "")
	g.pdf.CellFormat(0, 8, fmt.Sprintf("Commits (1y): %d", g.report.Repo2.CommitsLastYear), "", 0, "L", false, 0, "")
	g.pdf.Ln(8)

	// Contributors
	g.pdf.CellFormat(90, 8, fmt.Sprintf("Contributors: %d", g.report.Repo1.Contributors), "", 0, "L", false, 0, "")
	g.pdf.CellFormat(0, 8, fmt.Sprintf("Contributors: %d", g.report.Repo2.Contributors), "", 0, "L", false, 0, "")
	g.pdf.Ln(12)

	// Health scores
	g.pdf.SetFont("Arial", "B", 11)
	g.pdf.CellFormat(90, 8, fmt.Sprintf("Health Score: %d", g.report.Repo1.HealthScore), "", 0, "L", false, 0, "")
	g.pdf.CellFormat(0, 8, fmt.Sprintf("Health Score: %d", g.report.Repo2.HealthScore), "", 0, "L", false, 0, "")
	g.pdf.Ln(8)

	// Bus Factor
	g.pdf.SetFont("Arial", "", 11)
	g.pdf.CellFormat(90, 8, fmt.Sprintf("Bus Factor: %d (%s)", g.report.Repo1.BusFactor, g.report.Repo1.BusRisk), "", 0, "L", false, 0, "")
	g.pdf.CellFormat(0, 8, fmt.Sprintf("Bus Factor: %d (%s)", g.report.Repo2.BusFactor, g.report.Repo2.BusRisk), "", 0, "L", false, 0, "")
	g.pdf.Ln(8)

	// Maturity
	g.pdf.CellFormat(90, 8, fmt.Sprintf("Maturity: %s (%d)", g.report.Repo1.MaturityLevel, g.report.Repo1.MaturityScore), "", 0, "L", false, 0, "")
	g.pdf.CellFormat(0, 8, fmt.Sprintf("Maturity: %s (%d)", g.report.Repo2.MaturityLevel, g.report.Repo2.MaturityScore), "", 0, "L", false, 0, "")
	g.pdf.Ln(12)
}

// addDetailedComparison adds detailed metric comparisons
func (g *ComparisonPDFGenerator) addDetailedComparison() {
	g.pdf.AddPage()

	g.addSectionHeader("Detailed Metrics")

	// Create a simple comparison table
	g.pdf.SetFont("Arial", "B", 10)
	g.pdf.CellFormat(50, 8, "Metric", "1", 0, "C", false, 0, "")
	g.pdf.CellFormat(65, 8, g.report.Repo1.FullName, "1", 0, "C", false, 0, "")
	g.pdf.CellFormat(65, 8, g.report.Repo2.FullName, "1", 0, "C", false, 0, "")
	g.pdf.Ln(-1)

	// Metrics data
	g.pdf.SetFont("Arial", "", 9)

	metrics := []struct {
		label string
		val1  string
		val2  string
	}{
		{"Stars", fmt.Sprintf("%d", g.report.Repo1.Stars), fmt.Sprintf("%d", g.report.Repo2.Stars)},
		{"Forks", fmt.Sprintf("%d", g.report.Repo1.Forks), fmt.Sprintf("%d", g.report.Repo2.Forks)},
		{"Open Issues", fmt.Sprintf("%d", g.report.Repo1.OpenIssues), fmt.Sprintf("%d", g.report.Repo2.OpenIssues)},
		{"Commits (1Y)", fmt.Sprintf("%d", g.report.Repo1.CommitsLastYear), fmt.Sprintf("%d", g.report.Repo2.CommitsLastYear)},
		{"Contributors", fmt.Sprintf("%d", g.report.Repo1.Contributors), fmt.Sprintf("%d", g.report.Repo2.Contributors)},
		{"Health Score", fmt.Sprintf("%d/100", g.report.Repo1.HealthScore), fmt.Sprintf("%d/100", g.report.Repo2.HealthScore)},
		{"Bus Factor", fmt.Sprintf("%d (%s)", g.report.Repo1.BusFactor, g.report.Repo1.BusRisk), fmt.Sprintf("%d (%s)", g.report.Repo2.BusFactor, g.report.Repo2.BusRisk)},
		{"Maturity", fmt.Sprintf("%s (%d)", g.report.Repo1.MaturityLevel, g.report.Repo1.MaturityScore), fmt.Sprintf("%s (%d)", g.report.Repo2.MaturityLevel, g.report.Repo2.MaturityScore)},
	}

	for _, m := range metrics {
		g.pdf.CellFormat(50, 8, m.label, "1", 0, "L", false, 0, "")
		g.pdf.CellFormat(65, 8, m.val1, "1", 0, "C", false, 0, "")
		g.pdf.CellFormat(65, 8, m.val2, "1", 0, "C", false, 0, "")
		g.pdf.Ln(-1)
	}

	g.pdf.Ln(10)

	// Winner indicators
	g.pdf.SetFont("Arial", "B", 11)
	g.pdf.Cell(0, 10, "Overall Assessment:")
	g.pdf.Ln(8)

	g.pdf.SetFont("Arial", "", 10)
	winner := g.getWinner()
	g.pdf.MultiCell(0, 6, fmt.Sprintf(
		"Maturity comparison favors %s. %s",
		winner,
		g.report.Verdict,
	), "", "L", false)
}

// addComparisonVerdict adds the final verdict
func (g *ComparisonPDFGenerator) addComparisonVerdict() {
	g.pdf.AddPage()

	g.addSectionHeader("Comparison Verdict")

	winner := g.getWinner()

	g.pdf.SetFont("Arial", "B", 14)
	g.pdf.SetTextColor(30, 64, 175)
	g.pdf.Cell(0, 12, fmt.Sprintf("Winner: %s", winner))
	g.pdf.SetTextColor(0, 0, 0)
	g.pdf.Ln(15)

	g.pdf.SetFont("Arial", "B", 12)
	g.pdf.Cell(0, 10, fmt.Sprintf("Reason: %s", g.report.Verdict))
	g.pdf.Ln(15)

	g.pdf.SetFont("Arial", "", 11)
	g.pdf.MultiCell(0, 6, "This comparison is based on multiple factors including:\n"+
		"• Community engagement (stars, forks)\n"+
		"• Active development (commits, contributors)\n"+
		"• Repository health (health score)\n"+
		"• Risk assessment (bus factor)\n"+
		"• Project maturity and stability", "", "L", false)
	g.pdf.Ln(10)

	g.pdf.SetFont("Arial", "I", 10)
	g.pdf.SetTextColor(128, 128, 128)
	g.pdf.Cell(0, 10, "Note: This comparison provides general insights. Specific use cases may require different considerations.")
	g.pdf.SetTextColor(0, 0, 0)
}

// addSectionHeader adds a styled section header
func (g *ComparisonPDFGenerator) addSectionHeader(title string) {
	g.pdf.SetFont("Arial", "B", 16)
	g.pdf.SetTextColor(30, 64, 175)
	g.pdf.Cell(0, 12, title)
	g.pdf.SetTextColor(0, 0, 0)
	g.pdf.Ln(15)
}

// getWinner determines which repository is the winner based on maturity score
func (g *ComparisonPDFGenerator) getWinner() string {
	if g.report.Repo1.MaturityScore > g.report.Repo2.MaturityScore {
		return g.report.Repo1.FullName
	}
	if g.report.Repo2.MaturityScore > g.report.Repo1.MaturityScore {
		return g.report.Repo2.FullName
	}
	return "Both"
}
