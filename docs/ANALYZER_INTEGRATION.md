# Real Analyzer Data Integration Guide

## Overview

This document explains how the analyzer data is integrated with the UI components through the new `AnalyzerDataBridge` and supporting modules.

## Architecture

```
GitHub API
    ↓
github.Client
    ↓
analyzer.{Functions}
    ↓
AnalysisResult
    ↓
AnalyzerDataBridge ← Transforms & Enriches Data
    ↓
Dashboard/Export Models
    ↓
UI Display
```

## Core Components

### 1. AnalyzerDataBridge (analyzer_bridge.go)

**Purpose**: Clean interface between analyzer logic and UI components.

**Key Methods**:
- `NewAnalyzerDataBridge(result)`: Creates bridge from AnalysisResult
- `GetHealthMetrics()`: Returns health-related metrics with colors and status
- `GetRepositoryInfo()`: Returns repository metadata (stars, forks, etc.)
- `GetContributorMetrics()`: Returns contributor analysis and diversity scores
- `GetCommitMetrics()`: Returns commit activity and trends
- `GetLanguageMetrics()`: Returns language composition and diversity
- `GetCompleteAnalysis()`: Returns all metrics combined
- `GenerateSummary()`: Creates human-readable analysis summary
- `GenerateRecommendations()`: Creates actionable recommendations

**Data Flow**:
```go
// From app.go
result := AnalysisResult{ /* analyzer output */ }
bridge := NewAnalyzerDataBridge(result)

// Use in dashboard
metricsMap := bridge.GetHealthMetrics()
// Returns: {health_score, health_status, health_color, ...}
```

### 2. ProgressTracker (progress.go)

**Purpose**: Displays multi-step analysis progress with visual feedback.

**Key Features**:
- 6-stage progress tracking (Fetch → Analyze → Compute → Complete)
- Visual progress bar with percentage
- Elapsed time tracking
- Stage completion indicators

**Usage**:
```go
tracker := NewProgressTracker()
// ... perform work ...
tracker.NextStage() // Move to next stage
```

**Display**:
```
⚙️  Fetching repository data
✅ Analyzing commits
⏳ Analyzing contributors
...
⏱️  2s elapsed
```

### 3. KeyboardShortcuts (shortcuts.go)

**Purpose**: Screen-specific keyboard shortcut definitions.

**Available Functions**:
- `GetMainMenuShortcuts()`: Menu navigation shortcuts
- `GetInputShortcuts()`: Input field shortcuts
- `GetDashboardShortcuts()`: Dashboard/export shortcuts
- `GetSettingsShortcuts()`: Settings navigation
- `GetHistoryShortcuts()`: History browser
- `GetHelpShortcuts()`: Help navigation
- `GetFileTreeShortcuts()`: File tree navigation
- `GetUniversalShortcuts()`: Shortcuts that work everywhere
- `GetShortcutsForScreen(name)`: Route to appropriate shortcut set
- `FormatShortcutsForDisplay()`: Format shortcuts as readable string

**Example**:
```go
shortcuts := GetShortcutsForScreen("dashboard")
displayText := FormatShortcutsForDisplay(shortcuts, terminalWidth)
```

### 4. ResponsiveLayout (responsive.go)

**Purpose**: Handle responsive design for different terminal sizes.

**Key Methods**:
- `IsSmallTerminal()`: Check if width < 80 or height < 24
- `IsMobileTerminal()`: Check if width < 60
- `GetMaxContentWidth()`: Safe content width with padding
- `GetMaxContentHeight()`: Safe content height with padding
- `CenterText()` / `CenterContent()`: Center content on screen
- `WrapText()`: Wrap long text to fit width
- `GetLayoutMode()`: Returns "mobile", "compact", "wide", or "default"
- `AdjustSpacing()`: Returns appropriate spacing for layout
- `RenderResponsiveBox()`: Adaptive box rendering

**Layout Modes**:
- **Mobile** (< 60 width): Minimal styling, single column
- **Compact** (< 80 width): Reduced spacing, essential content
- **Default** (80-120 width): Standard layout
- **Wide** (> 120 width): Multi-column, sidebar support

## Integration Points

### Dashboard Enhancement

**File**: `internal/ui/dashboard.go`

The dashboard now uses the bridge to display rich metrics:

```go
type DashboardModel struct {
    data     AnalysisResult
    bridge   *AnalyzerDataBridge    // NEW: Data transformation layer
    // ... other fields ...
}

func (m *DashboardModel) SetData(data AnalysisResult) {
    m.data = data
    m.bridge = NewAnalyzerDataBridge(data)  // Initialize bridge
}
```

### App State Management

**File**: `internal/ui/app.go`

The main app now tracks progress during analysis:

```go
type MainModel struct {
    progress  *ProgressTracker      // NEW: Track analysis stages
    // ... other fields ...
}

case stateLoading:
    if m.progress != nil {
        stages := m.progress.GetAllStages()
        // Display progress stages
        for _, stage := range stages {
            // Show ✅/⚙️/⏳ indicators
        }
    }
```

### Analysis Flow Enhancement

The `analyzeRepo` function now uses a progress tracker:

```go
func (m MainModel) analyzeRepo(repoName string) tea.Cmd {
    return func() tea.Msg {
        tracker := NewProgressTracker()
        
        // Stage 1: Fetch repo
        repo, _ := client.GetRepo(owner, name)
        tracker.NextStage()
        
        // Stage 2: Analyze commits
        commits, _ := client.GetCommits(...)
        tracker.NextStage()
        
        // ... continue through all stages ...
        
        return AnalysisResult{ /* metrics */ }
    }
}
```

## Data Transformation Examples

### Health Metrics Transform

**Input** (from analyzer):
```go
HealthScore: 75
BusFactor: 5
```

**Transform** (via bridge):
```go
bridge.GetHealthMetrics() → map{
    "health_score": 75,
    "health_status": "Good",
    "health_color": "yellow",
    "bus_factor": 5,
    "bus_risk": "Low",
    // ... more fields
}
```

### Contributor Analysis

**Input**:
```go
Contributors: [
    {Login: "alice", Contributions: 150},
    {Login: "bob", Contributions: 120},
    // ... more
]
```

**Transform**:
```go
bridge.GetContributorMetrics() → map{
    "total_contributors": 5,
    "top_contributors": [
        {Login: "alice", Contributions: 150, AvatarURL: "..."},
        // ... top 5
    ],
    "diversity_score": 65.4,
}
```

## New Features

### 1. Progress Display During Analysis

Shows real-time progress with emoji indicators:
- ⚙️ Currently processing
- ✅ Completed step
- ⏳ Pending step

### 2. Data Enrichment

The bridge adds computed metrics not available from raw analyzer output:
- Health status (Excellent/Good/Fair/Poor)
- Color coding for metrics
- Diversity scores
- Trend analysis
- Recommendations

### 3. Screen-Specific Shortcuts

Each screen displays relevant shortcuts in footer:
```
Dashboard: e: export • ↑ ↓ select • q: back
Settings: ↑ ↓ navigate • Enter: toggle • s: save
```

### 4. Responsive Layout

Terminal size detection with automatic adjustments:
- Mobile mode for small terminals
- Compact layout for medium terminals
- Full layout for standard terminals
- Wide layout for large terminals

## Usage Patterns

### Pattern 1: Get Metrics for Display

```go
// In dashboard view
bridge := m.bridge
metrics := bridge.GetHealthMetrics()

display := fmt.Sprintf("Health: %d %s", 
    metrics["health_score"],
    metrics["health_status"],
)
```

### Pattern 2: Generate Recommendations

```go
bridge := NewAnalyzerDataBridge(result)
recommendations := bridge.GenerateRecommendations()

for _, rec := range recommendations {
    println("💡 " + rec)
}
```

### Pattern 3: Track Progress

```go
tracker := NewProgressTracker()
// ... perform work ...
tracker.NextStage()

// Display
percentage := tracker.GetProgress() // (3, 6) = 50%
bar := tracker.GetProgressBar(30)   // Visual bar
```

### Pattern 4: Responsive Display

```go
layout := NewResponsiveLayout(width, height)

if layout.IsSmallTerminal() {
    // Use minimal display
} else {
    // Use full featured display
}

wrapped := layout.WrapText(longContent, padding)
```

## File Structure

```
internal/ui/
├── analyzer_bridge.go  (NEW - Data transformation)
├── progress.go        (NEW - Progress tracking)
├── shortcuts.go       (NEW - Keyboard shortcuts)
├── responsive.go      (NEW - Responsive design)
├── app.go            (MODIFIED - Integrated bridge)
├── dashboard.go      (MODIFIED - Uses bridge)
├── export.go         (Unchanged)
├── menu.go           (Unchanged)
├── styles.go         (Unchanged)
├── settings.go       (Unchanged)
├── history.go        (Unchanged)
├── help.go           (Unchanged)
├── tree.go           (Unchanged)
└── types.go          (Unchanged)
```

## Error Handling

The bridge handles:
- Empty contributor lists
- Missing language data
- Zero commits
- Null pointers

Returns safe default values for all scenarios.

## Testing the Integration

### Test 1: Basic Analysis Flow
```bash
# Run app and analyze a real repository
repo-lyzer
# Input: torvalds/linux
# Observe progress display and final dashboard with metrics
```

### Test 2: Terminal Resizing
```bash
# Start app, then resize terminal
# Verify layout adapts correctly
# Check for text truncation/wrapping
```

### Test 3: Keyboard Shortcuts
```bash
# On dashboard, press '?'
# Verify context-appropriate shortcuts display
# Test each shortcut (e: export, q: back, etc.)
```

### Test 4: Small Terminal
```bash
# Start app in 60x20 terminal
# Verify mobile-friendly layout
# Confirm all essential info is visible
```

## Future Enhancements

1. **Caching**: Cache bridge results to avoid recomputation
2. **Streaming**: Stream progress updates from analyzer
3. **Custom Metrics**: Allow users to define custom computed metrics
4. **Comparison**: Compare metrics across multiple repositories
5. **Trend Tracking**: Track metric changes over time
6. **Export Templates**: Custom export format templates

## Troubleshooting

### Issue: Progress doesn't update
- Ensure `ProgressTracker` is instantiated before analysis
- Check that `tracker.NextStage()` is called in sequence
- Verify state transitions in `app.go` Update()

### Issue: Metrics not displaying
- Confirm `NewAnalyzerDataBridge(result)` is called
- Check that analyzer functions return valid data
- Verify Dashboard.SetData() is called with AnalysisResult

### Issue: Text gets cut off on small terminals
- `ResponsiveLayout` should be initialized with actual window size
- Call `GetMaxContentWidth()` before rendering long text
- Use `WrapText()` for multi-line content

### Issue: Keyboard shortcuts not showing
- Ensure `GetShortcutsForScreen()` is called with correct screen name
- Verify `FormatShortcutsForDisplay()` is rendering to footer
- Check terminal width to ensure shortcuts fit

## Related Files

- [CLI_IMPROVEMENTS.md](./CLI_IMPROVEMENTS.md) - User-facing guide
- [DEVELOPER_GUIDE.md](./DEVELOPER_GUIDE.md) - Development reference
- [IMPLEMENTATION_SUMMARY.md](./IMPLEMENTATION_SUMMARY.md) - Overview
