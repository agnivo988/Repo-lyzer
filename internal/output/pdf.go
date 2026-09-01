package output

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/agnivo988/Repo-lyzer/internal/analyzer"
	"github.com/agnivo988/Repo-lyzer/internal/github"
)

// GenerateAnalyzePDF generates a PDF report for a single repository analysis
// Parameters:
//   - repo: Repository information
//   - commits: List of commits
//   - contributors: List of contributors
//   - languages: Map of programming languages
//   - healthScore: Repository health score (0-100)
//   - busFactor: Bus factor value
//   - busRisk: Bus factor risk level
//   - maturityScore: Repository maturity score
//   - maturityLevel: Repository maturity level
//   - savePath: File path where the PDF should be saved
//   - config: Optional custom PDF configuration (uses default if nil)
//
// Returns error if PDF generation fails
func GenerateAnalyzePDF(
	repo *github.Repo,
	commits []github.Commit,
	contributors []github.Contributor,
	languages map[string]int,
	healthScore int,
	busFactor int,
	busRisk string,
	maturityScore int,
	maturityLevel string,
	savePath string,
	config *EnhancedPDFConfig,
) error {
	// Use default config if not provided
	if config == nil {
		defaultConfig := DefaultPDFConfig()
		config = &defaultConfig
	}

	// Create PDF data
	data := &PDFData{
		Repo:          repo,
		Commits:       commits,
		Contributors:  contributors,
		Languages:     languages,
		HealthScore:   healthScore,
		BusFactor:     busFactor,
		BusRisk:       busRisk,
		MaturityScore: maturityScore,
		MaturityLevel: maturityLevel,
		Security:      nil, // Not available in basic analyze command
		QualityDashboard: nil, // Not available in basic analyze command
	}

	// Create PDF generator
	gen := NewEnhancedPDFGenerator(data, config)

	// Ensure output directory exists
	outputDir := filepath.Dir(savePath)
	if outputDir != "" && outputDir != "." {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	// Generate and save PDF
	if err := gen.Generate(savePath); err != nil {
		return fmt.Errorf("failed to generate PDF: %w", err)
	}

	return nil
}

// GenerateAnalyzePDFWithSecurity generates a PDF report with security data
func GenerateAnalyzePDFWithSecurity(
	repo *github.Repo,
	commits []github.Commit,
	contributors []github.Contributor,
	languages map[string]int,
	healthScore int,
	busFactor int,
	busRisk string,
	maturityScore int,
	maturityLevel string,
	security *analyzer.SecurityScanResult,
	qualityDashboard *analyzer.QualityDashboard,
	savePath string,
	config *EnhancedPDFConfig,
) error {
	// Use default config if not provided
	if config == nil {
		defaultConfig := DefaultPDFConfig()
		config = &defaultConfig
	}

	// Create PDF data
	data := &PDFData{
		Repo:             repo,
		Commits:          commits,
		Contributors:     contributors,
		Languages:        languages,
		HealthScore:      healthScore,
		BusFactor:        busFactor,
		BusRisk:          busRisk,
		MaturityScore:    maturityScore,
		MaturityLevel:    maturityLevel,
		Security:         security,
		QualityDashboard: qualityDashboard,
	}

	// Create PDF generator
	gen := NewEnhancedPDFGenerator(data, config)

	// Ensure output directory exists
	outputDir := filepath.Dir(savePath)
	if outputDir != "" && outputDir != "." {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	// Generate and save PDF
	if err := gen.Generate(savePath); err != nil {
		return fmt.Errorf("failed to generate PDF: %w", err)
	}

	return nil
}

// GenerateComparePDF generates a PDF report comparing two repositories
// This creates a side-by-side comparison report
func GenerateComparePDF(
	report *CompareReport,
	savePath string,
	config *EnhancedPDFConfig,
) error {
	// Use default config if not provided
	if config == nil {
		defaultConfig := DefaultPDFConfig()
		config = &defaultConfig
	}

	// Create output directory
	outputDir := filepath.Dir(savePath)
	if outputDir != "" && outputDir != "." {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	// For comparison, we create a simple PDF-like structure using gofpdf
	// Since we can't reuse the EnhancedPDFGenerator directly for comparison,
	// we'll create a specialized comparison PDF generator

	gen := NewComparisonPDFGenerator(report, config)
	if err := gen.Generate(savePath); err != nil {
		return fmt.Errorf("failed to generate comparison PDF: %w", err)
	}

	return nil
}
