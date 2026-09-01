package predictive

import (
	"math"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/github"
)

// PRProbability represents the likelihood of a Pull Request being merged
type PRProbability struct {
	PRNumber int      `json:"pr_number"`
	Title    string   `json:"title"`
	Score    float64  `json:"score"` // 0 to 100 (100 = Guaranteed Merge)
	Level    string   `json:"level"` // "High", "Medium", "Low", "Critical Risk"
	Reasons  []string `json:"reasons"`
}

// PredictMergeProbability calculates the statistical likelihood of an open PR being merged
func PredictMergeProbability(pr github.PullRequest) PRProbability {
	score := 50.0 // Base probability
	var reasons []string

	// 1. Check CI/CD Status (if available in payload/approximated)
	// We don't have direct CI status in standard PullRequest struct right now,
	// so we use surrogate metrics like review state and mergeable flag if present.
	// But let's assume `Additions` and `Deletions` are available.

	// 2. Code Churn (Size)
	totalChurn := pr.Additions + pr.Deletions
	if totalChurn == 0 {
		// Not loaded, use a fallback
		score -= 5.0
	} else if totalChurn < 50 {
		score += 30.0
		reasons = append(reasons, "Tiny code footprint (high merge chance)")
	} else if totalChurn < 300 {
		score += 15.0
		reasons = append(reasons, "Small, review-friendly code footprint")
	} else if totalChurn > 1500 {
		score -= 30.0
		reasons = append(reasons, "Massive code footprint (high abandonment risk)")
	} else {
		score -= 10.0
		reasons = append(reasons, "Large code footprint (moderate review friction)")
	}

	// 3. Stale Risk
	now := time.Now()
	age := now.Sub(pr.CreatedAt)
	daysOpen := age.Hours() / 24.0

	if daysOpen > 30 {
		score -= 40.0
		reasons = append(reasons, "Open for > 30 days (very high abandonment risk)")
	} else if daysOpen > 14 {
		score -= 20.0
		reasons = append(reasons, "Open for > 14 days (moderate abandonment risk)")
	} else if daysOpen < 2 {
		score += 10.0
		reasons = append(reasons, "Fresh PR (high engagement chance)")
	}

	// Bound the score
	score = math.Max(0.0, math.Min(100.0, score))

	// Categorize Level
	level := "Medium"
	if score >= 80 {
		level = "High"
	} else if score >= 50 {
		level = "Medium"
	} else if score >= 30 {
		level = "Low"
	} else {
		level = "Critical Risk"
	}

	if len(reasons) == 0 {
		reasons = append(reasons, "Standard PR profile")
	}

	return PRProbability{
		PRNumber: pr.Number,
		Title:    pr.Title,
		Score:    score,
		Level:    level,
		Reasons:  reasons,
	}
}

// PredictTopPRs returns a sorted list of PR merge probabilities
func PredictTopPRs(prs []github.PullRequest, limit int) []PRProbability {
	var probabilities []PRProbability

	for _, pr := range prs {
		if pr.State != "open" {
			continue
		}
		probabilities = append(probabilities, PredictMergeProbability(pr))
	}

	// Sort highest probability first
	for i := 0; i < len(probabilities)-1; i++ {
		for j := 0; j < len(probabilities)-i-1; j++ {
			if probabilities[j].Score < probabilities[j+1].Score {
				probabilities[j], probabilities[j+1] = probabilities[j+1], probabilities[j]
			}
		}
	}

	if len(probabilities) > limit {
		return probabilities[:limit]
	}
	return probabilities
}
