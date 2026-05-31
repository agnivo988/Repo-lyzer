package graph

import (
	"fmt"
	"sync"
)

// cloneNode creates a deep copy of a Node to avoid exposing internal pointers.
func cloneNode(n *Node) *Node {
	if n == nil {
		return nil
	}
	m := make(map[string]interface{})
	for k, v := range n.Metadata {
		m[k] = v
	}
	return &Node{
		ID:          n.ID,
		Type:        n.Type,
		Timestamp:   n.Timestamp,
		Metadata:    m,
		Occurrences: n.Occurrences,
	}
}

// cloneEdge creates a deep copy of an Edge, cloning source/target nodes.
func cloneEdge(e *Edge) *Edge {
	if e == nil {
		return nil
	}
	return &Edge{
		Source:      cloneNode(e.Source),
		Target:      cloneNode(e.Target),
		Type:        e.Type,
		Weight:      e.Weight,
		Timestamp:   e.Timestamp,
		Occurrences: e.Occurrences,
	}
}

// getInternalNode returns the internal node pointer without cloning. The
// caller must hold the appropriate locks if necessary.
func (g *Graph) getInternalNode(id string) (*Node, bool) {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	n, ok := g.nodes[id]
	return n, ok
}

// Graph represents the temporal repository ecosystem as a directed graph.
// It maintains collections of nodes and edges with efficient lookup indices.
type Graph struct {
	// nodes maps node ID to node for O(1) lookup
	nodes map[string]*Node

	// edges stores all edges in the graph
	edges []*Edge

	// nodeIndex maps node type to list of nodes of that type
	nodeIndex map[NodeType][]*Node

	// edgeIndex maps source node ID to list of outgoing edges
	edgeIndex map[string][]*Edge

	// mutex protects concurrent access to the graph
	mutex sync.RWMutex
}

// NewGraph creates a new empty graph.
func NewGraph() *Graph {
	return &Graph{
		nodes:     make(map[string]*Node),
		edges:     make([]*Edge, 0),
		nodeIndex: make(map[NodeType][]*Node),
		edgeIndex: make(map[string][]*Edge),
	}
}

// AddNode adds a node to the graph, returning an error if a node with the same ID exists.
func (g *Graph) AddNode(node *Node) error {
	if node == nil {
		return fmt.Errorf("cannot add nil node")
	}

	g.mutex.Lock()
	defer g.mutex.Unlock()

	if _, exists := g.nodes[node.ID]; exists {
		return fmt.Errorf("node with ID %s already exists", node.ID)
	}

	g.nodes[node.ID] = node
	g.nodeIndex[node.Type] = append(g.nodeIndex[node.Type], node)

	return nil
}

// AddEdge adds an edge to the graph. If the edge already exists between the same source and target
// with the same type, it increments the weight instead of duplicating.
func (g *Graph) AddEdge(edge *Edge) error {
	if edge == nil {
		return fmt.Errorf("cannot add nil edge")
	}

	if edge.Source == nil || edge.Target == nil {
		return fmt.Errorf("edge must have non-nil source and target nodes")
	}

	g.mutex.Lock()
	defer g.mutex.Unlock()

	// Validate that both nodes exist in the graph
	if _, ok := g.nodes[edge.Source.ID]; !ok {
		return fmt.Errorf("source node %s does not exist in graph", edge.Source.ID)
	}
	if _, ok := g.nodes[edge.Target.ID]; !ok {
		return fmt.Errorf("target node %s does not exist in graph", edge.Target.ID)
	}

	// Check if edge already exists
	if existingEdges, ok := g.edgeIndex[edge.Source.ID]; ok {
		for _, e := range existingEdges {
			if e.Target.ID == edge.Target.ID && e.Type == edge.Type {
				// Edge exists; increase weight
				e.IncreaseWeight(edge.Weight)
				return nil
			}
		}
	}

	// Add new edge
	g.edges = append(g.edges, edge)
	g.edgeIndex[edge.Source.ID] = append(g.edgeIndex[edge.Source.ID], edge)

	return nil
}

// GetNode retrieves a node by ID, returning an error if not found.
func (g *Graph) GetNode(id string) (*Node, error) {
	g.mutex.RLock()
	node, exists := g.nodes[id]
	g.mutex.RUnlock()
	if !exists {
		return nil, fmt.Errorf("node with ID %s not found", id)
	}
	return cloneNode(node), nil
}

// GetEdges retrieves all outgoing edges from a source node.
func (g *Graph) GetEdges(sourceID string) []*Edge {
	g.mutex.RLock()
	edges := g.edgeIndex[sourceID]
	g.mutex.RUnlock()

	result := make([]*Edge, 0, len(edges))
	for _, e := range edges {
		result = append(result, cloneEdge(e))
	}
	return result
}

// Query returns all nodes matching the given predicate.
func (g *Graph) Query(predicate func(*Node) bool) []*Node {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	var results []*Node
	for _, node := range g.nodes {
		if predicate(node) {
			results = append(results, cloneNode(node))
		}
	}

	return results
}

// QueryByType returns all nodes of the specified type.
func (g *Graph) QueryByType(nodeType NodeType) []*Node {
	g.mutex.RLock()
	nodes := g.nodeIndex[nodeType]
	g.mutex.RUnlock()

	result := make([]*Node, 0, len(nodes))
	for _, n := range nodes {
		result = append(result, cloneNode(n))
	}
	return result
}

// NodeCount returns the total number of nodes in the graph.
func (g *Graph) NodeCount() int {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	return len(g.nodes)
}

// EdgeCount returns the total number of edges in the graph.
func (g *Graph) EdgeCount() int {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	return len(g.edges)
}

// GetAllNodes returns all nodes in the graph.
func (g *Graph) GetAllNodes() []*Node {
	g.mutex.RLock()
	nodes := make([]*Node, 0, len(g.nodes))
	for _, node := range g.nodes {
		nodes = append(nodes, cloneNode(node))
	}
	g.mutex.RUnlock()

	return nodes
}

// GetAllEdges returns all edges in the graph.
func (g *Graph) GetAllEdges() []*Edge {
	g.mutex.RLock()
	edges := make([]*Edge, 0, len(g.edges))
	for _, e := range g.edges {
		edges = append(edges, cloneEdge(e))
	}
	g.mutex.RUnlock()
	return edges
}

// Clear removes all nodes and edges from the graph.
func (g *Graph) Clear() {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	g.nodes = make(map[string]*Node)
	g.edges = make([]*Edge, 0)
	g.nodeIndex = make(map[NodeType][]*Node)
	g.edgeIndex = make(map[string][]*Edge)
}
