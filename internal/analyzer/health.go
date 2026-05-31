package analyzer

import (
	"github.com/agnivo988/Repo-lyzer/internal/analyzer/health"
	"github.com/agnivo988/Repo-lyzer/internal/github"
)

// CalculateHealth is the legacy wrapper for the new modular HealthAnalyzer
func CalculateHealth(repo *github.Repo, commits []github.Commit) int {
	analyzer := health.NewHealthAnalyzer()
	score, _ := analyzer.CalculateRepositoryHealth(repo, commits)
	return int(score)
}
