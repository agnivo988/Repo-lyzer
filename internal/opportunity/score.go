package opportunity

import "github.com/google/go-github/v55/github"

type IssueScore struct {
    Issue   *github.Issue
    Weight  int
    AgeDays int
}

// calculateWeight returns a base weight for an issue based on its classification.
func calculateWeight(issue *github.Issue) int {
    switch Classify(issue) {
    case ClassificationGoodFirstIssue:
        return 3
    case ClassificationHelpWanted:
        return 2
    case ClassificationBug:
        return 1
    case ClassificationEnhancement:
        return 1
    default:
        return 0
    }
}

// computeScore aggregates issue scores and contributor count into a final contribution score.
func computeScore(issues []*github.Issue, contributors [] *github.ContributorStats) int {
    total := 0
    for _, iss := range issues {
        weight := calculateWeight(iss)
        // simple age penalty: newer issues get a small bonus (max 2 points)
        ageDays := 0
        if iss.CreatedAt != nil {
            ageDays = int(time.Since(*iss.CreatedAt).Hours() / 24)
        }
        bonus := 0
        if ageDays < 30 {
            bonus = 2
        } else if ageDays < 90 {
            bonus = 1
        }
        total += weight + bonus
    }
    // Add a point per contributor (capped at 50 for sanity)
    contribPoints := len(contributors)
    if contribPoints > 50 {
        contribPoints = 50
    }
    total += contribPoints
    return total
}
