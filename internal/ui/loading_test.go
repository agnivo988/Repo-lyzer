package ui

import (
	"strings"
	"testing"
)

func TestLoadingModelView_RendersProgressAndProgressBar(t *testing.T) {
	m := NewLoadingModel()
	tracker := NewProgressTracker()
	m.SetProgress(tracker)
	m.SetRepoName("owner/repo")

	// Advance to complete 1 stage of progress
	tracker.NextStage()

	view := m.View(80, 20)

	// Since 1 of 9 stages is completed, the percentage should be 11%
	expectedPercentage := "11%"
	if !strings.Contains(view, expectedPercentage) {
		t.Errorf("expected view to contain %q, but did not. Got:\n%s", expectedPercentage, view)
	}

	expectedProgressLine := "Progress: 11%"
	if !strings.Contains(view, expectedProgressLine) {
		t.Errorf("expected view to contain %q, but did not. Got:\n%s", expectedProgressLine, view)
	}

	// The progress bar should be present in the view (it will have standard/shimmer blocks)
	if !strings.Contains(view, "[") || !strings.Contains(view, "]") {
		t.Errorf("expected view to contain progress bar brackets, but did not. Got:\n%s", view)
	}
}
