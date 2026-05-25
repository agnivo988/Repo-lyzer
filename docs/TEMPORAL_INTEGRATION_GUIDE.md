# Temporal Intelligence Integration Guide

**Document Version**: 1.0  
**Date**: May 18, 2026  
**Purpose**: Guide for integrating temporal modules with existing Repo-lyzer systems

---

## 1. Overview

The Temporal Intelligence system integrates with existing Repo-lyzer components while maintaining separation of concerns. This guide explains integration points and best practices.

### Architecture Layers

```
CLI Layer (cmd/temporal.go)
    ↓
Orchestration (temporal/coordinator.go)
    ↓
Analysis Modules (graph/, evolution/, predictive/, simulation/)
    ↓
Existing Repo-lyzer (github/, analyzer/, output/)
```

---

## 2. GitHub API Integration

### New GitHub Functions

The `internal/github/` package needs these new functions to support temporal analysis:

#### 2.1 Commit History Retrieval

```go
package github

// FetchCommitHistory retrieves the complete commit history for a repository
// This is essential for temporal reconstruction
func (c *Client) FetchCommitHistory(owner, repo string) ([]Commit, error) {
    // Implementation should:
    // 1. Paginate through all commits
    // 2. Extract timestamp, author, files changed
    // 3. Handle rate limiting gracefully
    // 4. Return commits in chronological order
}
```

**Usage in Temporal**:
```go
commits, err := client.FetchCommitHistory("golang", "go")
// Then convert to TemporalEvents for reconstruction
```

#### 2.2 Contributor Timeline

```go
// FetchContributorTimeline retrieves when contributors joined and their activity patterns
func (c *Client) FetchContributorTimeline(owner, repo string) ([]ContributorEvent, error) {
    // Implementation should:
    // 1. Track first commit timestamp per contributor
    // 2. Detect activity patterns (active, inactive)
    // 3. Record role changes if detectable
}
```

#### 2.3 File History

```go
// FetchFileHistory tracks changes to a specific file over time
func (c *Client) FetchFileHistory(owner, repo, path string) ([]FileEvent, error) {
    // Implementation should:
    // 1. Retrieve commits affecting the file
    // 2. Track complexity metrics over time
    // 3. Identify major refactorings
}
```

### Implementation Strategy

```go
// File: internal/github/temporal.go (new file)
package github

import (
    "context"
    "time"
)

// ContributorEvent tracks when a contributor acts
type ContributorEvent struct {
    Contributor string
    Timestamp   time.Time
    CommitCount int
    EventType   string // "first_commit", "active", "inactive"
}

// FileEvent tracks changes to a file
type FileEvent struct {
    FilePath      string
    Timestamp     time.Time
    LinesAdded    int
    LinesRemoved  int
    Author        string
    CommitMessage string
}

// TODO: Implement these functions using the existing GitHub API client
```

---

## 3. Analyzer Module Integration

### Enhanced Health Scoring

Integrate temporal analysis with existing health scoring:

```go
// File: internal/analyzer/temporal_health.go (new file)
package analyzer

import (
    "github.com/agnivo988/Repo-lyzer/internal/temporal"
)

// TemporalHealthReport extends health analysis with temporal insights
type TemporalHealthReport struct {
    CurrentHealth     int       // Existing health score
    HealthTrend       string    // "improving", "stable", "degrading"
    ProjectedHealth   int       // Forecast 6 months ahead
    HealthVolatility  float64   // How stable the health metric is
    CriticalRisks     []string  // Emerging risks
    Recommendations   []string  // Actions to improve health
}

// CalculateTemporalHealth augments existing health score with temporal analysis
func CalculateTemporalHealth(
    repo *github.Repo,
    timeline *temporal.Timeline,
    events []temporal.TemporalEvent,
) TemporalHealthReport {
    // 1. Get existing health score
    currentHealth := CalculateHealth(repo, []github.Commit{}) // Existing function
    
    // 2. Reconstruct temporal data
    coordinator := temporal.NewCoordinator(repo.Owner, repo.Name)
    if err := coordinator.ReconstructFromEvents(events); err != nil {
        // handle error appropriately
        return TemporalHealthReport{}
    }

    // 3. Analyze evolution
    if err := coordinator.AnalyzeEvolution(); err != nil {
        return TemporalHealthReport{}
    }
    if err := coordinator.ForecastHealth(6); err != nil {
        return TemporalHealthReport{}
    }

    // 4. Create combined report
    return TemporalHealthReport{
        CurrentHealth:    currentHealth,
        HealthTrend:      "stable", // From forecast
        ProjectedHealth:  75,        // From forecast
        HealthVolatility: 0.15,     // From metrics
    }
}
```

### Extend Existing Analyzers

```go
// Enhance bus_factor.go
func CalculateTemporalBusFactor(
    timeline *temporal.Timeline,
) (score int, trend string) {
    // Existing function computes current bus factor
    // Temporal version tracks how bus factor changes over time
    // and predicts risks of contributor departure
}

// Enhance maturity.go
func CalculateTemporalMaturity(
    timeline *temporal.Timeline,
) (score int, trajectory string) {
    // Tracks maturity over time instead of just snapshot
}
```

---

## 4. Output Module Integration

### Temporal Report Generation

Create new output formatters for temporal analysis:

```go
// File: internal/output/temporal.go (new file)
package output

import (
    "github.com/agnivo988/Repo-lyzer/internal/temporal"
)

// OutputTemporalJSON formats temporal analysis as JSON
func OutputTemporalJSON(result *temporal.AnalysisResult) (string, error) {
    // Marshal result to JSON with proper indentation
    // Include graphs, patterns, forecasts, recommendations
}

// OutputTemporalMarkdown creates readable markdown report
func OutputTemporalMarkdown(result *temporal.AnalysisResult) (string, error) {
    // Create markdown report with sections:
    // - Executive Summary
    // - Evolution Patterns
    // - Predictions
    // - Risks & Recommendations
    // - Simulation Results (if run)
}

// OutputTemporalChart renders timeline charts
func RenderTemporalChart(predictions []temporal.Prediction) string {
    // Use existing chart functionality to visualize:
    // - Health trajectory
    // - Risk progression
    // - Contributor trends
}
```

### Integration with Existing Output

```go
// In internal/output/json.go, add handler:
case "temporal":
    data, err := OutputTemporalJSON(temporalResult)
    // Handle output
    
// In internal/output/styles.go, add styling for temporal elements
```

---

## 5. CLI Integration

### Command Registration

The temporal commands are registered in `cmd/temporal.go`:

```go
// Root temporal command
repo-lyzer temporal analyze <owner>/<repo>
repo-lyzer temporal forecast <owner>/<repo>
repo-lyzer temporal contributors <owner>/<repo>
repo-lyzer temporal drift <owner>/<repo>
repo-lyzer temporal simulate <owner>/<repo> <scenario>
```

### Execution Flow

```
1. Parse repository URL
   ↓
2. Fetch repository data via github.Client
   ↓
3. Create temporal.Coordinator
   ↓
4. Reconstruct timeline from commit history
   ↓
5. Run analysis (evolution, prediction, etc.)
   ↓
6. Format and display results
```

### Example Integration

```go
// In cmd/temporal.go, analyzeTemporalCmd.RunE:
func(cmd *cobra.Command, args []string) error {
    repoURL := args[0]
    owner, repo := parseRepoURL(repoURL)
    
    // 1. Get GitHub client
    client := github.NewClient(token)
    
    // 2. Fetch repository data
    ghRepo, err := client.GetRepo(owner, repo)
    commits, err := client.FetchCommitHistory(owner, repo)
    
    // 3. Create temporal coordinator
    coordinator := temporal.NewCoordinator(owner, repo)
    
    // 4. Convert commits to events and reconstruct
    events := convertCommitsToEvents(commits)
    if err := coordinator.ReconstructFromEvents(events); err != nil {
        return err
    }

    // 5. Run analysis pipeline
    result, err := coordinator.FullAnalysisPipeline(events, 6)
    if err != nil {
        return err
    }

    // 6. Format and display
    out, err := output.OutputTemporalMarkdown(result)
    if err != nil {
        return err
    }
    fmt.Println(out)

    return nil
}
```

---

## 6. Data Flow Integration

### From GitHub API to Temporal Analysis

```
GitHub API
  ↓ (FetchCommitHistory)
  ↓ (FetchContributorTimeline)
  ↓ (FetchFileHistory)
Raw GitHub Data
  ↓ (convertToTemporalEvents)
TemporalEvent[]
  ↓ (ReconstructFromEvents)
Timeline with Snapshots
  ↓ (AnalyzeEvolution)
  ↓ (ForecastHealth)
  ↓ (ForecastContributorRisks)
Analysis Results
  ↓ (FormatOutput)
Display/Export
```

### Conversion Functions

```go
// File: internal/temporal/conversion.go (new file)
package temporal

import (
    "github.com/agnivo988/Repo-lyzer/internal/github"
)

// CommitsToEvents converts GitHub commits to temporal events
func CommitsToEvents(commits []github.Commit) []TemporalEvent {
    events := make([]TemporalEvent, 0, len(commits))
    for _, commit := range commits {
        event := TemporalEvent{
            Timestamp:   commit.CommitDate,
            EventType:   "commit",
            Contributor: commit.Author.Login,
            Files:       commit.ChangedFiles,
            Details: map[string]interface{}{
                "message":       commit.Message,
                "additions":     commit.Additions,
                "deletions":     commit.Deletions,
                "commit_hash":   commit.SHA,
            },
        }
        events = append(events, event)
    }
    return events
}
```

---

## 7. Testing Integration

### Unit Tests

```go
// File: internal/temporal/coordinator_test.go
func TestCoordinatorReconstruction(t *testing.T) {
    // Test complete reconstruction workflow
}

func TestCoordinatorAnalysisPipeline(t *testing.T) {
    // Test full analysis pipeline
}

// File: internal/graph/graph_test.go
func TestGraphAddNode(t *testing.T) {
    // Test node operations
}

func TestGraphTraversal(t *testing.T) {
    // Test DFS, BFS, shortest path
}
```

### Integration Tests

```go
// Test with real repositories or fixtures
// 1. Small repo for quick validation
// 2. Medium repo for performance baseline
// 3. Complex repo for pattern detection
```

---

## 8. Performance Considerations

### Integration Impact

- **Graph construction**: O(C × F) where C = commits, F = files per commit
  - For large repos: Consider streaming or windowing
  
- **API calls**: Multiple calls for history fetch
  - Implement rate limiting and caching
  - Use pagination efficiently

- **Memory usage**: Timeline snapshots can be large
  - Consider lazy loading for old snapshots
  - Implement garbage collection triggers

### Optimization Strategies

```go
// 1. Caching layer
type TimedCache struct {
    data      map[string]interface{}
    expiry    time.Time
}

// 2. Streaming graph construction
func StreamGraphConstruction(commitChan chan github.Commit) (*graph.Graph, error) {
    g := graph.NewGraph()
    for commit := range commitChan {
        // Process incrementally
    }
    return g, nil
}

// 3. Windowed analysis
func WindowedAnalysis(timeline *Timeline, windowDays int) {
    windows := timeline.WindowedSnapshots(windowDays)
    for _, window := range windows {
        // Analyze each window independently
    }
}
```

---

## 9. Error Handling

### Error Propagation

```go
// Maintain error chain for debugging
if err := coordinator.ReconstructFromEvents(events); err != nil {
    return fmt.Errorf("timeline reconstruction failed: %w", err)
}

if err := coordinator.AnalyzeEvolution(); err != nil {
    return fmt.Errorf("evolution analysis failed: %w", err)
}
```

### Graceful Degradation

```go
// If temporal analysis fails, provide partial results
coordinator.ReconstructFromEvents(events) // Might fail partially
patterns := coordinator.EvolutionPatterns    // Use what we have
risks := coordinator.RiskIndicators         // Might be incomplete
```

---

## 10. Migration Guide

### Existing Code Updates Needed

#### 1. Update `cmd/root.go` or `cmd/menu.go`

Add temporal commands to menu or help text.

#### 2. Update `go.mod` (if external dependencies needed)

Currently, no additional dependencies are required beyond what Repo-lyzer already uses.

#### 3. Update `internal/github/client.go`

Add new methods for commit history and contributor timeline fetching.

#### 4. Create Output Handlers

Add JSON and Markdown formatters for temporal results in `internal/output/`.

### Backward Compatibility

- All changes are additive (new commands, new modules)
- Existing commands remain unchanged
- No modifications to existing analyzers' function signatures
- Output formats for existing commands unaffected

---

## 11. Configuration

### Temporal Analysis Configuration

```go
// internal/config/temporal.go (new file, optional)
type TemporalConfig struct {
    // Analysis parameters
    MaxSnapshotCount     int           // Limit snapshots for memory
    SnapshotInterval     int           // Days between snapshots
    ForecastHorizon      int           // Months to forecast
    ConfidenceLevel      float64       // 0.95 for 95% confidence
    
    // Pattern detection thresholds
    DriftThreshold       float64       // 0.5
    ComplexityThreshold  float64       // 0.7
    RiskThreshold        float64       // 0.6
    
    // Performance
    EnableCaching        bool          // Cache results
    MaxGraphSize         int           // Max nodes in graph
    StreamingMode        bool          // For large repos
}
```

---

## 12. Documentation

### For Contributors

Developers adding new analysis types should:

1. Extend `evolution/detector.go` for pattern detection
2. Extend `predictive/forecaster.go` for new predictions
3. Extend `simulation/engine.go` for new scenarios
4. Add corresponding CLI commands in `cmd/temporal.go`
5. Add tests in respective `_test.go` files
6. Update this integration guide

### For Users

Users should refer to:
- `docs/TEMPORAL_INTELLIGENCE_FEATURE_SPEC.md` for feature overview
- `docs/TEMPORAL_ARCHITECTURE.md` for system design
- `docs/TEMPORAL_API_REFERENCE.md` for API details
- CLI help: `repo-lyzer temporal --help`

---

## 13. Troubleshooting

### Common Issues

**Issue**: Graph construction is slow
- **Solution**: Use streaming mode or windowed analysis

**Issue**: Out of memory on large repos
- **Solution**: Reduce snapshot count or implement garbage collection

**Issue**: Inconsistent temporal events
- **Solution**: Validate event timestamps and sort before processing

**Issue**: Prediction confidence is low
- **Solution**: Ensure sufficient historical data (minimum 100 commits recommended)

---

## 14. Future Extensions

### Phase 2 Integration Points

1. **GraphRAG Integration**
   - Extend output with LLM-generated insights
   - Add to `internal/output/temporal_insights.go`

2. **Repository Evolution Replay**
   - Integrate with UI for visualization
   - Add to `cmd/temporal.go` with new `replay` command

3. **Advanced Modeling**
   - Machine learning models in `internal/predictive/ml_models.go`
   - Simulation improvements in `internal/simulation/advanced_scenarios.go`

---

## 15. Maintenance Notes

- Keep temporal modules modular and independent
- Maintain separation between data reconstruction and analysis
- Use interfaces for extensibility (e.g., PredictiveModel interface)
- Document assumptions and limitations in each module
- Regular performance profiling on large repositories
