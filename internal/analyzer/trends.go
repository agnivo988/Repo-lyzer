// Package analyzer provides functions for analyzing GitHub repository data.
// This file implements Repository Trend Analysis & Forecasting capabilities.
package analyzer

import (
	"sort"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/github"
)

// TrendIndicator represents the direction of a metric
type TrendIndicator string

const (
	// TrendImproving indicates the metric is improving/upward
	TrendImproving TrendIndicator = "↗️ Improving"
	// TrendDeclining indicates the metric is declining/downward
	TrendDeclining TrendIndicator = "↘️ Declining"
	// TrendStable indicates the metric is relatively stable
	TrendStable TrendIndicator = "➡️ Stable"
)

// MonthlyMetric stores aggregated data for a single month
type MonthlyMetric struct {
	Month         time.Time `json:"month"`
	Commits       int       `json:"commits"`
	Contributors  int       `json:"contributors"`
	IssuesOpened  int       `json:"issues_opened"`
	IssuesClosed  int       `json:"issues_closed"`
	PRsOpened     int       `json:"prs_opened"`
	PRsMerged     int       `json:"prs_merged"`
	AvgCommitSize float64   `json:"avg_commit_size"`
}

// TrendMetrics contains comprehensive trend analysis for a repository
type TrendMetrics struct {
	Owner          string          `json:"owner"`
	Repo           string          `json:"repo"`
	AnalysisPeriod int             `json:"analysis_period_months"`
	MonthlyData    []MonthlyMetric `json:"monthly_data"`

	// Commit Trends
	CommitTrend        TrendIndicator `json:"commit_trend"`
	CommitChangeRate   float64        `json:"commit_change_rate"`
	AvgCommitsPerMonth float64        `json:"avg_commits_per_month"`
	CommitTrendValues  []int          `json:"commit_trend_values"` // For sparkline

	// Contributor Trends
	ContributorTrend       TrendIndicator `json:"contributor_trend"`
	ContributorChangeRate  float64        `json:"contributor_change_rate"`
	CurrentContributors    int            `json:"current_contributors"`
	NewContributors        int            `json:"new_contributors"`
	LostContributors       int            `json:"lost_contributors"`
	ContributorTrendValues []int          `json:"contributor_trend_values"`

	// Issue Resolution Trends
	IssueResolutionTrend TrendIndicator `json:"issue_resolution_trend"`
	AvgResolutionTime    time.Duration  `json:"avg_resolution_time"`
	ResolutionRate       float64        `json:"resolution_rate"`

	// PR Trends
	PRTrend       TrendIndicator `json:"pr_trend"`
	PRMergeRate   float64        `json:"pr_merge_rate"`
	PRTrendValues []int          `json:"pr_trend_values"`

	// Health Score Prediction
	PredictedHealthScore int            `json:"predicted_health_score"`
	HealthScoreTrend     TrendIndicator `json:"health_score_trend"`
	CurrentHealthScore   int            `json:"current_health_score"`

	// Overall Assessment
	OverallTrend TrendIndicator `json:"overall_trend"`
	Summary      string         `json:"summary"`
}

// AnalyzeTrends performs comprehensive trend analysis on a repository
// It analyzes commit patterns, contributor growth, issue resolution, and PR trends
func AnalyzeTrends(
	owner, repo string,
	commits []github.Commit,
	contributors []github.Contributor,
	issues []github.Issue,
	prs []github.PullRequest,
	months int,
) *TrendMetrics {
	metrics := &TrendMetrics{
		Owner:                  owner,
		Repo:                   repo,
		AnalysisPeriod:         months,
		MonthlyData:            []MonthlyMetric{},
		CommitTrendValues:      []int{},
		ContributorTrendValues: []int{},
		PRTrendValues:          []int{},
	}

	if len(commits) == 0 && len(contributors) == 0 {
		metrics.Summary = "Insufficient data for trend analysis"
		return metrics
	}

	// Aggregate data by month
	metrics.MonthlyData = aggregateMonthlyData(commits, issues, prs, months)

	// Sort by month
	sort.Slice(metrics.MonthlyData, func(i, j int) bool {
		return metrics.MonthlyData[i].Month.Before(metrics.MonthlyData[j].Month)
	})

	// Analyze commit trends
	metrics.CommitTrend, metrics.CommitChangeRate, metrics.AvgCommitsPerMonth = AnalyzeCommitFrequencyTrends(metrics.MonthlyData)
	for _, m := range metrics.MonthlyData {
		metrics.CommitTrendValues = append(metrics.CommitTrendValues, m.Commits)
	}

	// Analyze contributor trends
	metrics.ContributorTrend, metrics.ContributorChangeRate, metrics.CurrentContributors, metrics.NewContributors, metrics.LostContributors = AnalyzeContributorTrendsForTrends(contributors, metrics.MonthlyData)
	for _, m := range metrics.MonthlyData {
		metrics.ContributorTrendValues = append(metrics.ContributorTrendValues, m.Contributors)
	}

	// Analyze issue trends
	resolutionTrend, avgResTime, resRate := analyzeIssueTrendsForTrends(issues, metrics.MonthlyData)
	metrics.IssueResolutionTrend = resolutionTrend
	metrics.AvgResolutionTime = avgResTime
	metrics.ResolutionRate = resRate

	// Analyze PR trends
	metrics.PRTrend, metrics.PRMergeRate = AnalyzePRTrends(prs, metrics.MonthlyData)
	for _, m := range metrics.MonthlyData {
		metrics.PRTrendValues = append(metrics.PRTrendValues, m.PRsMerged)
	}

	// Calculate current health score (simplified)
	metrics.CurrentHealthScore = calculateCurrentHealthScore(commits, contributors, issues, prs)

	// Predict future health score
	metrics.PredictedHealthScore = PredictHealthScore(metrics)
	metrics.HealthScoreTrend = determineHealthScoreTrend(metrics)

	// Determine overall trend
	metrics.OverallTrend, metrics.Summary = DetermineOverallTrend(metrics)

	return metrics
}

// aggregateMonthlyData groups commits, issues, and PRs by month
func aggregateMonthlyData(commits []github.Commit, issues []github.Issue, prs []github.PullRequest, months int) []MonthlyMetric {
	now := time.Now()
	since := now.AddDate(0, -months, 0)

	// Create monthly buckets
	monthlyData := make(map[string]*MonthlyMetric)

	// Initialize all months in the range
	for i := 0; i < months; i++ {
		month := now.AddDate(0, -i, 0)
		monthKey := month.Format("2006-01")
		monthlyData[monthKey] = &MonthlyMetric{
			Month: time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, time.UTC),
		}
	}

	// Aggregate commits
	for _, commit := range commits {
		date := commit.Commit.Author.Date
		if date.Before(since) {
			continue
		}
		monthKey := date.Format("2006-01")
		if m, ok := monthlyData[monthKey]; ok {
			m.Commits++
		}
	}

	// Aggregate issues
	for _, issue := range issues {
		created := issue.CreatedAt
		if created.Before(since) {
			continue
		}
		monthKey := created.Format("2006-01")
		if m, ok := monthlyData[monthKey]; ok {
			m.IssuesOpened++
		}
		if issue.State == "closed" && issue.ClosedAt != nil && !issue.ClosedAt.Before(since) {
			closedKey := issue.ClosedAt.Format("2006-01")
			if m, ok := monthlyData[closedKey]; ok {
				m.IssuesClosed++
			}
		}
	}

	// Aggregate PRs
	for _, pr := range prs {
		created := pr.CreatedAt
		if created.Before(since) {
			continue
		}
		monthKey := created.Format("2006-01")
		if m, ok := monthlyData[monthKey]; ok {
			m.PRsOpened++
			if pr.MergedAt != nil {
				m.PRsMerged++
			}
		}
	}

	// Convert to slice
	result := make([]MonthlyMetric, 0, len(monthlyData))
	for _, m := range monthlyData {
		result = append(result, *m)
	}

	return result
}

// AnalyzeCommitFrequencyTrends analyzes commit patterns over time
// Returns the trend direction, percentage change rate, and average commits per month
func AnalyzeCommitFrequencyTrends(monthlyData []MonthlyMetric) (TrendIndicator, float64, float64) {
	if len(monthlyData) < 2 {
		return TrendStable, 0, 0
	}

	// Get first and last month commits
	firstMonth := monthlyData[0].Commits
	lastMonth := monthlyData[len(monthlyData)-1].Commits

	// Calculate total and average
	var total int
	for _, m := range monthlyData {
		total += m.Commits
	}
	avgCommits := float64(total) / float64(len(monthlyData))

	// Calculate change rate
	var changeRate float64
	if firstMonth > 0 {
		changeRate = float64(lastMonth-firstMonth) / float64(firstMonth) * 100
	} else if lastMonth > 0 {
		changeRate = 100 // Went from 0 to some commits
	}

	// Determine trend
	trend := determineTrendFromRate(changeRate)

	return trend, changeRate, avgCommits
}

// AnalyzeContributorTrendsForTrends analyzes contributor growth/decline for trend analysis
func AnalyzeContributorTrendsForTrends(contributors []github.Contributor, monthlyData []MonthlyMetric) (TrendIndicator, float64, int, int, int) {
	if len(monthlyData) < 2 || len(contributors) == 0 {
		return TrendStable, 0, len(contributors), 0, 0
	}

	// For simplicity, use current contributor count and compare first/last month
	currentContributors := len(contributors)

	// Estimate based on monthly data (this is simplified - real implementation would track historical contributor data)
	firstMonthContributors := monthlyData[0].Contributors
	if firstMonthContributors == 0 {
		firstMonthContributors = currentContributors
	}
	lastMonthContributors := monthlyData[len(monthlyData)-1].Contributors
	if lastMonthContributors == 0 {
		lastMonthContributors = currentContributors
	}

	// Calculate change rate
	var changeRate float64
	if firstMonthContributors > 0 {
		changeRate = float64(lastMonthContributors-firstMonthContributors) / float64(firstMonthContributors) * 100
	}

	trend := determineTrendFromRate(changeRate)

	// Estimate new and lost contributors (simplified)
	newContributors := 0
	lostContributors := 0
	if lastMonthContributors > firstMonthContributors {
		newContributors = lastMonthContributors - firstMonthContributors
	} else if firstMonthContributors > lastMonthContributors {
		lostContributors = firstMonthContributors - lastMonthContributors
	}

	return trend, changeRate, currentContributors, newContributors, lostContributors
}

// analyzeIssueTrendsForTrends analyzes issue resolution velocity for trend analysis
func analyzeIssueTrendsForTrends(issues []github.Issue, monthlyData []MonthlyMetric) (TrendIndicator, time.Duration, float64) {
	if len(monthlyData) < 2 {
		return TrendStable, 0, 0
	}

	// Calculate resolution rate (closed / opened)
	var totalOpened, totalClosed int
	for _, m := range monthlyData {
		totalOpened += m.IssuesOpened
		totalClosed += m.IssuesClosed
	}

	var resolutionRate float64
	if totalOpened > 0 {
		resolutionRate = float64(totalClosed) / float64(totalOpened) * 100
	}

	// For trend, compare recent months to earlier months
	midPoint := len(monthlyData) / 2
	var earlyClosed, recentClosed int
	var earlyOpened, recentOpened int

	for i, m := range monthlyData {
		if i < midPoint {
			earlyOpened += m.IssuesOpened
			earlyClosed += m.IssuesClosed
		} else {
			recentOpened += m.IssuesOpened
			recentClosed += m.IssuesClosed
		}
	}

	var earlyRate, recentRate float64
	if earlyOpened > 0 {
		earlyRate = float64(earlyClosed) / float64(earlyOpened) * 100
	}
	if recentOpened > 0 {
		recentRate = float64(recentClosed) / float64(recentOpened) * 100
	}

	var trend TrendIndicator
	if recentRate > earlyRate+10 {
		trend = TrendImproving
	} else if recentRate < earlyRate-10 {
		trend = TrendDeclining
	} else {
		trend = TrendStable
	}

	// Compute average resolution time from closed issues with valid timestamps.
	// Fall back to a 7-day estimate only when no closed-issue data is available.
	avgResolution := 7 * 24 * time.Hour // fallback when there is no data
	var totalDuration time.Duration
	var closedCount int
	for _, issue := range issues {
		if issue.State == "closed" && issue.ClosedAt != nil {
			totalDuration += issue.ClosedAt.Sub(issue.CreatedAt)
			closedCount++
		}
	}
	if closedCount > 0 {
		avgResolution = totalDuration / time.Duration(closedCount)
	}

	return trend, avgResolution, resolutionRate
}

// AnalyzePRTrends analyzes pull request trends
func AnalyzePRTrends(prs []github.PullRequest, monthlyData []MonthlyMetric) (TrendIndicator, float64) {
	if len(monthlyData) < 2 {
		return TrendStable, 0
	}

	// Calculate merge rate
	var totalOpened, totalMerged int
	for _, m := range monthlyData {
		totalOpened += m.PRsOpened
		totalMerged += m.PRsMerged
	}

	var mergeRate float64
	if totalOpened > 0 {
		mergeRate = float64(totalMerged) / float64(totalOpened) * 100
	}

	// Compare recent to early periods
	midPoint := len(monthlyData) / 2
	var earlyMerged, recentMerged int

	for i, m := range monthlyData {
		if i < midPoint {
			earlyMerged += m.PRsMerged
		} else {
			recentMerged += m.PRsMerged
		}
	}

	earlyAvg := float64(earlyMerged) / float64(midPoint)
	recentAvg := float64(recentMerged) / float64(len(monthlyData)-midPoint)

	var trend TrendIndicator
	if recentAvg > earlyAvg*1.1 {
		trend = TrendImproving
	} else if recentAvg < earlyAvg*0.9 {
		trend = TrendDeclining
	} else {
		trend = TrendStable
	}

	return trend, mergeRate
}

// PredictHealthScore uses simple linear regression to predict future health score
func PredictHealthScore(metrics *TrendMetrics) int {
	// Create a simple health score based on various trends
	score := 50 // Base score

	// Adjust based on commit trend
	switch metrics.CommitTrend {
	case TrendImproving:
		score += 15
	case TrendDeclining:
		score -= 15
	}

	// Adjust based on contributor trend
	switch metrics.ContributorTrend {
	case TrendImproving:
		score += 15
	case TrendDeclining:
		score -= 15
	}

	// Adjust based on resolution rate
	if metrics.ResolutionRate > 70 {
		score += 10
	} else if metrics.ResolutionRate < 50 {
		score -= 10
	}

	// Adjust based on PR merge rate
	if metrics.PRMergeRate > 60 {
		score += 10
	} else if metrics.PRMergeRate < 40 {
		score -= 10
	}

	// Clamp to 0-100
	if score > 100 {
		score = 100
	}
	if score < 0 {
		score = 0
	}

	return score
}

// calculateCurrentHealthScore calculates the current health score based on available data
func calculateCurrentHealthScore(commits []github.Commit, contributors []github.Contributor, issues []github.Issue, prs []github.PullRequest) int {
	score := 50

	// Factor in commit activity
	if len(commits) > 100 {
		score += 20
	} else if len(commits) > 50 {
		score += 10
	}

	// Factor in contributor count
	if len(contributors) > 10 {
		score += 15
	} else if len(contributors) > 5 {
		score += 10
	}

	// Factor in PR merge rate
	if len(prs) > 0 {
		merged := 0
		for _, pr := range prs {
			if pr.MergedAt != nil {
				merged++
			}
		}
		mergeRate := float64(merged) / float64(len(prs)) * 100
		if mergeRate > 60 {
			score += 15
		} else if mergeRate > 40 {
			score += 10
		}
	}

	// Clamp to 0-100
	if score > 100 {
		score = 100
	}

	return score
}

// determineHealthScoreTrend determines the trend of health score
func determineHealthScoreTrend(metrics *TrendMetrics) TrendIndicator {
	diff := metrics.PredictedHealthScore - metrics.CurrentHealthScore

	if diff > 10 {
		return TrendImproving
	} else if diff < -10 {
		return TrendDeclining
	}

	return TrendStable
}

// DetermineOverallTrend determines the overall repository trend
func DetermineOverallTrend(metrics *TrendMetrics) (TrendIndicator, string) {
	// Count trends
	improvingCount := 0
	decliningCount := 0

	if metrics.CommitTrend == TrendImproving {
		improvingCount++
	} else if metrics.CommitTrend == TrendDeclining {
		decliningCount++
	}

	if metrics.ContributorTrend == TrendImproving {
		improvingCount++
	} else if metrics.ContributorTrend == TrendDeclining {
		decliningCount++
	}

	if metrics.IssueResolutionTrend == TrendImproving {
		improvingCount++
	} else if metrics.IssueResolutionTrend == TrendDeclining {
		decliningCount++
	}

	if metrics.PRTrend == TrendImproving {
		improvingCount++
	} else if metrics.PRTrend == TrendDeclining {
		decliningCount++
	}

	// Determine overall trend
	var overallTrend TrendIndicator
	var summary string

	if improvingCount > decliningCount+1 {
		overallTrend = TrendImproving
		summary = "Repository shows positive momentum across multiple metrics"
	} else if decliningCount > improvingCount+1 {
		overallTrend = TrendDeclining
		summary = "Repository shows declining activity - may need attention"
	} else if improvingCount > 0 {
		overallTrend = TrendImproving
		summary = "Repository is maintaining steady activity with some improvements"
	} else if decliningCount > 0 {
		overallTrend = TrendDeclining
		summary = "Repository is experiencing some decline in activity"
	} else {
		overallTrend = TrendStable
		summary = "Repository activity is stable"
	}

	// Add specific insights
	if metrics.CommitTrend == TrendDeclining {
		summary += ". Commit frequency is decreasing."
	}
	if metrics.ContributorTrend == TrendDeclining {
		summary += " Contributor activity is declining."
	}
	if metrics.PRTrend == TrendImproving {
		summary += " Pull request activity is increasing."
	}

	return overallTrend, summary
}

// determineTrendFromRate determines trend indicator from percentage change rate
func determineTrendFromRate(changeRate float64) TrendIndicator {
	if changeRate > 10 {
		return TrendImproving
	} else if changeRate < -10 {
		return TrendDeclining
	}
	return TrendStable
}

// SimpleLinearRegression performs simple linear regression on data points
// Returns slope, intercept, and R-squared value
func SimpleLinearRegression(points []struct {
	X float64
	Y float64
}) (slope, intercept, rSquared float64) {
	if len(points) < 2 {
		return 0, 0, 0
	}

	n := float64(len(points))

	// Calculate means
	var sumX, sumY float64
	for _, p := range points {
		sumX += p.X
		sumY += p.Y
	}
	meanX := sumX / n
	meanY := sumY / n

	// Calculate slope and intercept
	var sumXY, sumX2, sumY2 float64
	for _, p := range points {
		dx := p.X - meanX
		dy := p.Y - meanY
		sumXY += dx * dy
		sumX2 += dx * dx
		sumY2 += dy * dy
	}

	if sumX2 == 0 {
		return 0, meanY, 0
	}

	slope = sumXY / sumX2
	intercept = meanY - slope*meanX

	// Calculate R-squared
	if sumY2 > 0 {
		rSquared = (sumXY * sumXY) / (sumX2 * sumY2)
	}

	return slope, intercept, rSquared
}

// PredictFutureValue predicts future Y value using linear regression
func PredictFutureValue(points []struct {
	X float64
	Y float64
}, futureX float64) float64 {
	slope, intercept, _ := SimpleLinearRegression(points)
	return slope*futureX + intercept
}
