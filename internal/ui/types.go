package ui

import "github.com/agnivo988/Repo-lyzer/internal/github"

type AnalysisResult struct {
	Repo          *github.Repo
	Commits       []github.Commit
	Contributors  []github.Contributor
	FileTree      []github.TreeEntry
	Languages     map[string]int
	HealthScore   int
	BusFactor     int
	BusRisk       string
	MaturityScore int
	MaturityLevel string
}

// CompareResult holds analysis data for two repositories
type CompareResult struct {
	Repo1 AnalysisResult
	Repo2 AnalysisResult
}
