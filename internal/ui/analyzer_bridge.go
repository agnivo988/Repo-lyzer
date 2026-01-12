package ui

import (
	"github.com/agnivo988/Repo-lyzer/internal/analyzer"
	"github.com/agnivo988/Repo-lyzer/internal/github"
)

// AnalyzerDataBridge provides clean interface between analyzer logic and UI
type AnalyzerDataBridge struct {
	repo          *github.Repo
	commits       []github.Commit
	contributors  []github.Contributor
	languages     map[string]int
	healthScore   int
	busFactor     int
	busRisk       string
	maturityScore int
	maturityLevel string
	fileTree      *FileNode
	cache         map[string]interface{}
}

// NewAnalyzerDataBridge creates a new data bridge with analyzer results.
// It initializes the bridge with repository data, metrics, and analysis results
// to provide a clean interface between the analyzer logic and UI components.
// Parameters:
//   - result: AnalysisResult containing all repository analysis data
// Returns a pointer to the initialized AnalyzerDataBridge.
func NewAnalyzerDataBridge(result AnalysisResult) *AnalyzerDataBridge {
	return &AnalyzerDataBridge{
		repo:          result.Repo,
		commits:       result.Commits,
		contributors:  result.Contributors,
		languages:     result.Languages,
		healthScore:   result.HealthScore,
		busFactor:     result.BusFactor,
		busRisk:       result.BusRisk,
		maturityScore: result.MaturityScore,
		maturityLevel: result.MaturityLevel,
		fileTree:      BuildFileTree(result),
	}
}

// GetHealthMetrics returns health-related metrics
func (b *AnalyzerDataBridge) GetHealthMetrics() map[string]interface{} {
	return map[string]interface{}{
		"health_score":   b.healthScore,
		"health_status":  b.getHealthStatus(),
		"bus_factor":     b.busFactor,
		"bus_risk":       b.busRisk,
		"maturity_level": b.maturityLevel,
		"maturity_score": b.maturityScore,
		"health_color":   b.getHealthColor(),
		"risk_color":     b.getRiskColor(),
	}
}

// GetRepositoryInfo returns repository metadata
func (b *AnalyzerDataBridge) GetRepositoryInfo() map[string]interface{} {
	return map[string]interface{}{
		"name":           b.repo.FullName,
		"description":    b.repo.Description,
		"stars":          b.repo.Stars,
		"forks":          b.repo.Forks,
		"open_issues":    b.repo.OpenIssues,
		"default_branch": b.repo.DefaultBranch,
	}
}

// GetContributorMetrics returns contributor analysis
func (b *AnalyzerDataBridge) GetContributorMetrics() map[string]interface{} {
	topContributors := b.getTopContributors(5)
	return map[string]interface{}{
		"total_contributors": len(b.contributors),
		"top_contributors":   topContributors,
		"contributor_count":  len(b.contributors),
		"diversity_score":    b.calculateDiversity(),
	}
}

// GetCommitMetrics returns commit-related metrics
func (b *AnalyzerDataBridge) GetCommitMetrics() map[string]interface{} {
	commitActivity := analyzer.CommitsPerDay(b.commits)
	recentActivity := b.getRecentActivity()

	return map[string]interface{}{
		"total_commits":    len(b.commits),
		"commits_per_day":  commitActivity,
		"recent_activity":  recentActivity,
		"commit_frequency": b.calculateCommitFrequency(),
		"last_commit":      b.getLastCommitInfo(),
		"activity_trend":   b.calculateActivityTrend(),
	}
}

// GetLanguageMetrics returns programming language information
func (b *AnalyzerDataBridge) GetLanguageMetrics() map[string]interface{} {
	return map[string]interface{}{
		"languages":          b.languages,
		"primary_language":   b.getPrimaryLanguage(),
		"language_count":     len(b.languages),
		"language_diversity": b.calculateLanguageDiversity(),
	}
}

// GetCompleteAnalysis returns all metrics combined
func (b *AnalyzerDataBridge) GetCompleteAnalysis() map[string]interface{} {
	return map[string]interface{}{
		"repository":      b.GetRepositoryInfo(),
		"health":          b.GetHealthMetrics(),
		"contributors":    b.GetContributorMetrics(),
		"commits":         b.GetCommitMetrics(),
		"languages":       b.GetLanguageMetrics(),
		"summary":         b.GenerateSummary(),
		"recommendations": b.GenerateRecommendations(),
	}
}

// GetFileTree returns the repository file structure
func (b *AnalyzerDataBridge) GetFileTree() *FileNode {
	return b.fileTree
}

// --- Helper Methods ---

func (b *AnalyzerDataBridge) getHealthStatus() string {
	if b.healthScore >= 80 {
		return "Excellent"
	} else if b.healthScore >= 60 {
		return "Good"
	} else if b.healthScore >= 40 {
		return "Fair"
	}
	return "Poor"
}

func (b *AnalyzerDataBridge) getHealthColor() string {
	if b.healthScore >= 80 {
		return "green"
	} else if b.healthScore >= 60 {
		return "yellow"
	}
	return "red"
}

func (b *AnalyzerDataBridge) getRiskColor() string {
	if b.busFactor >= 7 {
		return "green"
	} else if b.busFactor >= 4 {
		return "yellow"
	}
	return "red"
}

func (b *AnalyzerDataBridge) getTopContributors(count int) []map[string]interface{} {
	if count > len(b.contributors) {
		count = len(b.contributors)
	}

	top := []map[string]interface{}{}
	for i := 0; i < count; i++ {
		contrib := b.contributors[i]
		top = append(top, map[string]interface{}{
			"login":         contrib.Login,
			"contributions": contrib.Commits,
		})
	}
	return top
}

func (b *AnalyzerDataBridge) calculateDiversity() float64 {
	if len(b.contributors) == 0 {
		return 0
	}

	var sum int
	for _, contrib := range b.contributors {
		sum += contrib.Commits
	}

	// Calculate Herfindahl index (diversity measure)
	var diversity float64
	for _, contrib := range b.contributors {
		ratio := float64(contrib.Commits) / float64(sum)
		diversity += ratio * ratio
	}

	// Convert to 0-100 scale (1 = perfect diversity, 0 = concentration)
	return (1 - diversity) * 100
}

func (b *AnalyzerDataBridge) getRecentActivity() map[string]int {
	activity := make(map[string]int)
	if len(b.commits) == 0 {
		return activity
	}

	// This would be implemented based on commit timestamps
	return activity
}

func (b *AnalyzerDataBridge) calculateCommitFrequency() string {
	if len(b.commits) == 0 {
		return "No commits"
	}

	// Simplified frequency calculation
	avgPerDay := float64(len(b.commits)) / 365

	if avgPerDay >= 10 {
		return "Very High"
	} else if avgPerDay >= 5 {
		return "High"
	} else if avgPerDay >= 1 {
		return "Regular"
	}
	return "Sporadic"
}

func (b *AnalyzerDataBridge) getLastCommitInfo() map[string]interface{} {
	if len(b.commits) == 0 {
		return map[string]interface{}{}
	}

	lastCommit := b.commits[len(b.commits)-1]
	return map[string]interface{}{
		"sha":    lastCommit.SHA,
		"author": lastCommit.Commit.Author,
		"date":   lastCommit.Commit.Author.Date,
	}
}

func (b *AnalyzerDataBridge) calculateActivityTrend() string {
	if len(b.commits) < 2 {
		return "Unknown"
	}

	// Simplified trend calculation
	return "Stable"
}

func (b *AnalyzerDataBridge) getPrimaryLanguage() string {
	// Language field not available in Repo struct, use languages map instead
	maxBytes := 0
	primaryLang := "Unknown"
	for lang, bytes := range b.languages {
		if bytes > maxBytes {
			maxBytes = bytes
			primaryLang = lang
		}
	}

	return primaryLang
}

func (b *AnalyzerDataBridge) calculateLanguageDiversity() float64 {
	if len(b.languages) == 0 {
		return 0
	}

	var totalBytes int64
	for _, bytes := range b.languages {
		totalBytes += int64(bytes)
	}

	var diversity float64
	for _, bytes := range b.languages {
		ratio := float64(bytes) / float64(totalBytes)
		diversity += ratio * ratio
	}

	// Convert to 0-100 scale
	return (1 - diversity) * 100
}

// GenerateSummary creates a text summary of the analysis
func (b *AnalyzerDataBridge) GenerateSummary() string {
	summary := "📊 Analysis Summary:\n\n"

	// Health assessment
	if b.healthScore >= 80 {
		summary += "✅ This repository has excellent health metrics.\n"
	} else if b.healthScore >= 60 {
		summary += "⚠️ This repository has good health but room for improvement.\n"
	} else {
		summary += "❌ This repository needs attention in several areas.\n"
	}

	// Bus factor assessment
	if b.busFactor <= 2 {
		summary += "🚌 WARNING: High dependency on few contributors.\n"
	} else if b.busFactor <= 4 {
		summary += "⚠️ Some concentration of key contributors.\n"
	} else {
		summary += "✅ Good distribution of contributor responsibility.\n"
	}

	// Activity assessment
	switch b.calculateCommitFrequency() {
	case "Very High":
		summary += "📈 Very active development pace.\n"
	case "High":
		summary += "📈 Active development pace.\n"
	case "Regular":
		summary += "→ Regular maintenance activity.\n"
	default:
		summary += "→ Sporadic update activity.\n"
	}

	// Maturity assessment
	summary += "📚 Maturity Level: " + b.maturityLevel + "\n"

	return summary
}

// GenerateRecommendations creates actionable recommendations
func (b *AnalyzerDataBridge) GenerateRecommendations() []string {
	recommendations := []string{}

	// Health-based recommendations
	if b.healthScore < 60 {
		recommendations = append(recommendations, "Improve commit frequency and consistency")
		recommendations = append(recommendations, "Address open issues and manage problem backlog")
	}

	// Bus factor recommendations
	if b.busFactor <= 2 {
		recommendations = append(recommendations, "Recruit and onboard more contributors")
		recommendations = append(recommendations, "Document critical processes and architecture")
	}

	// Activity-based recommendations
	if b.calculateCommitFrequency() == "Sporadic" {
		recommendations = append(recommendations, "Establish regular development schedule")
		recommendations = append(recommendations, "Plan and track issues more systematically")
	}

	// Language diversity recommendations
	diversity := b.calculateLanguageDiversity()
	if diversity > 70 {
		recommendations = append(recommendations, "Consider consolidating technology stack")
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "Repository is well-maintained. Continue current practices.")
	}

	return recommendations
}
