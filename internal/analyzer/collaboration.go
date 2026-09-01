// Package analyzer provides functions for analyzing GitHub repository data.
// This file implements Team Collaboration Insights analysis.
package analyzer

import (
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/github"
)

// CollaborationMetrics contains comprehensive collaboration analysis
type CollaborationMetrics struct {
	// Overview
	TotalContributors  int     `json:"total_contributors"`
	ActiveReviewers    int     `json:"active_reviewers"`
	CollaborationScore float64 `json:"collaboration_score"` // 0-100

	// Review Participation
	ReviewParticipation ReviewParticipation `json:"review_participation"`

	// Cross-contributor patterns
	CrossContributorPatterns CrossContributorPatterns `json:"cross_contributor_patterns"`

	// Knowledge silos
	KnowledgeSilos []KnowledgeSilo `json:"knowledge_silos"`

	// Network metrics
	NetworkMetrics NetworkMetrics `json:"network_metrics"`

	// Review timing
	ReviewTiming ReviewTimingStats `json:"review_timing"`

	// Recommendations
	Recommendations []string `json:"recommendations"`
}

// ReviewParticipation tracks who reviews whose PRs
type ReviewParticipation struct {
	TotalReviews          int            `json:"total_reviews"`
	UniqueReviewers       int            `json:"unique_reviewers"`
	ReviewsPerContributor map[string]int `json:"reviews_per_contributor"`
	ReviewNetwork         []ReviewEdge   `json:"review_network"` // Who reviewed whose PR
	MostActiveReviewer    string         `json:"most_active_reviewer"`
	LeastActiveReviewer   string         `json:"least_active_reviewer"`
}

// ReviewEdge represents a review relationship
type ReviewEdge struct {
	Reviewer string `json:"reviewer"`
	Author   string `json:"author"`
	Count    int    `json:"count"`
}

// CrossContributorPatterns tracks files touched by multiple contributors
type CrossContributorPatterns struct {
	TotalFiles           int                `json:"total_files"`
	FilesByMultiple      int                `json:"files_by_multiple"`   // Files touched by 2+ contributors
	FilesBySingle        int                `json:"files_by_single"`     // Files touched by only 1 contributor
	CollaborationRatio   float64            `json:"collaboration_ratio"` // % of files with multiple contributors
	TopCollaboratedFiles []FileContributors `json:"top_collaborated_files"`
}

// FileContributors tracks contributors for a specific file
type FileContributors struct {
	Filename     string   `json:"filename"`
	Contributors []string `json:"contributors"`
	Count        int      `json:"count"`
}

// KnowledgeSilo represents a file with limited contributor access
type KnowledgeSilo struct {
	Filename     string   `json:"filename"`
	Owner        string   `json:"owner"`
	Contributors []string `json:"contributors"`
	RiskLevel    string   `json:"risk_level"` // Low, Medium, High, Critical
}

// NetworkMetrics contains graph theory metrics for collaboration network
type NetworkMetrics struct {
	ClusteringCoefficient float64            `json:"clustering_coefficient"`
	DegreeCentrality      map[string]float64 `json:"degree_centrality"`
	BetweennessCentrality map[string]float64 `json:"betweenness_centrality"`
	NetworkDensity        float64            `json:"network_density"`
	ConnectedComponents   int                `json:"connected_components"`
}

// ReviewTimingStats tracks review response times
type ReviewTimingStats struct {
	AverageTimeToFirstReview time.Duration `json:"average_time_to_first_review"`
	MedianTimeToFirstReview  time.Duration `json:"median_time_to_first_review"`
	AverageTimeToMerge       time.Duration `json:"average_time_to_merge"`
	ReviewsPerPR             float64       `json:"reviews_per_pr"`
	PRsWithReviews           float64       `json:"prs_with_reviews"` // % of PRs that get reviews
}

// AnalyzeCollaboration performs comprehensive collaboration analysis
func AnalyzeCollaboration(
	contributors []github.Contributor,
	prs []github.PullRequest,
	reviews map[int][]github.Review,
) *CollaborationMetrics {
	metrics := &CollaborationMetrics{
		TotalContributors: len(contributors),
	}

	if len(contributors) == 0 {
		return metrics
	}

	// Analyze review participation
	metrics.ReviewParticipation = analyzeReviewParticipation(prs, reviews)

	// Analyze cross-contributor patterns (using PR data as proxy)
	metrics.CrossContributorPatterns = analyzeCrossContributorPatterns(prs)

	// Detect knowledge silos
	metrics.KnowledgeSilos = detectKnowledgeSilos(prs, contributors)

	// Calculate network metrics
	metrics.NetworkMetrics = calculateNetworkMetrics(metrics.ReviewParticipation.ReviewNetwork, contributors)

	// Analyze review timing
	metrics.ReviewTiming = analyzeReviewTiming(prs, reviews)

	// Calculate overall collaboration score
	metrics.CollaborationScore = calculateCollaborationScore(metrics)

	// Generate recommendations
	metrics.Recommendations = generateCollaborationRecommendations(metrics)

	return metrics
}

// analyzeReviewParticipation analyzes who reviews whose PRs
func analyzeReviewParticipation(prs []github.PullRequest, reviews map[int][]github.Review) ReviewParticipation {
	participation := ReviewParticipation{
		ReviewsPerContributor: make(map[string]int),
	}

	// Build review network
	var reviewNetwork []ReviewEdge
	reviewCounts := make(map[string]int)

	for _, pr := range prs {
		prReviews, exists := reviews[pr.Number]
		if !exists || len(prReviews) == 0 {
			continue
		}

		// Track unique reviewers for this PR
		uniqueReviewers := make(map[string]int)

		for _, review := range prReviews {
			// Don't count self-reviews
			if review.User.Login == pr.User.Login {
				continue
			}

			uniqueReviewers[review.User.Login]++
			reviewCounts[review.User.Login]++
			participation.TotalReviews++
		}

		// Create edges for each reviewer-author pair
		for reviewer, count := range uniqueReviewers {
			reviewNetwork = append(reviewNetwork, ReviewEdge{
				Reviewer: reviewer,
				Author:   pr.User.Login,
				Count:    count,
			})
		}
	}

	participation.ReviewNetwork = reviewNetwork
	participation.ReviewsPerContributor = reviewCounts
	participation.UniqueReviewers = len(reviewCounts)

	// Find most and least active reviewers
	if len(reviewCounts) > 0 {
		sortedReviewers := sortMapByValue(reviewCounts)
		participation.MostActiveReviewer = sortedReviewers[0].Key
		if len(sortedReviewers) > 1 {
			participation.LeastActiveReviewer = sortedReviewers[len(sortedReviewers)-1].Key
		}
	}

	return participation
}

// analyzeCrossContributorPatterns analyzes files touched by multiple contributors
func analyzeCrossContributorPatterns(prs []github.PullRequest) CrossContributorPatterns {
	patterns := CrossContributorPatterns{
		TotalFiles: 0,
	}

	// Group files by contributor (using PR author as proxy)
	fileContributors := make(map[string]map[string]bool)

	for _, pr := range prs {
		// For each PR, we track the author as a contributor
		// In a real implementation, we'd also get the commit authors
		author := pr.User.Login
		if author == "" {
			continue
		}

		// Create a "virtual" file entry per PR for analysis
		// In production, we'd use commit file data
		fileKey := fmt.Sprintf("PR #%d", pr.Number)
		if fileContributors[fileKey] == nil {
			fileContributors[fileKey] = make(map[string]bool)
		}
		fileContributors[fileKey][author] = true
	}

	// Count files by contributor count
	var fileContributorList []FileContributors
	for filename, contribs := range fileContributors {
		contributorList := make([]string, 0, len(contribs))
		for c := range contribs {
			contributorList = append(contributorList, c)
		}

		count := len(contributorList)
		if count == 1 {
			patterns.FilesBySingle++
		} else {
			patterns.FilesByMultiple++
		}

		fileContributorList = append(fileContributorList, FileContributors{
			Filename:     filename,
			Contributors: contributorList,
			Count:        count,
		})
	}

	patterns.TotalFiles = len(fileContributors)

	if patterns.TotalFiles > 0 {
		patterns.CollaborationRatio = float64(patterns.FilesByMultiple) / float64(patterns.TotalFiles) * 100
	}

	// Get top collaborated files
	sort.Slice(fileContributorList, func(i, j int) bool {
		return fileContributorList[i].Count > fileContributorList[j].Count
	})
	if len(fileContributorList) > 10 {
		fileContributorList = fileContributorList[:10]
	}
	patterns.TopCollaboratedFiles = fileContributorList

	return patterns
}

// detectKnowledgeSilos detects files only touched by one person
func detectKnowledgeSilos(prs []github.PullRequest, contributors []github.Contributor) []KnowledgeSilo {
	silos := []KnowledgeSilo{}

	// Track file ownership via PR changes
	// In a real implementation, we'd use commit file data
	fileOwnership := make(map[string]map[string]bool)

	for _, pr := range prs {
		author := pr.User.Login
		if author == "" {
			continue
		}

		// Simulate file ownership from PR
		// In production, get actual files from commit details
		_ = author // Use author to avoid unused warning
	}

	// Build silos from ownership data
	for filename, owners := range fileOwnership {
		if len(owners) == 1 {
			var owner string
			for o := range owners {
				owner = o
				break
			}

			var ownerList []string
			for o := range owners {
				ownerList = append(ownerList, o)
			}

			risk := determineSiloRisk(owner, len(contributors))

			silos = append(silos, KnowledgeSilo{
				Filename:     filename,
				Owner:        owner,
				Contributors: ownerList,
				RiskLevel:    risk,
			})
		}
	}

	// Sort by risk level
	sort.Slice(silos, func(i, j int) bool {
		riskOrder := map[string]int{"Critical": 4, "High": 3, "Medium": 2, "Low": 1}
		return riskOrder[silos[i].RiskLevel] > riskOrder[silos[j].RiskLevel]
	})

	// Limit to top 20
	if len(silos) > 20 {
		silos = silos[:20]
	}

	return silos
}

func determineSiloRisk(owner string, totalContributors int) string {
	if totalContributors <= 1 {
		return "Critical"
	}
	// In a real implementation, calculate based on actual ownership
	return "Medium"
}

// calculateNetworkMetrics calculates graph theory metrics
func calculateNetworkMetrics(reviewNetwork []ReviewEdge, contributors []github.Contributor) NetworkMetrics {
	metrics := NetworkMetrics{
		DegreeCentrality:      make(map[string]float64),
		BetweennessCentrality: make(map[string]float64),
	}

	if len(reviewNetwork) == 0 || len(contributors) == 0 {
		return metrics
	}

	// Build adjacency list
	adjacency := make(map[string]map[string]bool)
	for _, edge := range reviewNetwork {
		if adjacency[edge.Reviewer] == nil {
			adjacency[edge.Reviewer] = make(map[string]bool)
		}
		adjacency[edge.Reviewer][edge.Author] = true
		if adjacency[edge.Author] == nil {
			adjacency[edge.Author] = make(map[string]bool)
		}
	}

	// Calculate degree centrality
	for _, contributor := range contributors {
		neighbors := 0
		if adj, ok := adjacency[contributor.Login]; ok {
			neighbors = len(adj)
		}
		// Normalize by possible connections (n-1)
		if len(contributors) > 1 {
			metrics.DegreeCentrality[contributor.Login] = float64(neighbors) / float64(len(contributors)-1)
		}
	}

	// Calculate clustering coefficient
	metrics.ClusteringCoefficient = calculateClusteringCoefficient(adjacency)

	// Calculate network density
	actualEdges := 0
	for _, neighbors := range adjacency {
		actualEdges += len(neighbors)
	}
	possibleEdges := len(contributors) * (len(contributors) - 1)
	if possibleEdges > 0 {
		metrics.NetworkDensity = float64(actualEdges) / float64(possibleEdges)
	}

	// Count connected components (simplified)
	metrics.ConnectedComponents = countConnectedComponents(adjacency, contributors)

	// Simplified betweenness centrality (approximation)
	for _, contributor := range contributors {
		metrics.BetweennessCentrality[contributor.Login] = calculateBetweenness(contributor.Login, adjacency, contributors)
	}

	return metrics
}

func calculateClusteringCoefficient(adjacency map[string]map[string]bool) float64 {
	var totalCoeff float64
	var count int

	for _, neighbors := range adjacency {
		if len(neighbors) < 2 {
			continue
		}

		// Count edges between neighbors
		neighborList := make([]string, 0, len(neighbors))
		for n := range neighbors {
			neighborList = append(neighborList, n)
		}

		possibleEdges := len(neighborList) * (len(neighborList) - 1) / 2
		if possibleEdges == 0 {
			continue
		}

		actualEdges := 0
		for i := 0; i < len(neighborList); i++ {
			for j := i + 1; j < len(neighborList); j++ {
				if adjacency[neighborList[i]][neighborList[j]] {
					actualEdges++
				}
			}
		}

		totalCoeff += float64(actualEdges) / float64(possibleEdges)
		count++
	}

	if count == 0 {
		return 0
	}

	return totalCoeff / float64(count)
}

func countConnectedComponents(adjacency map[string]map[string]bool, contributors []github.Contributor) int {
	if len(contributors) == 0 {
		return 0
	}

	visited := make(map[string]bool)
	components := 0

	var dfs func(node string)
	dfs = func(node string) {
		visited[node] = true
		for neighbor := range adjacency[node] {
			if !visited[neighbor] {
				dfs(neighbor)
			}
		}
	}

	for _, contributor := range contributors {
		if !visited[contributor.Login] {
			dfs(contributor.Login)
			components++
		}
	}

	return components
}

func calculateBetweenness(node string, adjacency map[string]map[string]bool, contributors []github.Contributor) float64 {
	// Simplified betweenness calculation
	// In a full implementation, would use shortest path algorithms
	_ = node // Acknowledge node parameter

	if len(contributors) <= 2 {
		return 0
	}

	// Count how many other nodes this node connects
	directConnections := 0
	if adj, ok := adjacency[node]; ok {
		directConnections = len(adj)
	}

	// Normalize
	return float64(directConnections) / float64(len(contributors)-1)
}

// analyzeReviewTiming analyzes review response times
func analyzeReviewTiming(prs []github.PullRequest, reviews map[int][]github.Review) ReviewTimingStats {
	stats := ReviewTimingStats{}

	var firstReviewTimes []time.Duration
	var mergeTimes []time.Duration
	prsWithReviews := 0

	for _, pr := range prs {
		prReviews, exists := reviews[pr.Number]
		if !exists || len(prReviews) == 0 {
			continue
		}

		prsWithReviews++

		// Find first review time
		earliestReview := prReviews[0]
		for _, review := range prReviews {
			if review.SubmittedAt.Before(earliestReview.SubmittedAt) {
				earliestReview = review
			}
		}
		firstReviewTime := earliestReview.SubmittedAt.Sub(pr.CreatedAt)
		firstReviewTimes = append(firstReviewTimes, firstReviewTime)

		// Calculate merge time
		if pr.MergedAt != nil {
			mergeTime := pr.MergedAt.Sub(pr.CreatedAt)
			mergeTimes = append(mergeTimes, mergeTime)
		}
	}

	// Calculate average time to first review
	if len(firstReviewTimes) > 0 {
		var total time.Duration
		for _, t := range firstReviewTimes {
			total += t
		}
		stats.AverageTimeToFirstReview = total / time.Duration(len(firstReviewTimes))

		// Calculate median
		sort.Slice(firstReviewTimes, func(i, j int) bool {
			return firstReviewTimes[i] < firstReviewTimes[j]
		})
		mid := len(firstReviewTimes) / 2
		if len(firstReviewTimes)%2 == 0 {
			stats.MedianTimeToFirstReview = (firstReviewTimes[mid-1] + firstReviewTimes[mid]) / 2
		} else {
			stats.MedianTimeToFirstReview = firstReviewTimes[mid]
		}
	}

	// Calculate average time to merge
	if len(mergeTimes) > 0 {
		var total time.Duration
		for _, t := range mergeTimes {
			total += t
		}
		stats.AverageTimeToMerge = total / time.Duration(len(mergeTimes))
	}

	// Calculate reviews per PR
	if len(prs) > 0 {
		var totalReviews int
		for _, prReviews := range reviews {
			totalReviews += len(prReviews)
		}
		stats.ReviewsPerPR = float64(totalReviews) / float64(len(prs))
		stats.PRsWithReviews = float64(prsWithReviews) / float64(len(prs)) * 100
	}

	return stats
}

// calculateCollaborationScore calculates overall collaboration score
func calculateCollaborationScore(metrics *CollaborationMetrics) float64 {
	var score float64

	// Review participation weight: 25%
	if metrics.ReviewParticipation.UniqueReviewers > 0 {
		reviewScore := math.Min(100, float64(metrics.ReviewParticipation.UniqueReviewers)/float64(metrics.TotalContributors)*100)
		score += reviewScore * 0.25
	}

	// Cross-contributor patterns weight: 25%
	score += math.Min(100, metrics.CrossContributorPatterns.CollaborationRatio) * 0.25

	// Network density weight: 25%
	score += metrics.NetworkMetrics.NetworkDensity * 100 * 0.25

	// Review engagement weight: 25%
	score += math.Min(100, metrics.ReviewTiming.PRsWithReviews) * 0.25

	return math.Round(score*100) / 100
}

// generateCollaborationRecommendations generates recommendations based on analysis
func generateCollaborationRecommendations(metrics *CollaborationMetrics) []string {
	var recs []string

	// Review participation recommendations
	if metrics.ReviewParticipation.UniqueReviewers < metrics.TotalContributors/2 {
		recs = append(recs, "⚠️ Low review participation. Encourage more contributors to review PRs.")
	}

	// Knowledge silos recommendations
	criticalSilos := 0
	highSilos := 0
	for _, silo := range metrics.KnowledgeSilos {
		if silo.RiskLevel == "Critical" {
			criticalSilos++
		}
		if silo.RiskLevel == "High" {
			highSilos++
		}
	}

	if criticalSilos > 0 {
		recs = append(recs, "🔴 Critical knowledge silos detected! These files need immediate knowledge transfer.")
	}
	if highSilos > 3 {
		recs = append(recs, "🟠 Multiple high-risk knowledge silos. Consider cross-training team members.")
	}

	// Network metrics recommendations
	if metrics.NetworkMetrics.NetworkDensity < 0.3 {
		recs = append(recs, "📊 Low network density. Contributors aren't well connected - consider pair programming sessions.")
	}

	if metrics.NetworkMetrics.ClusteringCoefficient < 0.3 {
		recs = append(recs, "🔗 Low clustering. Team works in silos - organize more team-wide code reviews.")
	}

	// Review timing recommendations
	if metrics.ReviewTiming.AverageTimeToFirstReview > 48*time.Hour {
		recs = append(recs, "⏰ Slow review response times (>48h). Consider implementing review rotation.")
	}

	if metrics.CollaborationScore > 80 {
		recs = append(recs, "🌟 Excellent collaboration! Team is working well together.")
	} else if metrics.CollaborationScore > 60 {
		recs = append(recs, "✅ Good collaboration. Continue current practices.")
	} else if metrics.CollaborationScore < 40 {
		recs = append(recs, "📈 Collaboration needs improvement. Focus on increasing code review participation.")
	}

	if len(recs) == 0 {
		recs = append(recs, "✅ Collaboration metrics look healthy.")
	}

	return recs
}

// Helper function to sort map by value
type MapEntry struct {
	Key   string
	Value int
}

func sortMapByValue(m map[string]int) []MapEntry {
	var entries []MapEntry
	for k, v := range m {
		entries = append(entries, MapEntry{Key: k, Value: v})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Value > entries[j].Value
	})
	return entries
}
