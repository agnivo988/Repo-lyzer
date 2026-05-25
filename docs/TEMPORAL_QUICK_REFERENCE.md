# Temporal Intelligence Quick Reference Guide

**Version**: 1.0  
**Last Updated**: May 18, 2026  
**For**: Developers & Users

---

## For Users: Getting Started

### Installation

```bash
# Clone and build with temporal features
git clone https://github.com/agnivo988/Repo-lyzer.git
cd Repo-lyzer
go mod tidy
go install
```

### Basic Commands

```bash
# Analyze repository evolution
repo-lyzer temporal analyze golang/go

# Forecast repository health (6 months)
repo-lyzer temporal forecast kubernetes/kubernetes

# Analyze contributor networks
repo-lyzer temporal contributors python/cpython

# Detect architectural drift
repo-lyzer temporal drift tensorflow/tensorflow

# Simulate scenarios
repo-lyzer temporal simulate rust/rust key_contributor_departure
```

### Understanding Output

**Health Forecast**:
- Trend: "improving", "stable", or "degrading"
- Confidence: 0-1 (higher = more reliable)
- Risk Level: "low", "medium", or "high"

**Evolution Patterns**:
- Severity: Importance of the pattern
- Affected: What subsystems/contributors are involved

**Risk Indicators**:
- Category: Type of risk (complexity, contributor, dependency, etc.)
- Severity: "low", "medium", "high", "critical"
- Trajectory: Direction of the risk (improving/stable/worsening)

---

## For Developers: Key Concepts

### Architecture Layers

```
┌─────────────────────────────────────┐
│         CLI Commands                 │  (cmd/temporal.go)
├─────────────────────────────────────┤
│    Temporal Coordinator              │  (temporal/coordinator.go)
│  (Orchestrates all modules)          │
├─────────────────────────────────────┤
│ Graph │ Temporal │ Evolution │ ...   │  (internal/graph/, temporal/, etc.)
│ (Data) │ (Time)   │(Analysis)│(Pred.)│
├─────────────────────────────────────┤
│         Existing Repo-lyzer          │  (github/, analyzer/, output/)
└─────────────────────────────────────┘
```

### Core Types

**Graph** (internal/graph/):
- `Node`: Represents contributors, files, subsystems
- `Edge`: Represents relationships (collaboration, modification, etc.)
- `Graph`: Container with query and traversal operations

**Temporal** (internal/temporal/):
- `Timeline`: Sequence of snapshots over time
- `Snapshot`: Repository state at a point in time
- `TemporalEvent`: A change event in the repository

**Evolution** (internal/evolution/):
- `EvolutionPattern`: Detected trend or pattern
- `DriftIndicator`: Architectural drift detection
- `RiskIndicator`: Identified risks

**Predictive** (internal/predictive/):
- `Prediction`: Single forecasted value with confidence
- `ForecastResult`: Complete forecast with trend
- `PredictiveModel`: Interface for forecasting models

**Simulation** (internal/simulation/):
- `SimulationScenario`: What-if scenario definition
- `SimulationResult`: Outcome of a simulation

### Complete Analysis Pipeline

```
1. GitHub API
   ↓ Fetch commits, contributors, files
   ↓ Create TemporalEvents
   
2. Temporal Reconstruction
   ↓ Create Timeline with Snapshots
   ↓ Build Graph at each point in time
   
3. Evolution Analysis
   ↓ Detect patterns
   ↓ Identify drift
   ↓ Detect risks
   
4. Predictive Forecasting
   ↓ Forecast health
   ↓ Forecast maturity
   ↓ Forecast contributor risks
   
5. Simulation (Optional)
   ↓ Run scenarios
   ↓ Analyze outcomes
   
6. Output Formatting
   ↓ Generate reports
   ↓ Display results
```

---

## Module Quick Reference

### Graph Module

```go
import "github.com/agnivo988/Repo-lyzer/internal/graph"

// Create graph
g := graph.NewGraph()

// Create nodes
contributor := graph.NewNode("alice@github", graph.NodeTypeContributor)
file := graph.NewNode("main.go", graph.NodeTypeFile)

// Add nodes
g.AddNode(contributor)
g.AddNode(file)

// Create and add edge
edge := graph.NewEdge(contributor, file, graph.EdgeTypeModification, 0.8)
g.AddEdge(edge)

// Query operations
nodes := g.Query(func(n *Node) bool { return n.Type == graph.NodeTypeFile })
neighbors := g.Neighborhood(contributor.ID)

// Traversal
g.DFS(contributor.ID, func(n *Node) error { 
    fmt.Println(n.ID)
    return nil
})

// Metrics
degree := g.DegreeCentrality(contributor)
clustering := g.ClusteringCoefficient(contributor)
density := g.Density()
```

### Temporal Module

```go
import "github.com/agnivo988/Repo-lyzer/internal/temporal"

// Create timeline
timeline := temporal.NewTimeline("golang", "go")

// Create snapshot
snap := temporal.NewSnapshot(time.Now(), g)
snap.Metrics.CommitCount = 50000

// Add to timeline
timeline.AddSnapshot(snap)

// Query operations
latest := timeline.LatestSnapshot()
earliest := timeline.EarliestSnapshot()
windows := timeline.WindowedSnapshots(7) // 7-day windows
```

### Evolution Module

```go
import "github.com/agnivo988/Repo-lyzer/internal/evolution"

// Create detector
detector := evolution.NewDetector()

// Analyze
patterns := detector.DetectPatterns(timeline)
drift := detector.DetectArchitecturalDrift(timeline)
risks := detector.IdentifyRisks(timeline)
complexity := detector.AnalyzeComplexityGrowth(timeline)

for _, risk := range risks {
    fmt.Printf("%s: %s (%s)\n", risk.Name, risk.Severity, risk.Trajectory)
}
```

### Predictive Module

```go
import "github.com/agnivo988/Repo-lyzer/internal/predictive"

// Create predictor
predictor := predictive.NewPredictor()

// Generate forecasts
healthForecast, err := predictor.ForecastHealth(timeline, 6)
contributorRisks, err := predictor.ForecastContributorRisk(timeline)

fmt.Printf("Health Trend: %s\n", healthForecast.Trend)
fmt.Printf("Predictions: %v\n", healthForecast.Predictions)
```

### Simulation Module

```go
import "github.com/agnivo988/Repo-lyzer/internal/simulation"

// Create runner
runner := simulation.NewScenarioRunner("golang", "go")

// Run predefined scenarios
result, err := runner.SimulateContributorDeparture(timeline, "alice", 6)
fmt.Printf("Health Change: %+.2f\n", result.HealthChange)

// Run custom scenario
scenario := simulation.NewScenario("Custom", "custom_type", 90*24*time.Hour)
result, err = runner.RunScenario(*scenario, timeline)
```

### Coordinator (Orchestration)

```go
import "github.com/agnivo988/Repo-lyzer/internal/temporal"

// Create coordinator
coordinator := temporal.NewCoordinator("golang", "go")

// Run complete pipeline
result, err := coordinator.FullAnalysisPipeline(events, 6)
if err != nil {
    log.Fatal(err)
}

// Access results
fmt.Printf("Health Score: %d\n", result.HealthScore)
fmt.Printf("Risk Level: %s\n", result.OverallRiskLevel)
fmt.Printf("Critical Issues: %v\n", result.CriticalIssues)
```

---

## Common Patterns

### Pattern 1: Analyze Repository

```go
coordinator := temporal.NewCoordinator(owner, repo)
events := fetchTemporalEvents(owner, repo)
result, _ := coordinator.FullAnalysisPipeline(events, 6)
displayAnalysisReport(result)
```

### Pattern 2: Run Specific Analysis

```go
coordinator := temporal.NewCoordinator(owner, repo)
coordinator.ReconstructFromEvents(events)
coordinator.AnalyzeEvolution()
patterns := coordinator.Detector.DetectPatterns(coordinator.Timeline)
```

### Pattern 3: Forecast Only

```go
coordinator := temporal.NewCoordinator(owner, repo)
coordinator.ReconstructFromEvents(events)
forecast, _ := coordinator.Predictor.ForecastHealth(coordinator.Timeline, 12)
displayForecast(forecast)
```

### Pattern 4: Run Simulation

```go
coordinator := temporal.NewCoordinator(owner, repo)
coordinator.ReconstructFromEvents(events)
scenario := simulation.NewScenario("Test", "contributor_departure", 90*24*time.Hour)
result, _ := coordinator.RunSimulation(scenario)
analyzeSimulationResult(result)
```

---

## Performance Guidelines

### Recommended Input Sizes

- **Commits**: 100 - 100,000+ (comfortable range)
- **Contributors**: 1 - 1,000+
- **Snapshots**: 10 - 1,000 (auto-generated)
- **Forecast Horizon**: 3 - 24 months

### Expected Performance

- **Reconstruction**: O(C × F) where C = commits, F = files
- **Graph metrics**: O(V²) for dense graphs
- **Pattern detection**: O(snapshots × metrics)
- **Forecasting**: O(historical_data_points)

### Memory Usage

- Typical repo (10K commits): 50-100 MB
- Large repo (100K commits): 200-400 MB
- Very large (500K+ commits): 500+ MB

### Timeouts

- Reconstruction: <2 minutes for 100K commits
- Analysis: <1 minute
- Forecasting: <30 seconds
- Simulation: <1 minute per scenario

---

## Troubleshooting

### Issue: Analysis Takes Too Long
**Solution**: Use windowed analysis or reduce time period

### Issue: Out of Memory
**Solution**: 
- Reduce snapshot count
- Use streaming mode
- Analyze subset of commits

### Issue: Low Prediction Confidence
**Solution**:
- Ensure at least 100 commits in history
- Use longer historical period (6+ months)
- Check for consistent commit patterns

### Issue: Pattern Detection Finds Nothing
**Solution**:
- Repository may be too new/stable
- Adjust detection thresholds
- Check for sufficient data

---

## File Structure

```
Repo-lyzer/
├── cmd/
│   └── temporal.go                 # CLI commands
├── internal/
│   ├── graph/
│   │   ├── types.go               # Node, Edge, Graph types
│   │   ├── graph.go               # Graph implementation
│   │   ├── metrics.go             # Centrality metrics
│   │   └── traversal.go           # DFS, BFS, shortest path
│   ├── temporal/
│   │   ├── types.go               # TemporalEvent, Metrics types
│   │   ├── snapshot.go            # Snapshot, Timeline
│   │   ├── coordinator.go         # Orchestration
│   │   └── conversion.go          # (To implement)
│   ├── evolution/
│   │   ├── types.go               # Evolution types
│   │   └── detector.go            # Pattern detection
│   ├── predictive/
│   │   ├── types.go               # Prediction types
│   │   └── forecaster.go          # Forecasting
│   └── simulation/
│       ├── types.go               # Scenario types
│       └── engine.go              # Simulation
└── docs/
    ├── TEMPORAL_INTELLIGENCE_FEATURE_SPEC.md
    ├── TEMPORAL_ARCHITECTURE.md
    ├── TEMPORAL_API_REFERENCE.md
    ├── TEMPORAL_INTEGRATION_GUIDE.md
    ├── TEMPORAL_MVP_ROADMAP.md
    └── TEMPORAL_QUICK_REFERENCE.md (this file)
```

---

## Documentation Index

1. **Feature Specification**: What temporal analysis does
2. **Architecture Design**: How the system works
3. **API Reference**: Function signatures and usage
4. **Integration Guide**: How to integrate with existing code
5. **MVP Roadmap**: Implementation plan
6. **Quick Reference**: This file (quick lookup)

---

## Key Metrics

### Repository Health
- **0-33**: High Risk
- **34-66**: Medium Risk
- **67-100**: Low Risk (Healthy)

### Bus Factor
- **1**: High Risk (single point of failure)
- **2**: Medium Risk
- **3**: Low Risk (distributed knowledge)

### Complexity
- **Low**: <30% of max
- **Medium**: 30-70% of max
- **High**: >70% of max

### Contributor Roles
- **Core**: 30%+ of commits
- **Active**: 5-30% of commits
- **Occasional**: 1-5% of commits
- **Inactive**: <1% recent commits

---

## Glossary

| Term | Definition |
|------|-----------|
| **Temporal Event** | A change event in repository (commit, issue, etc.) |
| **Snapshot** | Repository state at a point in time |
| **Timeline** | Sequence of snapshots over time |
| **Graph** | Network representation of repository entities |
| **Centrality** | Importance/influence of a node in the graph |
| **Drift** | Deviation from baseline in subsystem metrics |
| **Pattern** | Recurring trend or behavior detected over time |
| **Forecast** | Prediction of future metric values |
| **Scenario** | What-if simulation setup |
| **Bottleneck** | Person/area with concentrated knowledge |
| **Bus Factor** | How many people must leave to damage project |

---

## Resources

### Internal Documentation
- See `docs/` directory for complete documentation
- CLI help: `repo-lyzer temporal --help`
- Command-specific help: `repo-lyzer temporal analyze --help`

### Code Examples
- See `docs/TEMPORAL_API_REFERENCE.md` for code samples
- See test files (`*_test.go`) for usage examples

### Contact
- Open issues on GitHub
- Check existing issues for solutions
- Request features in discussions

---

## Summary

**Temporal Intelligence** transforms Repo-lyzer from static analysis into **predictive intelligence**:

✓ **Understand**: How repositories evolve over time  
✓ **Predict**: Future health, risks, and sustainability  
✓ **Simulate**: Impact of scenarios (departures, refactoring, etc.)  
✓ **Act**: Based on data-driven insights  

**Start with**: `repo-lyzer temporal analyze <owner>/<repo>`
