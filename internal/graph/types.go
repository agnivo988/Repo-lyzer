// Package graph provides graph-based modeling of repositories as evolving ecosystems.
// It models entities (contributors, files, subsystems) as nodes and relationships
// (collaboration, modification, dependency) as edges.
//
// The graph structure enables efficient querying, traversal, and computation of
// network metrics for temporal repository analysis.
package graph

import "time"

// NodeType enumerates different entity types in the temporal repository graph.
type NodeType int

const (
	// NodeTypeContributor represents a person who contributes to the repository
	NodeTypeContributor NodeType = iota
	// NodeTypeFile represents a source code file in the repository
	NodeTypeFile
	// NodeTypeSubsystem represents a logical grouping of files (module/package)
	NodeTypeSubsystem
	// NodeTypeDependency represents an external dependency
	NodeTypeDependency
	// NodeTypeIssue represents a GitHub issue
	NodeTypeIssue
)

// String returns the string representation of a NodeType.
func (nt NodeType) String() string {
	switch nt {
	case NodeTypeContributor:
		return "Contributor"
	case NodeTypeFile:
		return "File"
	case NodeTypeSubsystem:
		return "Subsystem"
	case NodeTypeDependency:
		return "Dependency"
	case NodeTypeIssue:
		return "Issue"
	default:
		return "Unknown"
	}
}

// EdgeType enumerates different relationship types in the temporal repository graph.
type EdgeType int

const (
	// EdgeTypeCollaboration represents contributors working together
	EdgeTypeCollaboration EdgeType = iota
	// EdgeTypeModification represents a contributor modifying a file
	EdgeTypeModification
	// EdgeTypeDependency represents a dependency relationship between files/packages
	EdgeTypeDependency
	// EdgeTypeIssueRelation represents a contributor/file related to an issue
	EdgeTypeIssueRelation
	// EdgeTypeContainment represents a file contained in a subsystem
	EdgeTypeContainment
)

// String returns the string representation of an EdgeType.
func (et EdgeType) String() string {
	switch et {
	case EdgeTypeCollaboration:
		return "Collaboration"
	case EdgeTypeModification:
		return "Modification"
	case EdgeTypeDependency:
		return "Dependency"
	case EdgeTypeIssueRelation:
		return "IssueRelation"
	case EdgeTypeContainment:
		return "Containment"
	default:
		return "Unknown"
	}
}

// Node represents an entity in the temporal repository graph.
// Nodes can be contributors, files, subsystems, dependencies, or issues.
type Node struct {
	// ID is the unique identifier for this node
	ID string

	// Type indicates what kind of entity this node represents
	Type NodeType

	// Timestamp indicates when this node was introduced
	Timestamp time.Time

	// Metadata stores additional information about the node
	// Examples: name, email, path, complexity_score, etc.
	Metadata map[string]interface{}

	// Occurrences tracks how many times this node appeared in analysis
	Occurrences int
}

// NewNode creates a new node with the given parameters.
func NewNode(id string, nodeType NodeType) *Node {
	return &Node{
		ID:          id,
		Type:        nodeType,
		Timestamp:   time.Now(),
		Metadata:    make(map[string]interface{}),
		Occurrences: 1,
	}
}

// SetMetadata sets a metadata key-value pair.
func (n *Node) SetMetadata(key string, value interface{}) {
	if n.Metadata == nil {
		n.Metadata = make(map[string]interface{})
	}
	n.Metadata[key] = value
}

// GetMetadata retrieves a metadata value, returning nil if not found.
func (n *Node) GetMetadata(key string) interface{} {
	return n.Metadata[key]
}

// Edge represents a relationship between two nodes in the temporal repository graph.
// Edges have weights representing the strength of the relationship.
type Edge struct {
	// Source is the originating node
	Source *Node

	// Target is the destination node
	Target *Node

	// Type indicates what kind of relationship this edge represents
	Type EdgeType

	// Weight represents the strength of the relationship (typically 0.0 to 1.0)
	Weight float64

	// Timestamp indicates when this edge was first observed
	Timestamp time.Time

	// Occurrences tracks how many times this relationship was observed
	Occurrences int
}

// NewEdge creates a new edge with the given parameters.
func NewEdge(source, target *Node, edgeType EdgeType, weight float64) *Edge {
	return &Edge{
		Source:      source,
		Target:      target,
		Type:        edgeType,
		Weight:      weight,
		Timestamp:   time.Now(),
		Occurrences: 1,
	}
}

// IncreaseWeight increases the edge weight using a running average.
// Updates the weight to reflect the average across all occurrences,
// clamped to [0, 1].
func (e *Edge) IncreaseWeight(delta float64) {
	// Compute new average: (old_avg * occurrences + delta) / (occurrences + 1)
	total := e.Weight*float64(e.Occurrences) + delta
	e.Occurrences++
	e.Weight = total / float64(e.Occurrences)

	// Clamp to valid range [0, 1]
	if e.Weight > 1.0 {
		e.Weight = 1.0
	}
	if e.Weight < 0.0 {
		e.Weight = 0.0
	}
}

// AverageWeight returns the current edge weight, which is maintained
// as a running average across all occurrences.
func (e *Edge) AverageWeight() float64 {
	return e.Weight
}

// GraphMetrics captures computed metrics about a graph.
type GraphMetrics struct {
	// NodeCount is the total number of nodes
	NodeCount int

	// EdgeCount is the total number of edges
	EdgeCount int

	// AverageNodeDegree is the average degree across all nodes
	AverageNodeDegree float64

	// AverageClusteringCoefficient is the average clustering coefficient
	AverageClusteringCoefficient float64

	// Density is the ratio of actual edges to possible edges
	Density float64

	// ConnectedComponentCount is the number of disconnected subgraphs
	ConnectedComponentCount int
}
