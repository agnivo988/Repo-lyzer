package ui

import (
	"testing"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/config"
	"github.com/agnivo988/Repo-lyzer/internal/github"
)

func TestAnalysisDeduplicationAndDebounce(t *testing.T) {
	// Initialize with defaults
	m := NewMainModel(nil, &config.AppSettings{})
	m.state = stateInput

	// 1. Test normal submission
	msg := AnalyzeRepoMsg{repoName: "owner/repo"}
	newModel, cmd := m.Update(msg)
	m = newModel.(MainModel)

	if !m.analysisInProgress {
		t.Error("Expected analysisInProgress to be true after first submission")
	}
	if m.state != stateLoading {
		t.Errorf("Expected state to be stateLoading, got %v", m.state)
	}
	if cmd == nil {
		t.Error("Expected a command to be returned for the first submission")
	}

	// 2. Test immediate duplicate submission (Debounce/Guard)
	newModel2, cmd2 := m.Update(msg)
	if cmd2 != nil {
		t.Error("Expected duplicate submission to be ignored (cmd should be nil)")
	}
	if !newModel2.(MainModel).analysisInProgress {
		t.Error("analysisInProgress should remain true")
	}

	// 3. Test submission after completion
	// Simulate success
	mModel, _ := m.Update(AnalysisResult{
		Repo: &github.Repo{FullName: "owner/repo"},
	})
	m = mModel.(MainModel)
	if m.analysisInProgress {
		t.Error("Expected analysisInProgress to be false after completion")
	}

	// 4. Test debounce with state transition
	m.state = stateInput
	m.lastSubmitTime = time.Now() // Just submitted
	mModel, cmd = m.Update(msg)
	m = mModel.(MainModel)
	if cmd != nil {
		t.Error("Expected submission within 2s to be ignored by debounce")
	}

	// 5. Test submission after debounce interval
	m.lastSubmitTime = time.Now().Add(-3 * time.Second)
	mModel, cmd = m.Update(msg)
	m = mModel.(MainModel)
	if cmd == nil {
		t.Error("Expected submission after debounce interval to be accepted")
	}
}

func TestCompareDeduplication(t *testing.T) {
	m := NewMainModel(nil, &config.AppSettings{})
	m.state = stateCompareInput

	msg := CompareReposMsg{Repo1: "owner/repo1", Repo2: "owner/repo2"}
	newModel, cmd := m.Update(msg)
	m = newModel.(MainModel)

	if !m.analysisInProgress {
		t.Error("Expected analysisInProgress to be true after compare request")
	}
	if cmd == nil {
		t.Error("Expected a command to be returned for the first compare request")
	}

	// Duplicate
	_, cmd2 := m.Update(msg)
	if cmd2 != nil {
		t.Error("Expected duplicate compare request to be ignored")
	}
}
