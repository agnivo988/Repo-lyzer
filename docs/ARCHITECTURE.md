# Repo-lyzer Architecture Guide

This document provides a comprehensive overview of the Repo-lyzer architecture, designed to help new contributors quickly understand the codebase and start contributing effectively.

## Table of Contents

1. [High-Level Overview](#high-level-overview)
2. [Directory Structure](#directory-structure)
3. [Core Components](#core-components)
4. [Data Flow](#data-flow)
5. [Adding New Features](#adding-new-features)
6. [Key Patterns](#key-patterns)

---

## High-Level Overview

Repo-lyzer follows a **modular architecture** with clear separation of concerns:

```
┌─────────────────────────────────────────────────────────────────┐
│                         main.go                                  │
│                    (Application Entry)                           │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                          cmd/                                    │
│              (CLI Commands & Menu System)                        │
│         root.go │ menu.go │ analyze.go │ compare.go             │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                       internal/ui/                               │
│              (Interactive Terminal UI - Bubble Tea)              │
│    app.go │ dashboard.go │ menu.go │ themes.go │ tree.go        │
└─────────────────────────────────────────────────────────────────┘
                              │
              ┌───────────────┼───────────────┐
              ▼               ▼               ▼
┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐
│ internal/github │ │internal/analyzer│ │ internal/output │
│   (API Client)  │ │ (Computations)  │ │  (Formatting)   │
└─────────────────┘ └─────────────────┘ └─────────────────┘
```

---

## Directory Structure

### Root Level

```
repo-lyzer/
├── main.go              # Application entry point
├── go.mod               # Go module definition
├── go.sum               # Dependency checksums
├── readme.md            # Project README
├── CONTRIBUTING.md      # Contribution guidelines
└── CODE_OF_CONDUCT.md   # Community guidelines
```

### `/cmd` - CLI Commands

The command layer handles CLI initialization and routing.

| File | Purpose |
|------|---------|
| `root.go` | Cobra root command setup |
| `menu.go` | Launches the interactive TUI menu |
| `analyze.go` | CLI command for single repo analysis |
| `compare.go` | CLI command for comparing two repos |

**When to modify:** Add new CLI commands or flags here.

### `/internal/github` - GitHub API Client

Handles all communication with the GitHub REST API.

| File | Purpose |
|------|---------|
| `client.go` | HTTP client with authentication |
| `repo.go` | Repository metadata (stars, forks, etc.) |
| `commits.go` | Commit history fetching |
| `contributor.go` | Contributor statistics |
| `languages.go` | Language breakdown |
| `issues.go` | Issue tracking data |
| `tree.go` | File tree structure |
| `rate_limit.go` | API rate limit handling |

**Key struct:** `Client` - Central API client with token support

```go
type Client struct {
    http  *http.Client
    token string  // From GITHUB_TOKEN env var
}
```

**When to modify:** Add new GitHub API endpoints or data fetching.

### `/internal/analyzer` - Analysis & Metrics

Contains all metric computation logic.

| File | Purpose | Output |
|------|---------|--------|
| `health.go` | Repository health score | 0-100 score |
| `bus_factor.go` | Contributor concentration risk | Factor + risk level |
| `maturity.go` | Project maturity assessment | Score + level (Emerging/Growing/Mature/Established) |
| `commit_activity.go` | Commit frequency analysis | Daily/weekly activity data |
| `recruiter_summary.go` | Quick evaluation summary | Formatted summary string |
| `security.go` | Security-related checks | Security indicators |
| `files.go` | File structure analysis | File statistics |

**When to modify:** Add new metrics or analysis algorithms.

### `/internal/ui` - Terminal User Interface

The interactive TUI built with Bubble Tea framework.

| File | Purpose |
|------|---------|
| `app.go` | **Main application state machine** - handles all state transitions |
| `menu.go` | Main menu component |
| `dashboard.go` | Analysis results dashboard with multiple views |
| `tree.go` | File tree viewer component |
| `file_edit.go` | File viewing and editing interface |
| `history.go` | Analysis history management |
| `export.go` | JSON/Markdown export functionality |
| `themes.go` | Color theme system (7 themes) |
| `styles.go` | Lipgloss style definitions |
| `charts.go` | Commit activity charts |
| `progress.go` | Loading progress indicators |
| `shortcuts.go` | Keyboard shortcut definitions |
| `types.go` | Shared type definitions |
| `responsive.go` | Terminal size handling |
| `analyzer_bridge.go` | Bridge between UI and analyzers |

**Key file:** `app.go` - Contains the main state machine:

```go
const (
    stateMenu sessionState = iota
    stateInput
    stateLoading
    stateDashboard
    stateTree
    stateFileEdit
    stateSettings
    stateHelp
    stateHistory
    stateCompareInput
    stateCompareLoading
    stateCompareResult
    stateCloneInput
    stateCloning
)
```

**When to modify:** Add new UI screens, views, or interactions.

### `/internal/output` - Output Formatting

Handles non-TUI output formatting (for CLI mode).

| File | Purpose |
|------|---------|
| `tables.go` | Table formatting |
| `charts.go` | ASCII chart generation |
| `health.go` | Health score display |
| `recruiter.go` | Recruiter summary formatting |
| `json.go` | JSON output |
| `styles.go` | CLI color styles |

### `/docs` - Documentation

Project documentation files.

### `/exports` - Export Output

Default location for exported analysis files and history.

---

## Data Flow

### 1. Application Startup

```
main.go
    └── cmd.RunMenu()
            └── ui.NewMainModel()
                    └── Initialize all UI components
                            └── Display main menu
```

### 2. Single Repository Analysis

```
User selects "Analyze Repository"
    │
    ▼
stateInput: User enters "owner/repo"
    │
    ▼
stateLoading: analyzeRepo() is called
    │
    ├── github.NewClient()
    │       └── Creates authenticated HTTP client
    │
    ├── client.GetRepo(owner, repo)
    │       └── Fetches: stars, forks, issues, dates
    │
    ├── client.GetCommits(owner, repo, 365)
    │       └── Fetches: last year of commits
    │
    ├── client.GetContributors(owner, repo)
    │       └── Fetches: contributor list with commit counts
    │
    ├── client.GetLanguages(owner, repo)
    │       └── Fetches: language byte counts
    │
    ├── client.GetFileTree(owner, repo, branch)
    │       └── Fetches: repository file structure
    │
    ├── analyzer.CalculateHealth(repo, commits)
    │       └── Computes: health score (0-100)
    │
    ├── analyzer.BusFactor(contributors)
    │       └── Computes: bus factor + risk level
    │
    └── analyzer.RepoMaturityScore(repo, commits, contributors)
            └── Computes: maturity score + level
    │
    ▼
AnalysisResult struct is created
    │
    ▼
stateDashboard: Results displayed in interactive dashboard
    │
    ├── View 1: Overview (health, bus factor, maturity)
    ├── View 2: Repository Details (stars, forks, issues)
    ├── View 3: Languages (breakdown with percentages)
    ├── View 4: Activity (commit chart)
    ├── View 5: Contributors (top contributors list)
    ├── View 6: Recruiter Summary
    └── View 7: API Status
```

### 3. Repository Comparison

```
User selects "Compare Repositories"
    │
    ▼
stateCompareInput: User enters two repos
    │
    ▼
stateCompareLoading: Both repos analyzed in parallel
    │
    ▼
stateCompareResult: Side-by-side comparison displayed
```

---

## Adding New Features

### Adding a New Analyzer

1. **Create the analyzer file** in `/internal/analyzer/`:

```go
// internal/analyzer/my_metric.go
package analyzer

func CalculateMyMetric(repo *github.Repo, commits []github.Commit) int {
    // Your computation logic
    return score
}
```

2. **Add to AnalysisResult** in `/internal/ui/types.go`:

```go
type AnalysisResult struct {
    // ... existing fields
    MyMetricScore int
    MyMetricLevel string
}
```

3. **Call from analyzeRepo()** in `/internal/ui/app.go`:

```go
myScore, myLevel := analyzer.CalculateMyMetric(repo, commits)
return AnalysisResult{
    // ... existing fields
    MyMetricScore: myScore,
    MyMetricLevel: myLevel,
}
```

4. **Display in dashboard** in `/internal/ui/dashboard.go`:

```go
func (m DashboardModel) myMetricView() string {
    // Format and return the view
}
```

### Adding a New GitHub API Endpoint

1. **Add the data type** in `/internal/github/`:

```go
// internal/github/my_data.go
package github

type MyData struct {
    Field1 string `json:"field1"`
    Field2 int    `json:"field2"`
}

func (c *Client) GetMyData(owner, repo string) (*MyData, error) {
    var data MyData
    err := c.get("https://api.github.com/repos/"+owner+"/"+repo+"/myendpoint", &data)
    return &data, err
}
```

### Adding a New UI View

1. **Add state constant** in `/internal/ui/app.go`:

```go
const (
    // ... existing states
    stateMyNewView
)
```

2. **Handle state in Update()** method
3. **Add view rendering in View()** method
4. **Create view function**:

```go
func (m MainModel) myNewView() string {
    // Build and return the view
}
```

### Adding a New Theme

Add to `/internal/ui/themes.go`:

```go
var MyTheme = Theme{
    Name:       "My Theme",
    Primary:    lipgloss.Color("#XXXXXX"),
    Secondary:  lipgloss.Color("#XXXXXX"),
    // ... other colors
}

// Add to AvailableThemes slice
var AvailableThemes = []Theme{
    // ... existing themes
    MyTheme,
}
```

---

## Key Patterns

### Bubble Tea Model Pattern

All UI components follow the Bubble Tea pattern:

```go
type MyModel struct {
    // State fields
}

func (m MyModel) Init() tea.Cmd { return nil }

func (m MyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Handle keyboard input
    case tea.WindowSizeMsg:
        // Handle resize
    }
    return m, nil
}

func (m MyModel) View() string {
    // Return rendered view
}
```

### State Machine Pattern

The main app uses a state machine for navigation:

```go
switch m.state {
case stateMenu:
    // Handle menu state
case stateInput:
    // Handle input state
case stateDashboard:
    // Handle dashboard state
}
```

### Error Handling Pattern

Errors are passed through the message system:

```go
func (m MainModel) analyzeRepo(name string) tea.Cmd {
    return func() tea.Msg {
        // ... do work
        if err != nil {
            return err  // Error is a valid message
        }
        return AnalysisResult{...}
    }
}
```

---

## Quick Reference

| Task | Location |
|------|----------|
| Add CLI command | `/cmd/` |
| Add GitHub API call | `/internal/github/` |
| Add new metric | `/internal/analyzer/` |
| Add UI component | `/internal/ui/` |
| Add theme | `/internal/ui/themes.go` |
| Modify dashboard | `/internal/ui/dashboard.go` |
| Change styles | `/internal/ui/styles.go` |
| Add export format | `/internal/ui/export.go` |

---

## Getting Help

- Check existing code for patterns
- Read the [Bubble Tea documentation](https://github.com/charmbracelet/bubbletea)
- Review [Lipgloss styling guide](https://github.com/charmbracelet/lipgloss)
- Open an issue for questions

Happy contributing! 🚀
