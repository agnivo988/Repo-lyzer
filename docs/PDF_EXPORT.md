# PDF Export Feature Documentation

## Overview

The Repo-lyzer PDF export feature enables users to generate professional, well-formatted PDF reports for both single repository analysis and repository comparisons. PDF reports are ideal for sharing with stakeholders, recruiters, and team members.

## Features

### Single Repository Analysis PDF

When using the `analyze` command with `--format pdf`, Repo-lyzer generates a comprehensive PDF report containing:

- **Cover Page**: Professional header with repository name, generation date, and prepared by information
- **Executive Summary**: Overall health score, grade, key findings, and risk assessment
- **Repository Overview**: Repository details, stars, forks, open issues, creation/update dates
- **Programming Languages**: Technology breakdown with percentage distribution
- **Commit Activity**: Visual chart showing commit patterns over time
- **Code Quality & Metrics**: Health score, maturity level, bus factor, contributors count
- **Security Analysis**: Vulnerability counts by severity, top security issues
- **Contributors**: Detailed contributor table with top 15 contributors
- **Recommendations**: Actionable recommendations based on repository metrics

### Repository Comparison PDF

When using the `compare` command with `--format pdf`, Repo-lyzer generates a side-by-side comparison report:

- **Cover Page**: Both repository names with "vs" separator, verdict summary
- **Comparison Overview**: Side-by-side metric comparison (stars, forks, commits, etc.)
- **Detailed Metrics**: Full table with all comparison metrics
- **Comparison Verdict**: Overall assessment and winner determination based on maturity

## Usage

### Generate PDF Report for Single Repository

```bash
# Generate PDF with auto-generated filename
repo-lyzer analyze owner/repository --format pdf

# Generate PDF with custom filename
repo-lyzer analyze owner/repository --format pdf --save reports/my-report.pdf
```

**Example Output:**
```text
🔍 Fetching repository information
📚 Fetching programming languages
📝 Analyzing commit history (365d)
📁 Scanning repository structure
💪 Computing repository health
👥 Fetching contributor information
⚠️  Analyzing bus factor and risk
📈 Calculating repository maturity
📄 Generating PDF report: owner-repository-report.pdf
✅ PDF report saved to: owner-repository-report.pdf
```

### Generate PDF Comparison Report

```bash
# Generate comparison PDF with auto-generated filename
repo-lyzer compare owner1/repo1 owner2/repo2 --format pdf

# Generate comparison PDF with custom filename
repo-lyzer compare owner1/repo1 owner2/repo2 --format pdf --save reports/comparison.pdf
```

**Example Output:**
```text
🔍 Analyzing owner1/repo1...
Analyzed owner1/repo1
🔍 Analyzing owner2/repo2...
Analyzed owner2/repo2
📄 Generating PDF comparison report: owner1-repo1-vs-owner2-repo2-comparison.pdf
✅ PDF comparison report saved to: owner1-repo1-vs-owner2-repo2-comparison.pdf
```

## File Naming

### Automatic Filename Generation

When `--save` is not specified:

**For single repository analysis:**
- Pattern: `{owner}-{repository}-report.pdf`
- Example: `facebook-react-report.pdf`

**For repository comparison:**
- Pattern: `{owner1}-{repo1}-vs-{owner2}-{repo2}-comparison.pdf`
- Example: `facebook-react-vs-angular-angular-comparison.pdf`

### Custom Filenames

Use the `--save` flag to specify a custom path:

```bash
repo-lyzer analyze angular/angular --format pdf --save ./reports/angular-analysis-2024.pdf
repo-lyzer compare react vue --format pdf --save ./comparisons/react-vs-vue.pdf
```

## Configuration

### Default Configuration

PDF reports use sensible defaults:
- **Cover Page**: Enabled
- **Table of Contents**: Enabled
- **Charts**: Enabled
- **Color Scheme**: Default professional blue/green
- **All Sections**: Enabled

### Advanced Configuration (Future)

While not implemented yet, the framework is designed to support file-based configuration for:
- Custom color schemes
- Company branding and logo
- Section visibility control
- Custom footer text

## Supported Data

Each PDF report includes:

### Repository Metrics
- Stars and forks count
- Open issues
- Commit activity
- Contributor count and breakdown
- Languages used
- Health score (0-100)
- Bus factor (1-n)
- Maturity score and level

### Analysis Results
- Code quality assessment
- Security vulnerabilities (if scanned)
- Risk levels and recommendations
- Quality dashboard metrics

### Charts and Visualizations
- Commit activity timeline
- Language distribution pie chart
- Contributor contribution bar chart
- Health score gauge

## Advantages Over Other Formats

| Feature | PDF | JSON | HTML | Markdown | CSV |
|---------|-----|------|------|----------|-----|
| **Professional Appearance** | ✅ | ❌ | ✅ | ⚠️ | ❌ |
| **Portable** | ✅ | ❌ | ⚠️ | ⚠️ | ❌ |
| **Print-Friendly** | ✅ | ❌ | ⚠️ | ⚠️ | ❌ |
| **Shareable** | ✅ | ⚠️ | ⚠️ | ⚠️ | ⚠️ |
| **Machine-Readable** | ❌ | ✅ | ⚠️ | ⚠️ | ✅ |
| **Charts & Visuals** | ✅ | ❌ | ✅ | ❌ | ❌ |
| **Mobile-Friendly** | ✅ | ❌ | ⚠️ | ✅ | ❌ |

## Implementation Details

### Architecture

The PDF export feature uses:

1. **gofpdf**: Go library for PDF generation (already in dependencies)
2. **EnhancedPDFGenerator**: Reusable PDF generation class with professional styling
3. **ComparisonPDFGenerator**: Specialized generator for comparison reports
4. **Chart Generation**: Existing chart functions for visualizations

### File Structure

```text
internal/output/
├── pdf.go                 # Main PDF export functions
├── pdf_enhanced.go        # EnhancedPDFGenerator implementation
├── pdf_config.go          # PDF configuration management
├── pdf_compare.go         # Comparison PDF generator
├── pdf_charts.go          # Chart generation functions
└── ...
```

### Key Functions

#### GenerateAnalyzePDF
```go
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
) error
```

#### GenerateComparePDF
```go
func GenerateComparePDF(
    report *CompareReport,
    savePath string,
    config *EnhancedPDFConfig,
) error
```

## Use Cases

### For Recruiters
- Professional repository analysis for candidate evaluation
- Shareable with hiring teams
- Print-friendly format for documentation

### For Teams
- Baseline repository health metrics
- Archive analysis reports
- Comparison before/after refactoring
- Technology stack documentation

### For Organizations
- Repository portfolio analysis
- Standardized reporting format
- Integration with documentation systems
- Executive summaries for stakeholders

## Examples

### Example 1: Analyze Popular Repository
```bash
repo-lyzer analyze kubernetes/kubernetes --format pdf --save reports/k8s-analysis.pdf
```

### Example 2: Compare Web Frameworks
```bash
repo-lyzer compare react vue --format pdf --save output.pdf
```

### Example 3: Monitor Project Health
```bash
# Generate weekly reports
repo-lyzer analyze myorg/myproject --format pdf --save reports/health-$(date +%Y-%m-%d).pdf
```

## Technical Specifications

- **Format**: PDF/A compatible (printable and archivable)
- **Page Size**: A4 (210mm × 297mm)
- **Orientation**: Portrait
- **Font**: Arial (standard, doesn't require embedding)
- **File Size**: Typically 100KB - 500KB depending on content
- **Generation Time**: 1-3 seconds on typical hardware

## Future Enhancements

Potential improvements for future versions:

1. **Configuration Files**: File-based PDF configuration
2. **Themes**: Multiple color schemes (corporate, dark, light)
3. **Branding**: Company logo and custom headers/footers
4. **Scheduling**: Automated PDF report generation
5. **Archiving**: Built-in report versioning and storage
6. **Multi-Repository Reports**: Analyze multiple repositories in one PDF
7. **Custom Sections**: User-defined report sections
8. **Watermarks**: Confidentiality and version watermarks
9. **Signatures**: Digital signatures for official reports
10. **Language Support**: Multi-language PDF generation

## Troubleshooting

### PDF File Not Created
- Check write permissions in the target directory
- Ensure disk space is available
- Verify the repository data was fetched successfully

### Large File Size
- PDF size increases with number of commits and contributors
- Charts and images are embedded in the PDF
- This is normal and expected

### Missing Charts
- Charts may not appear if data is insufficient
- Reports still generate successfully, just without visual elements
- This ensures robustness with edge cases

## Backward Compatibility

The PDF export feature is fully backward compatible:
- Existing commands continue to work unchanged
- PDF is an additional optional format
- No breaking changes to existing APIs
- Default behavior remains the same

## Related Commands

```bash
# View all available formats
repo-lyzer analyze --help

# Other export formats still available
repo-lyzer analyze owner/repo --format json --save output.json
repo-lyzer compare repo1 repo2 --format html --save comparison.html
repo-lyzer compare repo1 repo2 --format markdown
```

## Support

For issues or questions regarding PDF export:
1. Check the examples in this document
2. Review the project's GitHub issues
3. Contribute improvements via pull requests
4. Consult the main README for general usage
