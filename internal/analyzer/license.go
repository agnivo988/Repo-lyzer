// Package analyzer provides license detection and compliance analysis.
// This file implements license detection from LICENSE files and dependency licenses.
package analyzer

import (
	"encoding/base64"
	"strings"

	"github.com/agnivo988/Repo-lyzer/internal/github"
)

// LicenseInfo represents detected license information
type LicenseInfo struct {
	Name       string `json:"name"`        // License name (e.g., "MIT", "Apache-2.0")
	SPDX       string `json:"spdx"`        // SPDX identifier
	Category   string `json:"category"`    // "permissive", "copyleft", "proprietary"
	Commercial bool   `json:"commercial"`  // Allows commercial use
	Modify     bool   `json:"modify"`      // Allows modification
	Distribute bool   `json:"distribute"`  // Allows distribution
	Patent     bool   `json:"patent"`      // Includes patent grant
	SourceFile string `json:"source_file"` // File where license was found
}

// LicenseAnalysis holds complete license analysis results
type LicenseAnalysis struct {
	MainLicense    *LicenseInfo   `json:"main_license"`    // Primary project license
	OtherLicenses  []LicenseInfo  `json:"other_licenses"`  // Other licenses found
	Compatibility  string         `json:"compatibility"`   // "compatible", "warning", "conflict"
	Warnings       []string       `json:"warnings"`        // Potential issues
	LicenseScore   int            `json:"license_score"`   // 0-100 score
}

// Known license patterns for detection
var licensePatterns = map[string]LicenseInfo{
	"mit": {
		Name:       "MIT License",
		SPDX:       "MIT",
		Category:   "permissive",
		Commercial: true,
		Modify:     true,
		Distribute: true,
		Patent:     false,
	},
	"apache": {
		Name:       "Apache License 2.0",
		SPDX:       "Apache-2.0",
		Category:   "permissive",
		Commercial: true,
		Modify:     true,
		Distribute: true,
		Patent:     true,
	},
	"gpl-3": {
		Name:       "GNU GPL v3",
		SPDX:       "GPL-3.0",
		Category:   "copyleft",
		Commercial: true,
		Modify:     true,
		Distribute: true,
		Patent:     true,
	},
	"gpl-2": {
		Name:       "GNU GPL v2",
		SPDX:       "GPL-2.0",
		Category:   "copyleft",
		Commercial: true,
		Modify:     true,
		Distribute: true,
		Patent:     false,
	},
	"lgpl": {
		Name:       "GNU LGPL",
		SPDX:       "LGPL-3.0",
		Category:   "copyleft",
		Commercial: true,
		Modify:     true,
		Distribute: true,
		Patent:     false,
	},
	"bsd-3": {
		Name:       "BSD 3-Clause",
		SPDX:       "BSD-3-Clause",
		Category:   "permissive",
		Commercial: true,
		Modify:     true,
		Distribute: true,
		Patent:     false,
	},
	"bsd-2": {
		Name:       "BSD 2-Clause",
		SPDX:       "BSD-2-Clause",
		Category:   "permissive",
		Commercial: true,
		Modify:     true,
		Distribute: true,
		Patent:     false,
	},
	"isc": {
		Name:       "ISC License",
		SPDX:       "ISC",
		Category:   "permissive",
		Commercial: true,
		Modify:     true,
		Distribute: true,
		Patent:     false,
	},
	"mpl": {
		Name:       "Mozilla Public License 2.0",
		SPDX:       "MPL-2.0",
		Category:   "copyleft",
		Commercial: true,
		Modify:     true,
		Distribute: true,
		Patent:     true,
	},
	"unlicense": {
		Name:       "The Unlicense",
		SPDX:       "Unlicense",
		Category:   "permissive",
		Commercial: true,
		Modify:     true,
		Distribute: true,
		Patent:     false,
	},
	"cc0": {
		Name:       "CC0 1.0 Universal",
		SPDX:       "CC0-1.0",
		Category:   "permissive",
		Commercial: true,
		Modify:     true,
		Distribute: true,
		Patent:     false,
	},
	"agpl": {
		Name:       "GNU AGPL v3",
		SPDX:       "AGPL-3.0",
		Category:   "copyleft",
		Commercial: true,
		Modify:     true,
		Distribute: true,
		Patent:     true,
	},
}

// AnalyzeLicense detects and analyzes licenses in a repository
func AnalyzeLicense(client *github.Client, owner, repo string, fileTree []github.TreeEntry) (*LicenseAnalysis, error) {
	analysis := &LicenseAnalysis{
		OtherLicenses: []LicenseInfo{},
		Warnings:      []string{},
		Compatibility: "compatible",
		LicenseScore:  100,
	}

	// Find license files
	licenseFiles := findLicenseFiles(fileTree)

	if len(licenseFiles) == 0 {
		analysis.Warnings = append(analysis.Warnings, "No LICENSE file found")
		analysis.LicenseScore = 50
		return analysis, nil
	}

	// Analyze each license file
	for i, path := range licenseFiles {
		content, err := client.GetFileContent(owner, repo, path)
		if err != nil {
			continue
		}

		decoded, err := base64.StdEncoding.DecodeString(content)
		if err != nil {
			continue
		}

		license := detectLicense(string(decoded), path)
		if license != nil {
			if i == 0 {
				analysis.MainLicense = license
			} else {
				analysis.OtherLicenses = append(analysis.OtherLicenses, *license)
			}
		}
	}

	// Check for compatibility issues
	if analysis.MainLicense != nil {
		checkCompatibility(analysis)
	} else {
		analysis.Warnings = append(analysis.Warnings, "Could not detect license type")
		analysis.LicenseScore = 60
	}

	return analysis, nil
}

// findLicenseFiles finds license files in the file tree
func findLicenseFiles(tree []github.TreeEntry) []string {
	var files []string
	licenseNames := []string{
		"LICENSE",
		"LICENSE.md",
		"LICENSE.txt",
		"LICENCE",
		"LICENCE.md",
		"LICENCE.txt",
		"COPYING",
		"COPYING.md",
		"COPYING.txt",
	}

	for _, entry := range tree {
		if entry.Type != "blob" {
			continue
		}

		// Get filename from path
		parts := strings.Split(entry.Path, "/")
		filename := strings.ToUpper(parts[len(parts)-1])

		for _, licenseName := range licenseNames {
			if filename == licenseName || strings.ToUpper(filename) == licenseName {
				files = append(files, entry.Path)
				break
			}
		}
	}

	return files
}

// detectLicense detects license type from content
func detectLicense(content, sourcePath string) *LicenseInfo {
	contentLower := strings.ToLower(content)

	// Check for specific license patterns
	if strings.Contains(contentLower, "mit license") ||
		strings.Contains(contentLower, "permission is hereby granted, free of charge") {
		license := licensePatterns["mit"]
		license.SourceFile = sourcePath
		return &license
	}

	if strings.Contains(contentLower, "apache license") &&
		strings.Contains(contentLower, "version 2.0") {
		license := licensePatterns["apache"]
		license.SourceFile = sourcePath
		return &license
	}

	if strings.Contains(contentLower, "gnu general public license") {
		if strings.Contains(contentLower, "version 3") {
			license := licensePatterns["gpl-3"]
			license.SourceFile = sourcePath
			return &license
		}
		if strings.Contains(contentLower, "version 2") {
			license := licensePatterns["gpl-2"]
			license.SourceFile = sourcePath
			return &license
		}
	}

	if strings.Contains(contentLower, "gnu lesser general public license") ||
		strings.Contains(contentLower, "gnu library general public license") {
		license := licensePatterns["lgpl"]
		license.SourceFile = sourcePath
		return &license
	}

	if strings.Contains(contentLower, "gnu affero general public license") {
		license := licensePatterns["agpl"]
		license.SourceFile = sourcePath
		return &license
	}

	if strings.Contains(contentLower, "bsd 3-clause") ||
		(strings.Contains(contentLower, "redistribution") &&
			strings.Contains(contentLower, "neither the name")) {
		license := licensePatterns["bsd-3"]
		license.SourceFile = sourcePath
		return &license
	}

	if strings.Contains(contentLower, "bsd 2-clause") ||
		(strings.Contains(contentLower, "redistribution") &&
			!strings.Contains(contentLower, "neither the name") &&
			strings.Contains(contentLower, "binary form")) {
		license := licensePatterns["bsd-2"]
		license.SourceFile = sourcePath
		return &license
	}

	if strings.Contains(contentLower, "isc license") {
		license := licensePatterns["isc"]
		license.SourceFile = sourcePath
		return &license
	}

	if strings.Contains(contentLower, "mozilla public license") {
		license := licensePatterns["mpl"]
		license.SourceFile = sourcePath
		return &license
	}

	if strings.Contains(contentLower, "unlicense") ||
		strings.Contains(contentLower, "this is free and unencumbered software") {
		license := licensePatterns["unlicense"]
		license.SourceFile = sourcePath
		return &license
	}

	if strings.Contains(contentLower, "cc0") ||
		strings.Contains(contentLower, "creative commons zero") {
		license := licensePatterns["cc0"]
		license.SourceFile = sourcePath
		return &license
	}

	return nil
}

// checkCompatibility checks for license compatibility issues
func checkCompatibility(analysis *LicenseAnalysis) {
	if analysis.MainLicense == nil {
		return
	}

	mainCategory := analysis.MainLicense.Category

	// Check if copyleft license is used
	if mainCategory == "copyleft" {
		analysis.Warnings = append(analysis.Warnings,
			"Copyleft license requires derivative works to use same license")
	}

	// Check for AGPL
	if analysis.MainLicense.SPDX == "AGPL-3.0" {
		analysis.Warnings = append(analysis.Warnings,
			"AGPL requires source disclosure for network use")
	}

	// Check for mixed licenses
	for _, other := range analysis.OtherLicenses {
		if other.Category != mainCategory {
			analysis.Compatibility = "warning"
			analysis.Warnings = append(analysis.Warnings,
				"Mixed license categories detected: "+analysis.MainLicense.Name+" and "+other.Name)
			analysis.LicenseScore -= 10
		}

		// GPL incompatibility
		if mainCategory == "copyleft" && other.Category == "permissive" {
			// This is usually fine
		} else if mainCategory == "permissive" && other.Category == "copyleft" {
			analysis.Compatibility = "conflict"
			analysis.Warnings = append(analysis.Warnings,
				"Potential conflict: permissive main license with copyleft dependency")
			analysis.LicenseScore -= 20
		}
	}

	if analysis.LicenseScore < 0 {
		analysis.LicenseScore = 0
	}
}

// GetLicenseEmoji returns emoji for license category
func GetLicenseEmoji(category string) string {
	switch category {
	case "permissive":
		return "🟢"
	case "copyleft":
		return "🟡"
	case "proprietary":
		return "🔴"
	default:
		return "⚪"
	}
}

// GetLicenseGrade returns letter grade based on score
func GetLicenseGrade(score int) string {
	switch {
	case score >= 90:
		return "A"
	case score >= 80:
		return "B"
	case score >= 70:
		return "C"
	case score >= 60:
		return "D"
	default:
		return "F"
	}
}
