package monitor

import (
	"fmt"
	"testing"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/github"
)

// filterIssues returns only true issues (PullRequest == nil) from a mixed list.
func filterIssues(items []github.Issue) int {
	count := 0
	for _, item := range items {
		if item.PullRequest == nil {
			count++
		}
	}
	return count
}

// filterPRs returns only pull requests (PullRequest != nil) from a mixed list.
func filterPRs(items []github.Issue) int {
	count := 0
	for _, item := range items {
		if item.PullRequest != nil {
			count++
		}
	}
	return count
}

// TestFilterIssues_ExcludesPRs verifies that pull requests are excluded from the
// issue count when filtering by the PullRequest field.
func TestFilterIssues_ExcludesPRs(t *testing.T) {
	prMarker := &struct{}{}

	tests := []struct {
		name      string
		items     []github.Issue
		wantIssue int
		wantPR    int
	}{
		{
			name:      "empty list",
			items:     []github.Issue{},
			wantIssue: 0,
			wantPR:    0,
		},
		{
			name: "all issues",
			items: []github.Issue{
				{Number: 1, Title: "Bug report", PullRequest: nil},
				{Number: 2, Title: "Feature request", PullRequest: nil},
				{Number: 3, Title: "Docs update", PullRequest: nil},
			},
			wantIssue: 3,
			wantPR:    0,
		},
		{
			name: "all pull requests",
			items: []github.Issue{
				{Number: 10, Title: "Fix typo", PullRequest: prMarker},
				{Number: 11, Title: "Add feature", PullRequest: prMarker},
			},
			wantIssue: 0,
			wantPR:    2,
		},
		{
			name: "mixed issues and PRs",
			items: []github.Issue{
				{Number: 1, Title: "Bug", PullRequest: nil},
				{Number: 2, Title: "PR fix", PullRequest: prMarker},
				{Number: 3, Title: "Feature", PullRequest: nil},
				{Number: 4, Title: "PR feature", PullRequest: prMarker},
				{Number: 5, Title: "Question", PullRequest: nil},
			},
			wantIssue: 3,
			wantPR:    2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotIssue := filterIssues(tc.items)
			gotPR := filterPRs(tc.items)

			if gotIssue != tc.wantIssue {
				t.Errorf("filterIssues() = %d, want %d", gotIssue, tc.wantIssue)
			}
			if gotPR != tc.wantPR {
				t.Errorf("filterPRs() = %d, want %d", gotPR, tc.wantPR)
			}

			// The total should always equal the original length
			if gotIssue+gotPR != len(tc.items) {
				t.Errorf("filterIssues() + filterPRs() = %d, want %d (total items)",
					gotIssue+gotPR, len(tc.items))
			}
		})
	}
}

// TestNotification_OnlyOnStateChange verifies that notifications are only sent
// when the count changes, preventing notification flooding.
func TestNotification_OnlyOnStateChange(t *testing.T) {
	notifications := make(chan Notification, 100)

	state := &MonitorState{
		LastIssueCount: 3,
		LastPRCount:    2,
	}

	prMarker := &struct{}{}
	items := []github.Issue{
		{Number: 1, Title: "Bug", PullRequest: nil},
		{Number: 2, Title: "PR", PullRequest: prMarker},
		{Number: 3, Title: "Feature", PullRequest: nil},
		{Number: 4, Title: "PR2", PullRequest: prMarker},
		{Number: 5, Title: "Task", PullRequest: nil},
	}

	// Simulate checkIssues logic
	issueCount := filterIssues(items)
	if issueCount != state.LastIssueCount {
		notifications <- Notification{
			Type:      "issue",
			Title:     "Issues Update",
			Message:   fmt.Sprintf("Repository has %d open issues", issueCount),
			Timestamp: time.Now(),
			Severity:  "info",
		}
		state.LastIssueCount = issueCount
	}

	// Simulate checkPullRequests logic
	prCount := filterPRs(items)
	if prCount != state.LastPRCount {
		notifications <- Notification{
			Type:      "pr",
			Title:     "Pull Requests Update",
			Message:   fmt.Sprintf("Repository has %d open pull requests", prCount),
			Timestamp: time.Now(),
			Severity:  "info",
		}
		state.LastPRCount = prCount
	}

	// Verify state was updated correctly
	if state.LastIssueCount != 3 {
		t.Errorf("LastIssueCount = %d, want 3", state.LastIssueCount)
	}
	if state.LastPRCount != 2 {
		t.Errorf("LastPRCount = %d, want 2", state.LastPRCount)
	}

	// No notifications should have been sent because the counts match
	if len(notifications) != 0 {
		t.Errorf("Expected 0 notifications (no state change), got %d", len(notifications))
	}

	// Now simulate a state change (new issue added)
	items = append(items, github.Issue{Number: 6, Title: "New Bug", PullRequest: nil})
	issueCount = filterIssues(items)
	if issueCount != state.LastIssueCount {
		notifications <- Notification{
			Type:      "issue",
			Title:     "Issues Update",
			Message:   fmt.Sprintf("Repository has %d open issues", issueCount),
			Timestamp: time.Now(),
			Severity:  "info",
		}
		state.LastIssueCount = issueCount
	}

	if len(notifications) != 1 {
		t.Errorf("Expected 1 notification after state change, got %d", len(notifications))
	}

	n := <-notifications
	if n.Type != "issue" {
		t.Errorf("Notification type = %q, want %q", n.Type, "issue")
	}
	if state.LastIssueCount != 4 {
		t.Errorf("LastIssueCount after update = %d, want 4", state.LastIssueCount)
	}
}
