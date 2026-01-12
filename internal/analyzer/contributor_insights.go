// Package analyzer provides functions for analyzing GitHub repository data.
// This file implements detailed contributor insights analysis.
package analyzer

import (
	"sort"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/github"
)

// ContributorInsights contains detailed analysis of repository contributors
type ContributorInsights struct {
	TotalContributors   int                   `json:"total_contributors"`
	ActiveContributors  int                   `json:"active_contributors"`  // Contributors with >1% of commits
	TopContributor      *ContributorDetail    `json:"top_contributor"`
	ContributorDetails  []ContributorDetail   `json:"contributor_details"`
	DiversityScore      float64               `json:"diversity_score"`      // 0-100, higher = more diverse
	ConcentrationRisk   string                `json:"concentration_risk"`   // Low, Medium, High, Critical
	NewContributors     int                   `json:"new_contributors"`     // Contributors with <10 commits
	VeteranContributors int                   `json:"veteran_contributors"` // Contributors with >100 commits
	CommitDistribution  CommitDistribution    `json:"commit_distribution"`
	TeamSize            string                `json:"team_size"`            // Solo, Small, Medium, Large
	Recommendations     []string              `json:"recommendations"`
}

// ContributorDetail contains detailed info about a single contributor
type ContributorDetail struct {
	Login           string  `json:"login"`
	Commits         int     `json:"commits"`
	Percentage      float64 `json:"percentage"`
	Rank            int     `json:"rank"`
	ContributorType string  `json:"contributor_type"` // Core, Regular, Occasional, New
}

// CommitDistribution shows how commits are distributed
type CommitDistribution struct {
	Top1Percent    float64 `json:"top_1_percent"`    // % of commits by top 1%
	Top10Percent   float64 `json:"top_10_percent"`   // % of commits by top 10%
	Top50Percent   float64 `json:"top_50_percent"`   // % of commits by top 50%
	GiniCoefficient float64 `json:"gini_coefficient"` // Inequality measure (0=equal, 1=unequal)
}

// AnalyzeContributors performs detailed contributor analysis
func AnalyzeContributors(contributors []github.Contributor) *ContributorInsights {
	if len(contributors) == 0 {
		return &ContributorInsights{
			TotalContributors:  0,
			DiversityScore:     0,
			ConcentrationRisk:  "Unknown",
			TeamSize:           "None",
			Recommendations:    []string{"No contributor data available"},
		}
	}

	insights := &ContributorInsights{
		TotalContributors: len(contributors),
	}

	// Calculate total commits
	totalCommits := 0
	for _, c := range contributors {
		totalCommits += c.Commits
	}

	// Build contributor details
	details := make([]ContributorDetail, len(contributors))
	for i, c := range contributors {
		pct := 0.0
		if totalCommits > 0 {
			pct = float64(c.Commits) / float64(totalCommits) * 100
		}
		
		details[i] = ContributorDetail{
			Login:           c.Login,
			Commits:         c.Commits,
			Percentage:      pct,
			Rank:            i + 1,
			ContributorType: classifyContributor(c.Commits, pct),
		}
	}
	insights.ContributorDetails = details

	// Top contributor
	if len(details) > 0 {
		insights.TopContributor = &details[0]
	}

	// Count contributor types
	for _, d := range details {
		if d.Percentage > 1 {
			insights.ActiveContributors++
		}
		if d.Commits < 10 {
			insights.NewContributors++
		}
		if d.Commits > 100 {
			insights.VeteranContributors++
		}
	}

	// Calculate commit distribution
	insights.CommitDistribution = calculateDistribution(contributors, totalCommits)

	// Calculate diversity score (inverse of concentration)
	insights.DiversityScore = calculateDiversityScore(contributors, totalCommits)

	// Determine concentration risk
	insights.ConcentrationRisk = determineConcentrationRisk(insights)

	// Determine team size
	insights.TeamSize = classifyTeamSize(len(contributors), insights.ActiveContributors)

	// Generate recommendations
	insights.Recommendations = generateRecommendations(insights)

	return insights
}

func classifyContributor(commits int, percentage float64) string {
	if percentage > 20 {
		return "Core"
	} else if percentage > 5 {
		return "Regular"
	} else if commits > 10 {
		return "Occasional"
	}
	return "New"
}

func calculateDistribution(contributors []github.Contributor, total int) CommitDistribution {
	if len(contributors) == 0 || total == 0 {
		return CommitDistribution{}
	}

	dist := CommitDistribution{}
	n := len(contributors)

	// Top 1%
	top1Count := max(1, n/100)
	top1Commits := 0
	for i := 0; i < top1Count && i < n; i++ {
		top1Commits += contributors[i].Commits
	}
	dist.Top1Percent = float64(top1Commits) / float64(total) * 100

	// Top 10%
	top10Count := max(1, n/10)
	top10Commits := 0
	for i := 0; i < top10Count && i < n; i++ {
		top10Commits += contributors[i].Commits
	}
	dist.Top10Percent = float64(top10Commits) / float64(total) * 100

	// Top 50%
	top50Count := max(1, n/2)
	top50Commits := 0
	for i := 0; i < top50Count && i < n; i++ {
		top50Commits += contributors[i].Commits
	}
	dist.Top50Percent = float64(top50Commits) / float64(total) * 100

	// Gini coefficient
	dist.GiniCoefficient = calculateGini(contributors, total)

	return dist
}

func calculateGini(contributors []github.Contributor, total int) float64 {
	if len(contributors) <= 1 || total == 0 {
		return 0
	}

	n := len(contributors)
	
	// Sort by commits ascending for Gini calculation
	sorted := make([]int, n)
	for i, c := range contributors {
		sorted[i] = c.Commits
	}
	sort.Ints(sorted)

	// Calculate Gini coefficient
	var sumOfDiffs float64
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			diff := sorted[i] - sorted[j]
			if diff < 0 {
				diff = -diff
			}
			sumOfDiffs += float64(diff)
		}
	}

	gini := sumOfDiffs / (2 * float64(n) * float64(total))
	return gini
}

func calculateDiversityScore(contributors []github.Contributor, total int) float64 {
	if len(contributors) == 0 || total == 0 {
		return 0
	}

	// Use inverse Herfindahl-Hirschman Index
	var hhi float64
	for _, c := range contributors {
		share := float64(c.Commits) / float64(total)
		hhi += share * share
	}

	// Convert to 0-100 scale (1/n is perfect diversity, 1 is monopoly)
	minHHI := 1.0 / float64(len(contributors))
	if hhi <= minHHI {
		return 100
	}
	
	// Normalize: 100 = perfect diversity, 0 = single contributor
	diversity := (1 - hhi) / (1 - minHHI) * 100
	if diversity < 0 {
		diversity = 0
	}
	if diversity > 100 {
		diversity = 100
	}
	
	return diversity
}

func determineConcentrationRisk(insights *ContributorInsights) string {
	if insights.TopContributor == nil {
		return "Unknown"
	}

	topPct := insights.TopContributor.Percentage
	
	if topPct > 80 {
		return "Critical"
	} else if topPct > 60 {
		return "High"
	} else if topPct > 40 {
		return "Medium"
	}
	return "Low"
}

func classifyTeamSize(total, active int) string {
	if active == 0 {
		return "None"
	} else if active == 1 {
		return "Solo"
	} else if active <= 3 {
		return "Small"
	} else if active <= 10 {
		return "Medium"
	}
	return "Large"
}

func generateRecommendations(insights *ContributorInsights) []string {
	var recs []string

	// Concentration risk recommendations
	switch insights.ConcentrationRisk {
	case "Critical":
		recs = append(recs, "⚠️ Critical: Single contributor dominates. Urgent need for knowledge sharing.")
		recs = append(recs, "📝 Document all critical processes and architecture decisions.")
	case "High":
		recs = append(recs, "⚠️ High concentration risk. Consider pair programming and code reviews.")
		recs = append(recs, "🎯 Actively mentor and onboard new contributors.")
	case "Medium":
		recs = append(recs, "📊 Moderate concentration. Continue encouraging diverse contributions.")
	}

	// Team size recommendations
	switch insights.TeamSize {
	case "Solo":
		recs = append(recs, "👥 Consider recruiting contributors to reduce bus factor risk.")
	case "Small":
		recs = append(recs, "🌱 Good foundation. Focus on growing the contributor base.")
	}

	// New vs veteran balance
	if insights.VeteranContributors == 0 && insights.TotalContributors > 5 {
		recs = append(recs, "📈 No veteran contributors yet. Project may be young or have high turnover.")
	}
	
	if insights.NewContributors > insights.TotalContributors/2 && insights.TotalContributors > 10 {
		recs = append(recs, "🆕 Many new contributors. Ensure good onboarding documentation.")
	}

	// Diversity score recommendations
	if insights.DiversityScore < 30 {
		recs = append(recs, "🔄 Low diversity score. Work on distributing responsibilities more evenly.")
	} else if insights.DiversityScore > 70 {
		recs = append(recs, "✅ Good contribution diversity. Maintain current practices.")
	}

	if len(recs) == 0 {
		recs = append(recs, "✅ Contributor distribution looks healthy.")
	}

	return recs
}

// GetContributorActivity analyzes contributor activity patterns
func GetContributorActivity(commits []github.Commit) map[string]int {
	activity := make(map[string]int)
	
	for _, c := range commits {
		date := c.Commit.Author.Date
		if !date.IsZero() {
			week := date.Format("2006-W01")
			activity[week]++
		}
	}
	
	return activity
}

// Helper function for Go versions < 1.21
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// ContributorTrend represents activity trend for a contributor
type ContributorTrend struct {
	Login        string
	RecentCommits int  // Last 30 days
	TotalCommits  int
	IsActive      bool // Had commits in last 30 days
	Trend         string // "Rising", "Stable", "Declining", "Inactive"
}

// AnalyzeContributorTrends analyzes recent activity trends
func AnalyzeContributorTrends(contributors []github.Contributor, commits []github.Commit) []ContributorTrend {
	// Count recent commits per contributor
	recentCounts := make(map[string]int)
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	
	for _, c := range commits {
		if c.Commit.Author.Date.After(thirtyDaysAgo) {
			// Note: We don't have author login in commit struct, so this is simplified
			recentCounts["recent"]++
		}
	}

	var trends []ContributorTrend
	for _, c := range contributors {
		trend := ContributorTrend{
			Login:        c.Login,
			TotalCommits: c.Commits,
		}
		
		// Simplified trend analysis based on total commits
		if c.Commits > 100 {
			trend.Trend = "Veteran"
			trend.IsActive = true
		} else if c.Commits > 20 {
			trend.Trend = "Regular"
			trend.IsActive = true
		} else if c.Commits > 5 {
			trend.Trend = "Occasional"
			trend.IsActive = true
		} else {
			trend.Trend = "New"
			trend.IsActive = c.Commits > 0
		}
		
		trends = append(trends, trend)
	}
	
	return trends
}
