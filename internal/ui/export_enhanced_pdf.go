package ui

import (
	"path/filepath"

	"github.com/agnivo988/Repo-lyzer/internal/output"
)

// ExportEnhancedPDF exports analysis data as an enhanced PDF with charts and branding
func ExportEnhancedPDF(data AnalysisResult, filename string, config *output.EnhancedPDFConfig) (string, error) {
	if config == nil {
		config = output.LoadPDFConfigOrDefault()
	}

	// Convert AnalysisResult to PDFData
	pdfData := &output.PDFData{
		Repo:             data.Repo,
		Commits:          data.Commits,
		Contributors:     data.Contributors,
		Languages:        data.Languages,
		HealthScore:      data.HealthScore,
		BusFactor:        data.BusFactor,
		BusRisk:          data.BusRisk,
		MaturityScore:    data.MaturityScore,
		MaturityLevel:    data.MaturityLevel,
		Security:         data.Security,
		QualityDashboard: data.QualityDashboard,
	}

	// Determine filename
	downloadsDir, err := getDownloadsDir()
	if err != nil {
		return "", err
	}

	if filename == "" {
		filename = filepath.Join(downloadsDir, generateFilename(data.Repo.FullName, "pdf"))
	} else if !filepath.IsAbs(filename) {
		filename = filepath.Join(downloadsDir, filename)
	}

	// Generate PDF
	gen := output.NewEnhancedPDFGenerator(pdfData, config)
	err = gen.Generate(filename)
	if err != nil {
		return "", err
	}

	// Open file manager to show the exported file
	_ = openFileManager(filename)

	return filename, nil
}

// ExportEnhancedPDFWithDefaults exports using default configuration
func ExportEnhancedPDFWithDefaults(data AnalysisResult) (string, error) {
	return ExportEnhancedPDF(data, "", nil)
}
