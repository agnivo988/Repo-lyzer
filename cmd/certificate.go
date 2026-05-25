package cmd

import (
	"fmt"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/analyzer"
	"github.com/agnivo988/Repo-lyzer/internal/github"
	"github.com/agnivo988/Repo-lyzer/internal/output"
	"github.com/spf13/cobra"
)

// certificateCmd defines the "certificate" command for the CLI.
// It analyzes a GitHub repository and generates a comprehensive certificate
// showing the repository's overall score and potential uses.
// Usage example:
//
//	repo-lyzer certificate octocat/Hello-World
//
// This will perform full analysis and display a formatted certificate.
var certificateCmd = &cobra.Command{
	Use:   "certificate owner/repo",
	Short: "Generate a repository certificate with overall score and uses",
	Long: `Analyze a GitHub repository and generate a certificate that summarizes
the repository's health, maturity, activity, and overall quality score.
The certificate also suggests potential uses for the repository.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Validate the repository URL format
		owner, repo, err := validateRepoURL(args[0])
		if err != nil {
			return fmt.Errorf("invalid repository URL: %w", err)
		}

		// Record start time for analysis timing
		startTime := time.Now()

		fmt.Printf("🔍 Analyzing repository %s/%s for certificate generation...\n", owner, repo)

		// Initialize GitHub client
		client := github.NewClient()

		// Generate certificate data
		certificate, err := analyzer.GenerateCertificate(owner, repo, client)
		if err != nil {
			return fmt.Errorf("failed to generate certificate: %w", err)
		}

		// Generate and save PDF certificate
		pdfPath, err := output.ExportCertificatePDF(certificate)
		if err != nil {
			return fmt.Errorf("failed to export certificate PDF: %w", err)
		}

		// Display analysis time and PDF location
		duration := time.Since(startTime)
		fmt.Printf("\n⏱️  Certificate generated in %.2f seconds\n", duration.Seconds())
		fmt.Printf("📄 PDF saved to: %s\n", pdfPath)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(certificateCmd)
}
