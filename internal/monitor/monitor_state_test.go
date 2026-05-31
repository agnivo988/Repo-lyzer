package monitor

import (
	"os"
	"testing"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/cache"
)

func setTempHome(t *testing.T) {
	t.Helper()

	tempHome := t.TempDir()
	originalHome := os.Getenv("HOME")
	originalUserProfile := os.Getenv("USERPROFILE")

	if err := os.Setenv("HOME", tempHome); err != nil {
		t.Fatalf("failed to set HOME: %v", err)
	}
	if err := os.Setenv("USERPROFILE", tempHome); err != nil {
		t.Fatalf("failed to set USERPROFILE: %v", err)
	}

	t.Cleanup(func() {
		_ = os.Setenv("HOME", originalHome)
		_ = os.Setenv("USERPROFILE", originalUserProfile)
	})
}

// TestMonitorState_SaveLoadRoundTrip verifies monitor state is persisted and restored end to end.
func TestMonitorState_SaveLoadRoundTrip(t *testing.T) {
	setTempHome(t)

	c, err := cache.NewCache()
	if err != nil {
		t.Fatalf("NewCache() error = %v", err)
	}

	timestamp := time.Unix(1716076800, 0).UTC()
	m := &Monitor{
		cache: c,
		owner: "octocat",
		repo:  "hello-world",
		state: &MonitorState{
			Owner:                "octocat",
			Repo:                 "hello-world",
			LastCommitSHA:        "abcdef123456",
			LastIssueCount:       4,
			LastPRCount:          2,
			LastContributorCount: 7,
			LastUpdated:          timestamp,
		},
	}

	m.saveState()

	reloadedCache, err := cache.NewCache()
	if err != nil {
		t.Fatalf("NewCache() reload error = %v", err)
	}

	restored := &Monitor{
		cache: reloadedCache,
		owner: "octocat",
		repo:  "hello-world",
		state: &MonitorState{Owner: "octocat", Repo: "hello-world"},
	}
	restored.loadState()

	if restored.state.LastCommitSHA != "abcdef123456" {
		t.Fatalf("LastCommitSHA = %q, want %q", restored.state.LastCommitSHA, "abcdef123456")
	}
	if restored.state.LastIssueCount != 4 {
		t.Fatalf("LastIssueCount = %d, want 4", restored.state.LastIssueCount)
	}
	if restored.state.LastPRCount != 2 {
		t.Fatalf("LastPRCount = %d, want 2", restored.state.LastPRCount)
	}
	if restored.state.LastContributorCount != 7 {
		t.Fatalf("LastContributorCount = %d, want 7", restored.state.LastContributorCount)
	}
	if !restored.state.LastUpdated.Equal(timestamp) {
		t.Fatalf("LastUpdated = %v, want %v", restored.state.LastUpdated, timestamp)
	}
}

// TestMonitorState_LoadState_InvalidCachePayload verifies invalid cached payloads do not overwrite state.
func TestMonitorState_LoadState_InvalidCachePayload(t *testing.T) {
	setTempHome(t)

	c, err := cache.NewCache()
	if err != nil {
		t.Fatalf("NewCache() error = %v", err)
	}

	if err := c.Set("octocat/hello-world", "invalid-state"); err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	m := &Monitor{
		cache: c,
		owner: "octocat",
		repo:  "hello-world",
		state: &MonitorState{
			Owner:         "octocat",
			Repo:          "hello-world",
			LastCommitSHA: "existing-sha",
		},
	}

	m.loadState()

	if m.state.LastCommitSHA != "existing-sha" {
		t.Fatalf("LastCommitSHA was unexpectedly overwritten: %q", m.state.LastCommitSHA)
	}
}

// TestMonitorState_LoadState_FillsMissingIdentity verifies missing owner and repo fields are restored.
func TestMonitorState_LoadState_FillsMissingIdentity(t *testing.T) {
	setTempHome(t)

	c, err := cache.NewCache()
	if err != nil {
		t.Fatalf("NewCache() error = %v", err)
	}

	if err := c.Set("octocat/hello-world", MonitorState{LastCommitSHA: "xyz"}); err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	m := &Monitor{
		cache: c,
		owner: "octocat",
		repo:  "hello-world",
		state: &MonitorState{},
	}

	m.loadState()

	if m.state.Owner != "octocat" {
		t.Fatalf("Owner = %q, want %q", m.state.Owner, "octocat")
	}
	if m.state.Repo != "hello-world" {
		t.Fatalf("Repo = %q, want %q", m.state.Repo, "hello-world")
	}
	if m.state.LastCommitSHA != "xyz" {
		t.Fatalf("LastCommitSHA = %q, want %q", m.state.LastCommitSHA, "xyz")
	}
}
