package predictive

import (
	"math"
	"strings"

	"github.com/agnivo988/Repo-lyzer/internal/github"
)

// IssueDifficulty represents the estimated complexity of an issue
type IssueDifficulty struct {
	IssueNumber int      `json:"issue_number"`
	Title       string   `json:"title"`
	Score       float64  `json:"score"` // 0 to 100 (100 = Extremely Hard)
	Level       string   `json:"level"` // "Trivial", "Easy", "Medium", "Hard", "Expert"
	Reasons     []string `json:"reasons"`
}

var easyLabels = map[string]bool{
	"good first issue":  true,
	"good-first-issue":  true,
	"beginner-friendly": true,
	"easy":              true,
	"docs":              true,
	"documentation":     true,
	"typo":              true,
}

var hardLabels = map[string]bool{
	"architecture": true,
	"core":         true,
	"refactor":     true,
	"performance":  true,
	"breaking":     true,
	"security":     true,
}

// EstimateDifficulty calculates a deterministic difficulty score for an issue
func EstimateDifficulty(issue github.Issue) IssueDifficulty {
	score := 50.0 // Base score is Medium
	var reasons []string

	// 1. Analyze Labels
	hasEasyLabel := false
	hasHardLabel := false

	for _, label := range issue.Labels {
		lname := strings.ToLower(label.Name)
		if easyLabels[lname] {
			hasEasyLabel = true
			score -= 30.0
			reasons = append(reasons, "Contains beginner-friendly labels")
		}
		if hardLabels[lname] {
			hasHardLabel = true
			score += 30.0
			reasons = append(reasons, "Contains complex/architecture labels")
		}
	}

	// 2. Analyze Body Length and Complexity
	bodyLen := len(issue.Body)
	if bodyLen < 100 {
		// Too short, might be ambiguous/hard to understand without context
		score += 10.0
		reasons = append(reasons, "Extremely short description (ambiguity risk)")
	} else if bodyLen > 2000 {
		// Very long, likely a massive feature request or complex bug
		score += 20.0
		reasons = append(reasons, "Very large issue description (high complexity)")
	} else if !hasEasyLabel && !hasHardLabel {
		// Good length, well detailed
		score -= 5.0
		reasons = append(reasons, "Well-detailed issue length")
	}

	// 3. Presence of Code Blocks
	if strings.Contains(issue.Body, "```") {
		score -= 10.0
		reasons = append(reasons, "Contains code blocks (reproducible steps provided)")
	}

	// 4. Presence of Stack Traces (often complex bugs)
	if strings.Contains(strings.ToLower(issue.Body), "panic:") || strings.Contains(strings.ToLower(issue.Body), "exception") {
		score += 15.0
		reasons = append(reasons, "Contains stack traces (potential core bug)")
	}

	// Bound the score between 0 and 100
	score = math.Max(0.0, math.Min(100.0, score))

	// Categorize Level
	level := "Medium"
	if score <= 20 {
		level = "Trivial"
	} else if score <= 40 {
		level = "Easy"
	} else if score <= 60 {
		level = "Medium"
	} else if score <= 80 {
		level = "Hard"
	} else {
		level = "Expert"
	}

	// If no reasons exist, add a generic one
	if len(reasons) == 0 {
		reasons = append(reasons, "Standard issue with no extreme modifiers")
	}

	return IssueDifficulty{
		IssueNumber: issue.Number,
		Title:       issue.Title,
		Score:       score,
		Level:       level,
		Reasons:     reasons,
	}
}

// PredictBestFirstIssues returns a sorted slice of the easiest open issues
func PredictBestFirstIssues(issues []github.Issue, limit int) []IssueDifficulty {
	var difficulties []IssueDifficulty

	for _, issue := range issues {
		// Skip Pull Requests
		if issue.PullRequest != nil || issue.State != "open" {
			continue
		}
		difficulties = append(difficulties, EstimateDifficulty(issue))
	}

	// Simple bubble sort for difficulty (lowest score first)
	for i := 0; i < len(difficulties)-1; i++ {
		for j := 0; j < len(difficulties)-i-1; j++ {
			if difficulties[j].Score > difficulties[j+1].Score {
				difficulties[j], difficulties[j+1] = difficulties[j+1], difficulties[j]
			}
		}
	}

	if len(difficulties) > limit {
		return difficulties[:limit]
	}
	return difficulties
}
