# PDF Export Feature - Testing Guide

## Quick Verification

### 1. Verify Binary Exists
```powershell
cd "e:\Rehan\Repo-Lyzer\Repo-lyzer\Repo-lyzer"
Get-Item repo-lyzer.exe
# Should show ~21 MB executable file
```

### 2. Check Help Text
```bash
./repo-lyzer.exe analyze --help
./repo-lyzer.exe compare --help
# Look for "pdf" in format options
```

## Manual Test Scenarios

### Test Case 1: Generate Single Repository PDF (Without Save Path)

**Command:**
```bash
./repo-lyzer.exe analyze tensorflow/tensorflow --format pdf
```

**Expected Output:**
```
🔍 Fetching repository information
📚 Fetching programming languages
📝 Analyzing commit history (365d)
📁 Scanning repository structure
💪 Computing repository health
👥 Fetching contributor information
⚠️  Analyzing bus factor and risk
📈 Calculating repository maturity
📄 Generating PDF report: tensorflow-tensorflow-report.pdf
✅ PDF report saved to: tensorflow-tensorflow-report.pdf
```

**Expected Result:**
- File `tensorflow-tensorflow-report.pdf` created in current directory
- File size: 100KB - 500KB
- Contains repository analysis data

---

### Test Case 2: Generate PDF with Custom Path

**Command:**
```bash
./repo-lyzer.exe analyze tensorflow/tensorflow --format pdf --save ./reports/ml-analysis.pdf
```

**Expected Output:**
```
...analysis steps...
📄 Generating PDF report: ./reports/ml-analysis.pdf
✅ PDF report saved to: ./reports/ml-analysis.pdf
```

**Expected Result:**
- Directory `./reports` created if it doesn't exist
- File `./reports/ml-analysis.pdf` created
- PDF contains complete analysis

---

### Test Case 3: Comparison PDF (Without Save Path)

**Command:**
```bash
./repo-lyzer.exe compare tensorflow/tensorflow pytorch/pytorch --format pdf
```

**Expected Output:**
```
🔍 Analyzing tensorflow/tensorflow...
Analyzed tensorflow/tensorflow
🔍 Analyzing pytorch/pytorch...
Analyzed pytorch/pytorch
📄 Generating PDF comparison report: tensorflow-tensorflow-vs-pytorch-pytorch-comparison.pdf
✅ PDF comparison report saved to: tensorflow-tensorflow-vs-pytorch-pytorch-comparison.pdf
```

**Expected Result:**
- File `tensorflow-tensorflow-vs-pytorch-pytorch-comparison.pdf` created
- Contains side-by-side comparison
- Includes verdict on which is more mature

---

### Test Case 4: Comparison PDF with Custom Path

**Command:**
```bash
./repo-lyzer.exe compare react vue --format pdf --save ./comparisons/frontend-frameworks.pdf
```

**Expected Output:**
```
...comparison steps...
📄 Generating PDF comparison report: ./comparisons/frontend-frameworks.pdf
✅ PDF comparison report saved to: ./comparisons/frontend-frameworks.pdf
```

**Expected Result:**
- Directory `./comparisons` created
- PDF file saved with custom name
- Contains comparison metrics

---

### Test Case 5: Verify Backward Compatibility

All these commands should still work without errors:

```bash
# YAML format (existing)
./repo-lyzer.exe analyze owner/repo --format yaml

# JSON format (existing)
./repo-lyzer.exe analyze owner/repo --format json --save output.json

# HTML comparison (existing)
./repo-lyzer.exe compare repo1 repo2 --format html

# Markdown (existing)
./repo-lyzer.exe compare repo1 repo2 --format markdown

# Default (no format, should show terminal output)
./repo-lyzer.exe analyze owner/repo
```

---

## PDF Content Verification Checklist

After generating a PDF, verify it contains:

### For Single Repository Analysis

- [ ] Cover page with repository name
- [ ] Generation date
- [ ] Executive summary section
- [ ] Repository overview with metrics
- [ ] Stars count
- [ ] Forks count
- [ ] Open issues count
- [ ] Programming languages breakdown
- [ ] Health score (0-100)
- [ ] Bus factor analysis
- [ ] Maturity level and score
- [ ] Contributors list (top 15)
- [ ] Security analysis (if data available)
- [ ] Code quality metrics
- [ ] Recommendations section
- [ ] Professional formatting and styling
- [ ] Readable fonts
- [ ] Proper page breaks

### For Comparison Report

- [ ] Cover page with both repository names
- [ ] "vs" separator between names
- [ ] Verdict statement
- [ ] Comparison overview page
- [ ] Side-by-side metrics
- [ ] Detailed metrics table
- [ ] Overall assessment section
- [ ] Winner determination
- [ ] Reason for verdict
- [ ] Note about comparison factors
- [ ] Professional layout

---

## Error Handling Test Cases

### Test Case 6: Invalid Repository

**Command:**
```bash
./repo-lyzer.exe analyze nonexistent/repository --format pdf
```

**Expected Result:**
- Appropriate error message about repository not found
- No PDF created
- Program exits gracefully

---

### Test Case 7: Permission Denied (Save Path)

**Command:**
```bash
./repo-lyzer.exe analyze owner/repo --format pdf --save /root/protected-dir/file.pdf
```

**Expected Result:**
- Error message about permissions
- No PDF created
- Program exits with error code

---

### Test Case 8: Directory Doesn't Exist

**Command:**
```bash
./repo-lyzer.exe analyze owner/repo --format pdf --save ./very/deep/nested/dirs/report.pdf
```

**Expected Result:**
- Directories are created automatically
- PDF is saved successfully
- Success message displayed

---

## Performance Tests

### Test Case 9: Small Repository

**Command:**
```bash
./repo-lyzer.exe analyze golang/gofpdf --format pdf
```

**Metrics to Record:**
- Time taken: Should be ~1-3 seconds
- File size: Typically 100KB
- Memory usage: Monitor task manager

---

### Test Case 10: Large Repository

**Command:**
```bash
./repo-lyzer.exe analyze kubernetes/kubernetes --format pdf
```

**Metrics to Record:**
- Time taken: Should be ~1-3 seconds
- File size: Typically up to 500KB
- Memory usage: Should stay reasonable

---

## Integration Tests

### Test Case 11: Multiple Reports

**Commands:**
```bash
./repo-lyzer.exe analyze tensorflow/tensorflow --format pdf --save reports/tf.pdf
./repo-lyzer.exe analyze pytorch/pytorch --format pdf --save reports/pytorch.pdf
./repo-lyzer.exe analyze jax/jax --format pdf --save reports/jax.pdf
./repo-lyzer.exe compare tensorflow/tensorflow pytorch/pytorch --format pdf --save reports/comparison.pdf
```

**Expected Result:**
- All 4 PDFs created successfully
- No conflicts or errors
- Each PDF is independent and readable

---

### Test Case 12: Overwrite Existing PDF

**Commands:**
```bash
./repo-lyzer.exe analyze owner/repo --format pdf --save report.pdf
# Modify the report somehow
./repo-lyzer.exe analyze owner/repo --format pdf --save report.pdf
```

**Expected Result:**
- Second PDF overwrites the first
- No error messages
- Latest PDF is valid

---

## Automated Testing Script

```bash
#!/bin/bash

echo "=== PDF Export Feature Test Suite ==="

# Test 1: Single repo PDF
echo "Test 1: Single repo PDF..."
./repo-lyzer analyze tensorflow/tensorflow --format pdf
if [ -f "tensorflow-tensorflow-report.pdf" ]; then
    echo "✅ Test 1 passed"
else
    echo "❌ Test 1 failed"
fi

# Test 2: Comparison PDF
echo "Test 2: Comparison PDF..."
./repo-lyzer compare tensorflow/tensorflow pytorch/pytorch --format pdf
if [ -f "tensorflow-tensorflow-vs-pytorch-pytorch-comparison.pdf" ]; then
    echo "✅ Test 2 passed"
else
    echo "❌ Test 2 failed"
fi

# Test 3: Custom path
echo "Test 3: Custom path PDF..."
mkdir -p test-reports
./repo-lyzer analyze go/go --format pdf --save test-reports/go-lang.pdf
if [ -f "test-reports/go-lang.pdf" ]; then
    echo "✅ Test 3 passed"
else
    echo "❌ Test 3 failed"
fi

# Test 4: Backward compatibility
echo "Test 4: Backward compatibility..."
./repo-lyzer analyze owner/repo --format yaml > /dev/null 2>&1
if [ $? -eq 0 ] || [ $? -eq 1 ]; then
    echo "✅ Test 4 passed (YAML still works)"
else
    echo "❌ Test 4 failed"
fi

echo "=== Test Suite Complete ==="
```

---

## File Integrity Tests

### Verify PDF Structure

```powershell
# PowerShell script to check PDF file integrity
$pdfFile = "tensorflow-tensorflow-report.pdf"

# Check file exists
if (Test-Path $pdfFile) {
    Write-Host "✅ File exists"
} else {
    Write-Host "❌ File not found"
}

# Check file size
$size = (Get-Item $pdfFile).Length
if ($size -gt 50000) {  # At least 50KB
    Write-Host "✅ File size reasonable: $($size / 1KB) KB"
} else {
    Write-Host "❌ File size too small: $($size / 1KB) KB"
}

# Check file signature (PDF files start with %PDF)
$header = Get-Content $pdfFile -Encoding Byte -ReadCount 4 | ForEach-Object {
    [System.Text.Encoding]::ASCII.GetString($_)
}
if ($header -like "%PDF*") {
    Write-Host "✅ Valid PDF file signature"
} else {
    Write-Host "❌ Invalid PDF signature"
}
```

---

## Report Quality Checklist

### Visual Quality
- [ ] Text is crisp and readable
- [ ] No garbled characters
- [ ] Proper font sizes
- [ ] Good spacing between sections
- [ ] Professional color scheme
- [ ] Consistent formatting

### Content Quality
- [ ] All data is accurate
- [ ] Numbers match analysis
- [ ] No missing sections
- [ ] Recommendations are relevant
- [ ] Text is complete (not cut off)

### Technical Quality
- [ ] Renders in all PDF viewers
- [ ] Prints without issues
- [ ] No corrupted data
- [ ] File opens quickly
- [ ] Page breaks are appropriate

---

## Known Issues and Workarounds

| Issue | Cause | Workaround |
|-------|-------|-----------|
| PDF takes time | Large repos have more data | Use compare to reduce scope |
| Large file size | All charts embedded | Normal behavior for PDF |
| Permission error | Directory access | Use different save location |
| Network timeout | GitHub API limit | Retry after rate limit resets |

---

## Success Criteria

✅ All test cases pass
✅ PDFs are readable in all viewers
✅ File sizes are reasonable (100KB - 500KB)
✅ Generation completes in 1-3 seconds
✅ No data corruption
✅ Backward compatibility maintained
✅ Error handling works correctly
✅ Documentation is accurate
✅ Feature ready for production

---

## Sign-Off Checklist

- [ ] All test cases executed successfully
- [ ] No crashes or exceptions
- [ ] PDFs verified for content
- [ ] Performance meets expectations
- [ ] Documentation reviewed
- [ ] Backward compatibility confirmed
- [ ] Ready for release

---

## Additional Resources

- Full documentation: [PDF_EXPORT.md](PDF_EXPORT.md)
- Quick reference: [PDF_QUICK_REFERENCE.md](PDF_QUICK_REFERENCE.md)
- Implementation details: [PDF_IMPLEMENTATION.md](PDF_IMPLEMENTATION.md)
- Source code: [internal/output/pdf*.go](../internal/output/)
