package ui

import (
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/analyzer"
	"github.com/agnivo988/Repo-lyzer/internal/github"
)

type AnalysisResult struct {
	Repo                 *github.Repo
	Commits              []github.Commit
	Contributors         []github.Contributor
	FileTree             []github.TreeEntry
	Languages            map[string]int
	HealthScore          int
	BusFactor            int
	BusRisk              string
	MaturityScore        int
	MaturityLevel        string
	Dependencies         *analyzer.DependencyAnalysis
	ContributorInsights  *analyzer.ContributorInsights
	Security             *analyzer.SecurityScanResult
	CodeQuality          *analyzer.CodeQualityMetrics
	License              *analyzer.LicenseAnalysis
}

// CachedAnalysisResult wraps AnalysisResult with cache metadata
type CachedAnalysisResult struct {
	Result   AnalysisResult
	IsCached bool
	CachedAt time.Time
}

// CompareResult holds analysis data for two repositories
type CompareResult struct {
	Repo1 AnalysisResult
	Repo2 AnalysisResult
}
