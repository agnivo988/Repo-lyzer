# PDF Export Quick Reference

## One-Minute Guide

### Analyze Repository as PDF
```bash
repo-lyzer analyze owner/repo --format pdf
```

### Compare Repositories as PDF
```bash
repo-lyzer compare repo1 repo2 --format pdf
```

### Save to Custom Location
```bash
repo-lyzer analyze owner/repo --format pdf --save ./reports/analysis.pdf
repo-lyzer compare repo1 repo2 --format pdf --save ./reports/comparison.pdf
```

## What Gets Included

### Single Repository Report
✅ Repository metadata (stars, forks, issues, languages)
✅ Health score and grade (0-100 scale)
✅ Maturity assessment and bus factor
✅ Top 15 contributors with commit counts
✅ Commit activity visualization
✅ Language distribution chart
✅ Security vulnerabilities (if scanned)
✅ Code quality metrics
✅ Professional recommendations

### Comparison Report
✅ Side-by-side repository comparison
✅ Key metrics comparison table
✅ Overall assessment and winner
✅ Health score comparison
✅ Maturity and bus factor analysis

## Common Tasks

### Create Weekly Health Report
```bash
repo-lyzer analyze myorg/myrepo --format pdf --save reports/health-$(date +%Y-%m-%d).pdf
```

### Share with Recruiter
```bash
repo-lyzer analyze myorg/myrepo --format pdf --save projects/MyProject-Report.pdf
# Share the PDF file via email
```

### Archive Baseline Comparison
```bash
repo-lyzer compare original/repo fork/repo --format pdf --save archive/fork-comparison-baseline.pdf
```

### Evaluate Multiple Projects
```bash
for repo in repo1 repo2 repo3; do
  repo-lyzer analyze myorg/$repo --format pdf --save reports/$repo-analysis.pdf
done
```

## File Output

| Command | Default Filename |
|---------|-----------------|
| `repo-lyzer analyze owner/repo --format pdf` | `owner-repo-report.pdf` |
| `repo-lyzer compare repo1 repo2 --format pdf` | `repo1-vs-repo2-comparison.pdf` |

## Features

| Feature | Included |
|---------|----------|
| Cover page | ✅ |
| Executive summary | ✅ |
| Repository metrics | ✅ |
| Charts/visualizations | ✅ |
| Contributor analysis | ✅ |
| Security findings | ✅ |
| Recommendations | ✅ |
| Professional formatting | ✅ |
| Print-friendly | ✅ |

## Tips

- 💡 Use `--save` flag to control where PDFs are saved
- 💡 PDF filenames are auto-generated based on repository names if no path specified
- 💡 PDFs are typically 100KB - 500KB in size
- 💡 Generation takes 1-3 seconds on typical hardware
- 💡 PDFs include all analysis data without network dependency after generation
- 💡 Perfect for sharing with non-technical stakeholders
- 💡 Prints beautifully on standard paper

## Examples in Real Scenarios

### Scenario 1: Repository Portfolio Review
```bash
# Generate reports for a portfolio of repositories
repo-lyzer analyze tensorflow/tensorflow --format pdf --save portfolio/tensorflow.pdf
repo-lyzer analyze pytorch/pytorch --format pdf --save portfolio/pytorch.pdf
repo-lyzer analyze scikit-learn/scikit-learn --format pdf --save portfolio/scikit-learn.pdf
# Now you have professional reports to share
```

### Scenario 2: Framework Evaluation
```bash
# Compare popular web frameworks
repo-lyzer compare facebook/react vuejs/vue --format pdf --save evaluation/framework-comparison.pdf
```

### Scenario 3: Fork Evaluation
```bash
# Check if a fork is better maintained than original
repo-lyzer compare original/repo mycompany/repo-fork --format pdf --save audit/fork-audit.pdf
```

### Scenario 4: Automated Daily Reports
```bash
#!/bin/bash
# Generate daily health reports
REPO="myorg/myproject"
DATE=$(date +%Y-%m-%d)
repo-lyzer analyze $REPO --format pdf --save daily-reports/health-$DATE.pdf
```

## Integration

### With Documentation
```bash
# Generate PDF for documentation
repo-lyzer analyze owner/repo --format pdf --save docs/repository-analysis.pdf
```

### With Version Control
```bash
# Generate PDF and commit to reports branch
repo-lyzer analyze owner/repo --format pdf --save reports/latest-analysis.pdf
git add reports/latest-analysis.pdf
git commit -m "Update repository analysis report"
```

### With CI/CD
```yaml
# GitHub Actions example
- name: Generate Repository Analysis
  run: repo-lyzer analyze ${{ github.repository }} --format pdf --save analysis-${{ github.run_id }}.pdf

- name: Upload Report
  uses: actions/upload-artifact@v2
  with:
    name: repository-analysis
    path: analysis-*.pdf
```

## Comparison: Export Formats

| Format | Best For | Portable | Visual | Shareable |
|--------|----------|----------|--------|-----------|
| **PDF** | Reports, sharing, printing | ✅ | ✅ | ✅ |
| **HTML** | Web viewing, interactive | ⚠️ | ✅ | ⚠️ |
| **JSON** | Data processing, APIs | ✅ | ❌ | ⚠️ |
| **Markdown** | Documentation, git | ✅ | ⚠️ | ✅ |
| **YAML** | Config, structured data | ✅ | ❌ | ⚠️ |
| **CSV** | Spreadsheets, data analysis | ✅ | ❌ | ✅ |

## Troubleshooting Quick Fix

**Q: PDF not saving?**
A: Check directory permissions and disk space. Use `--save ./reports/file.pdf` instead of just `--save file.pdf`

**Q: PDF is taking too long?**
A: Normal for large repositories. Be patient, it should complete in 1-3 seconds.

**Q: Charts not appearing?**
A: Some edge cases may skip charts, but report still generates. Check network connection during analysis.

**Q: Want custom branding?**
A: Future feature! Currently using default professional styling.

## Next Steps

1. Generate your first PDF report
2. Share it with your team
3. Use it in documentation
4. Archive it for compliance
5. Enjoy professional repository analysis!

## Learn More

- Full documentation: [PDF_EXPORT.md](PDF_EXPORT.md)
- Command help: `repo-lyzer analyze --help`
- Examples: See the main README
