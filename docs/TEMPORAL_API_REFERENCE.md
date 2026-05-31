# Temporal Intelligence API Reference

**Document Version**: 1.0  
**Last Updated**: May 18, 2026  
**Module**: `internal/temporal`, `internal/graph`, `internal/evolution`, `internal/predictive`, `internal/simulation`

---

## Graph Module API

### Package: `internal/graph`

Provides graph-based modeling of repository ecosystems.

#### Types

##### `NodeType` (int const)

Enumerates entity types in the temporal graph.

```go
const (
    NodeTypeContributor NodeType = iota
    NodeTypeFile
    NodeTypeSubsystem
    NodeTypeDependency
    NodeTypeIssue
)
```

##### `Node`

Represents an entity in the temporal repository graph.

```go
type Node struct {
    ID          string                 // Unique identifier
    Type        NodeType               // Entity type
    Timestamp   time.Time              // Introduction time
    Metadata    map[string]interface{} // Custom attributes
    Occurrences int                    // Occurrence count
}
```

**Methods**:
- `NewNode(id string, nodeType NodeType) *Node` - Creates new node
- `SetMetadata(key string, value interface{})` - Sets metadata
- `GetMetadata(key string) interface{}` - Gets metadata value

##### `EdgeType` (int const)

Enumerates relationship types in the temporal graph.

```go
const (
    EdgeTypeCollaboration EdgeType = iota
    EdgeTypeModification
    EdgeTypeDependency
    EdgeTypeIssueRelation
    EdgeTypeContainment
    EdgeTypeContainment
)
```

##### `Edge`

Represents a relationship between two nodes.

```go
type Edge struct {
    Source      *Node     // Starting node
    Target      *Node     // Ending node
    Type        EdgeType  // Relationship type
    Weight      float64   // Relationship strength [0, 1]
    Timestamp   time.Time // First observed
    Occurrences int       // Observation count
}
```

**Methods**:
- `NewEdge(source, target *Node, type EdgeType, weight float64) *Edge` - Creates new edge
- `IncreaseWeight(delta float64)` - Increases edge weight (capped at 1.0)
- `AverageWeight() float64` - Computes average weight

##### `Graph`

Main repository ecosystem model.

```go
type Graph struct {
    // Private fields managed by Graph methods
}
```

**Constructor**:
- `NewGraph() *Graph` - Creates empty graph

**Node Operations**:
- `AddNode(node *Node) error` - Adds node to graph
- `GetNode(id string) (*Node, error)` - Retrieves node by ID
- `QueryByType(nodeType NodeType) []*Node` - Gets all nodes of type
- `Query(predicate func(*Node) bool) []*Node` - Gets nodes matching predicate
- `GetAllNodes() []*Node` - Gets all nodes
- `NodeCount() int` - Returns node count

**Edge Operations**:
- `AddEdge(edge *Edge) error` - Adds edge (merges duplicates)
- `GetEdges(sourceID string) []*Edge` - Gets outgoing edges
- `IncomingEdges(nodeID string) []*Edge` - Gets incoming edges
- `GetAllEdges() []*Edge` - Gets all edges
- `EdgeCount() int` - Returns edge count

**Traversal**:
- `DFS(startID string, visitor TraversalVisitor) error` - Depth-first search
- `BFS(startID string, visitor TraversalVisitor) error` - Breadth-first search
- `ShortestPath(fromID, toID string) ([]*Node, error)` - Finds shortest path
- `AllPaths(fromID, toID string, maxDepth int) ([][]*Node, error)` - Finds all paths
- `DistanceTo(fromID, toID string) int` - Computes shortest distance
- `Neighborhood(nodeID string) []*Node` - Gets adjacent nodes
- `NodesByDistance(sourceID string) ([]*Node, []int, error)` - Nodes sorted by distance

**Metrics**:
- `DegreeCentrality(node *Node) float64` - Node degree centrality
- `BetweennessCentrality(node *Node) float64` - Node betweenness centrality
- `ClusteringCoefficient(node *Node) float64` - Node clustering coefficient
- `AverageDegreeCentrality() float64` - Average degree centrality
- `AverageClusteringCoefficient() float64` - Average clustering coefficient
- `Density() float64` - Graph density [0, 1]
- `ConnectedComponents() [][]Node` - Finds disconnected subgraphs
- `ComputeMetrics() GraphMetrics` - Computes all metrics

**Utilities**:
- `Clear()` - Removes all nodes and edges

---

## Temporal Module API

### Package: `internal/temporal`

Manages temporal repository states and timeline reconstruction.

#### Types

##### `TemporalEvent`

Represents a change event in the repository.

```go
type TemporalEvent struct {
    Timestamp   time.Time
    EventType   string                 // "commit", "contributor_added", etc.
    Contributor string
    Files       []string
    Details     map[string]interface{}
}
```

##### `RepositoryMetrics`

Metrics at a point in time.

```go
type RepositoryMetrics struct {
    CommitCount         int
    ContributorCount    int
    ActiveContributors  int
    AverageBusFactor    float64
    AverageHealth       int
    FilesChanged        int
    LinesAdded          int
    LinesRemoved        int
    AverageCommitFrequency float64
    DependencyCount     int
    IssuesOpen          int
    PullRequestsOpen    int
}
```

##### `Snapshot`

Repository state at a specific time.

```go
type Snapshot struct {
    Timestamp    time.Time
    CommitHash   string
    Graph        *graph.Graph
    Metrics      RepositoryMetrics
    Contributors []string
    Files        int
}
```

**Methods**:
- `NewSnapshot(timestamp time.Time, g *graph.Graph) *Snapshot` - Creates snapshot

##### `Timeline`

Evolution of repository over time.

```go
type Timeline struct {
    RepoName     string
    Owner        string
    Snapshots    []*Snapshot
    StartTime    time.Time
    EndTime      time.Time
    IntervalDays int
    Events       []TemporalEvent
}
```

**Methods**:
- `NewTimeline(owner, repoName string) *Timeline` - Creates timeline
- `AddSnapshot(snapshot *Snapshot) error` - Adds snapshot (must be chronological)
- `AddEvent(event TemporalEvent) error` - Adds temporal event
- `GetSnapshot(timestamp time.Time, exact bool) (*Snapshot, error)` - Retrieves snapshot
- `SnapshotsBetween(startTime, endTime time.Time) []*Snapshot` - Gets time-range snapshots
- `EventsBetween(startTime, endTime time.Time) []TemporalEvent` - Gets time-range events
- `SnapshotCount() int` - Returns snapshot count
- `EventCount() int` - Returns event count
- `Duration() time.Duration` - Time span of timeline
- `IsEmpty() bool` - Checks if empty
- `LatestSnapshot() *Snapshot` - Gets most recent snapshot
- `EarliestSnapshot() *Snapshot` - Gets earliest snapshot
- `SnapshotAtIndex(index int) (*Snapshot, error)` - Gets snapshot by index
- `WindowedSnapshots(windowDays int) [][]*Snapshot` - Groups snapshots by time window

##### `Coordinator`

Orchestrates complete temporal analysis workflow.

```go
type Coordinator struct {
    Timeline          *Timeline
    Detector          *evolution.Detector
    Predictor         *predictive.Predictor
    SimulationRunner  *simulation.ScenarioRunner
    // Results fields...
}
```

**Methods**:
- `NewCoordinator(owner, repoName string) *Coordinator` - Creates coordinator
- `ReconstructFromEvents(events []TemporalEvent) error` - Builds timeline
- `AnalyzeEvolution() error` - Detects patterns and drift
- `ForecastHealth(monthsAhead int) error` - Predicts health metrics
- `ForecastContributorRisks() error` - Predicts contributor risks
- `RunSimulation(scenario *simulation.SimulationScenario) (*simulation.SimulationResult, error)` - Executes scenario
- `Finalize() *AnalysisResult` - Generates final report
- `FullAnalysisPipeline(events []TemporalEvent, forecastMonths int) (*AnalysisResult, error)` - Complete workflow

---

## Evolution Module API

### Package: `internal/evolution`

Analyzes repository evolution patterns and risks.

#### Types

##### `EvolutionPattern`

Detected evolution pattern.

```go
type EvolutionPattern struct {
    Name        string
    Description string
    StartTime   time.Time
    EndTime     time.Time
    Indicators  map[string]float64
    Severity    string    // "low", "medium", "high"
    Confidence  float64   // [0, 1]
    Affected    []string
}
```

##### `DriftIndicator`

Architectural drift detection.

```go
type DriftIndicator struct {
    SubsystemID string
    MetricName  string
    Direction   string    // "increasing", "decreasing", "stable"
    Magnitude   float64
    StartValue  float64
    EndValue    float64
    TimeSpan    time.Duration
    Threshold   float64
    Severity    string
}
```

##### `RiskIndicator`

Detected risk in repository.

```go
type RiskIndicator struct {
    Category       string    // "complexity", "contributor", "dependency", etc.
    Name           string
    Severity       string    // "low", "medium", "high", "critical"
    Affected       []string
    Current        float64
    Threshold      float64
    Trajectory     string    // "improving", "stable", "worsening"
    Recommendations []string
}
```

##### `Detector`

Evolution pattern detector.

```go
type Detector struct {
    ComplexityThreshold float64
    DriftThreshold      float64
    RiskThreshold       float64
    // Results fields...
}
```

**Methods**:
- `NewDetector() *Detector` - Creates detector
- `DetectPatterns(timeline *temporal.Timeline) []EvolutionPattern` - Detects evolution patterns
- `DetectArchitecturalDrift(timeline *temporal.Timeline) []DriftIndicator` - Detects drift
- `AnalyzeComplexityGrowth(timeline *temporal.Timeline) ComplexityReport` - Analyzes complexity
- `TrackContributorEvolution(timeline *temporal.Timeline) []ContributorRole` - Tracks contributors
- `DetectKnowledgeSilos(timeline *temporal.Timeline) []Bottleneck` - Identifies knowledge silos
- `IdentifyRisks(timeline *temporal.Timeline) []RiskIndicator` - Identifies risks
- `ComputeRiskScore(timeline *temporal.Timeline) float64` - Computes overall risk

---

## Predictive Module API

### Package: `internal/predictive`

Provides forecasting and predictive modeling.

#### Types

##### `Prediction`

Single forecasted value.

```go
type Prediction struct {
    Timestamp  time.Time
    Value      float64
    LowerBound float64
    UpperBound float64
    Confidence float64
    Method     string
}
```

##### `ForecastResult`

Complete forecasting output.

```go
type ForecastResult struct {
    Metric         string
    Predictions    []Prediction
    Trend          string    // "improving", "stable", "degrading"
    RiskLevel      string    // "low", "medium", "high"
    Recommendations []string
    ConfidenceScore float64
    BaselineMean   float64
    BaselineStdDev float64
}
```

##### `Predictor`

Main forecasting engine.

```go
type Predictor struct {
    Models         map[string]PredictiveModel
    HistoricalData map[string][]float64
    ForecastHorizon int  // periods to forecast
    ConfidenceLevel float64 // [0.8, 0.99]
}
```

**Methods**:
- `NewPredictor() *Predictor` - Creates predictor
- `ForecastHealth(timeline *temporal.Timeline, months int) (*ForecastResult, error)` - Predicts health
- `ForecastMaturity(timeline *temporal.Timeline, months int) (*ForecastResult, error)` - Predicts maturity
- `ForecastContributorRisk(timeline *temporal.Timeline) ([]ContributorRiskForecast, error)` - Predicts contributor risks
- `EstimateBurnoutRisk(contributor string, timeline *temporal.Timeline) (float64, error)` - Estimates burnout
- `ForecastDependencyStability(timeline *temporal.Timeline, months int) (*ForecastResult, error)` - Predicts dependency stability
- `ProjectTechnicalDebt(timeline *temporal.Timeline, months int) (*ForecastResult, error)` - Projects technical debt

---

## Simulation Module API

### Package: `internal/simulation`

Provides what-if scenario simulation.

#### Types

##### `SimulationScenario`

What-if scenario definition.

```go
type SimulationScenario struct {
    Name         string
    Description  string
    ScenarioType string
    Parameters   map[string]interface{}
    Duration     time.Duration
    StartTime    time.Time
}
```

**Constructor**:
- `NewScenario(name, scenarioType string, duration time.Duration) *SimulationScenario`

##### `SimulationResult`

Simulation outcome.

```go
type SimulationResult struct {
    Scenario                 SimulationScenario
    InitialState             map[string]float64
    FinalState               map[string]float64
    HealthTrajectory         []float64
    RiskTrajectory           []float64
    ComplexityTrajectory     []float64
    Timestamps               []time.Time
    KeyFindings              []string
    Recommendations          []string
    HealthChange             float64
    RiskChange               float64
    Success                  bool
}
```

##### `ScenarioRunner`

Scenario execution engine.

```go
type ScenarioRunner struct {
    RepoName     string
    Owner        string
    TimestepDays int
    RandomSeed   int64
    Results      []SimulationResult
}
```

**Methods**:
- `NewScenarioRunner(owner, repoName string) *ScenarioRunner` - Creates runner
- `RunScenario(scenario SimulationScenario, timeline *temporal.Timeline) (*SimulationResult, error)` - Executes scenario
- `RunMultipleScenarios(scenarios []SimulationScenario, timeline *temporal.Timeline) ([]SimulationResult, error)` - Executes multiple scenarios
- `SimulateContributorDeparture(timeline *temporal.Timeline, contributor string, replacementMonths int) (*SimulationResult, error)` - Simulates departure
- `SimulateMajorRefactoring(timeline *temporal.Timeline, subsystem string, effortHours int) (*SimulationResult, error)` - Simulates refactoring
- `SimulateDependencyUpgrade(timeline *temporal.Timeline, dependency string, breakingChange bool) (*SimulationResult, error)` - Simulates upgrade
- `SimulateRapidGrowth(timeline *temporal.Timeline, subsystem string, growthRate float64, teamSize int) (*SimulationResult, error)` - Simulates growth
- `CompareScenarios(scenarios []SimulationScenario, timeline *temporal.Timeline) (string, error)` - Compares multiple scenarios

---

## Usage Examples

### Example 1: Basic Graph Construction

```go
// Create graph
g := graph.NewGraph()

// Add nodes
contributor := graph.NewNode("alice@github", graph.NodeTypeContributor)
file := graph.NewNode("main.go", graph.NodeTypeFile)

g.AddNode(contributor)
g.AddNode(file)

// Add edge
edge := graph.NewEdge(contributor, file, graph.EdgeTypeModification, 0.8)
g.AddEdge(edge)

// Query
fmt.Println("Nodes:", g.NodeCount())
fmt.Println("Edges:", g.EdgeCount())
fmt.Println("Density:", g.Density())
```

### Example 2: Timeline Creation

```go
timeline := temporal.NewTimeline("golang", "go")

// Create and add snapshots
snap1 := temporal.NewSnapshot(time.Now(), g)
snap1.Metrics.CommitCount = 50000
timeline.AddSnapshot(snap1)

// Query timeline
latestSnapshot := timeline.LatestSnapshot()
fmt.Printf("Latest snapshot: %v\n", latestSnapshot.Timestamp)
```

### Example 3: Evolution Analysis

```go
detector := evolution.NewDetector()
patterns := detector.DetectPatterns(timeline)
risks := detector.IdentifyRisks(timeline)

for _, risk := range risks {
    fmt.Printf("%s: %s\n", risk.Name, risk.Severity)
}
```

### Example 4: Health Forecasting

```go
predictor := predictive.NewPredictor()
forecast, err := predictor.ForecastHealth(timeline, 6)

fmt.Printf("Trend: %s\n", forecast.Trend)
fmt.Printf("Risk Level: %s\n", forecast.RiskLevel)
```

### Example 5: Scenario Simulation

```go
runner := simulation.NewScenarioRunner("golang", "go")
result, err := runner.SimulateContributorDeparture(timeline, "alice", 6)

fmt.Printf("Health Change: %+.2f\n", result.HealthChange)
fmt.Printf("Findings: %v\n", result.KeyFindings)
```

### Example 6: Complete Pipeline

```go
coordinator := temporal.NewCoordinator("kubernetes", "kubernetes")

result, err := coordinator.FullAnalysisPipeline(events, 6)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Health Score: %d\n", result.HealthScore)
fmt.Printf("Risk Level: %s\n", result.OverallRiskLevel)
fmt.Printf("Critical Issues: %v\n", result.CriticalIssues)
```

---

## Error Handling

Operations that can fail return `error` as the last return value (idiomatic Go). Always check for errors:

```go
coordinator := temporal.NewCoordinator("owner", "repo")
err := coordinator.ReconstructFromEvents(events)
if err != nil {
    fmt.Fprintf(os.Stderr, "Reconstruction failed: %v\n", err)
    return
}
```

---

## Thread Safety

The `Graph` type uses mutexes for thread safety. All operations are safe for concurrent access. Other types (`Timeline`, `Coordinator`, etc.) are not inherently thread-safe; synchronize access at the application level if needed.

---

## Performance Notes

- **Graph operations**: Most are O(V + E) or O(V²) for metric computations
- **Traversal**: DFS/BFS are O(V + E)
- **Shortest path**: Uses BFS, O(V + E)
- **Timeline queries**: O(N) where N is number of snapshots
- **Prediction**: O(M) where M is historical data points

For very large graphs (>100K nodes), consider:
- Streaming graph construction
- Windowed analysis for long timelines
- Batch processing of scenarios
