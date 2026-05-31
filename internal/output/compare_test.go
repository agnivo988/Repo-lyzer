package output

import (
	"strings"
	"testing"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/github"
)

func TestBuildCompareReport(t *testing.T) {
	report := BuildCompareReport(
		CompareInput{
			Repo: &github.Repo{
				FullName:      "owner/alpha",
				Description:   "Alpha repo",
				Stars:         120,
				Forks:         18,
				OpenIssues:    7,
				Language:      "Go",
				DefaultBranch: "main",
				PushedAt:      time.Now().Add(-24 * time.Hour),
			},
			Commits: []github.Commit{{}, {}, {}},
			Contributors: []github.Contributor{{}, {}},
			Languages: map[string]int{"Go": 1000, "Shell": 100},
		},
		CompareInput{
			Repo: &github.Repo{
				FullName:      "owner/beta",
				Description:   "Beta repo",
				Stars:         80,
				Forks:         12,
				OpenIssues:    10,
				Language:      "TypeScript",
				DefaultBranch: "main",
				PushedAt:      time.Now().Add(-48 * time.Hour),
			},
			Commits: []github.Commit{{}, {}},
			Contributors: []github.Contributor{{}},
			Languages: map[string]int{"TypeScript": 900, "CSS": 50},
		},
	)

	if report.Repo1.FullName != "owner/alpha" || report.Repo2.FullName != "owner/beta" {
		t.Fatalf("unexpected repo names: %#v", report)
	}

	if report.Verdict == "" {
		t.Fatal("expected a non-empty verdict")
	}

	if report.Repo1.HealthScore == 0 || report.Repo2.HealthScore == 0 {
		t.Fatal("expected health scores to be computed")
	}
}

func TestRenderCompareOutputs(t *testing.T) {
	report := BuildCompareReport(
		CompareInput{
			Repo: &github.Repo{FullName: "owner/alpha", Stars: 10, Forks: 4, OpenIssues: 1, Language: "Go", DefaultBranch: "main"},
			Commits: []github.Commit{{}},
			Contributors: []github.Contributor{{}},
			Languages: map[string]int{"Go": 1},
		},
		CompareInput{
			Repo: &github.Repo{FullName: "owner/beta", Stars: 20, Forks: 6, OpenIssues: 2, Language: "Rust", DefaultBranch: "main"},
			Commits: []github.Commit{{}, {}},
			Contributors: []github.Contributor{{}, {}},
			Languages: map[string]int{"Rust": 2},
		},
	)

	terminal := RenderCompareTerminal(report)
	if !strings.Contains(terminal, "Repository Comparison") {
		t.Fatalf("terminal output missing header: %s", terminal)
	}
	lowerTerm := strings.ToLower(terminal)
	normTerm := strings.ReplaceAll(lowerTerm, " / ", "/")
	if !strings.Contains(normTerm, "owner/alpha") || !strings.Contains(normTerm, "owner/beta") {
		t.Fatalf("terminal output missing one of repo names: %s", terminal)
	}
	if !strings.Contains(normTerm, strings.ToLower(report.Verdict)) {
		t.Fatalf("terminal output missing verdict: %s", terminal)
	}

	jsonData, err := RenderCompareJSON(report)
	if err != nil {
		t.Fatalf("RenderCompareJSON failed: %v", err)
	}
	lowerJSON := strings.ToLower(string(jsonData))
	normJSON := strings.ReplaceAll(lowerJSON, " / ", "/")
	if !strings.Contains(normJSON, "owner/alpha") || !strings.Contains(normJSON, "owner/beta") {
		t.Fatalf("JSON output missing one of repo names: %s", string(jsonData))
	}
	if !strings.Contains(normJSON, strings.ToLower(report.Verdict)) {
		t.Fatalf("JSON output missing verdict: %s", string(jsonData))
	}

	markdown, err := RenderCompareMarkdown(report)
	if err != nil {
		t.Fatalf("RenderCompareMarkdown failed: %v", err)
	}
	if !strings.Contains(string(markdown), "# Repository Comparison") {
		t.Fatalf("markdown output missing title: %s", string(markdown))
	}
	lowerMD := strings.ToLower(string(markdown))
	normMD := strings.ReplaceAll(lowerMD, " / ", "/")
	if !strings.Contains(normMD, "owner/alpha") || !strings.Contains(normMD, "owner/beta") {
		t.Fatalf("markdown output missing one of repo names: %s", string(markdown))
	}
	if !strings.Contains(normMD, strings.ToLower(report.Verdict)) {
		t.Fatalf("markdown output missing verdict: %s", string(markdown))
	}

	htmlData, err := RenderCompareHTML(report)
	if err != nil {
		t.Fatalf("RenderCompareHTML failed: %v", err)
	}
	if !strings.Contains(string(htmlData), "Repo-lyzer comparison report") {
		t.Fatalf("html output missing title: %s", string(htmlData))
	}
	lowerHTML := strings.ToLower(string(htmlData))
	normHTML := strings.ReplaceAll(lowerHTML, " / ", "/")
	if !strings.Contains(normHTML, "owner/alpha") || !strings.Contains(normHTML, "owner/beta") {
		t.Fatalf("html output missing one of repo names: %s", string(htmlData))
	}
	if !strings.Contains(normHTML, strings.ToLower(report.Verdict)) {
		t.Fatalf("html output missing verdict: %s", string(htmlData))
	}
}