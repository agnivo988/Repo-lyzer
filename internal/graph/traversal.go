package graph

import (
	"fmt"
	"sort"
)

// TraversalVisitor is a function type for visiting nodes during graph traversal.
// If it returns an error, traversal stops immediately and the error is returned.
type TraversalVisitor func(*Node) error

// DFS performs a depth-first search starting from the given node.
// It calls the visitor function for each visited node.
// Returns an error if the visitor returns an error or if the start node is not found.
func (g *Graph) DFS(startID string, visitor TraversalVisitor) error {
	if visitor == nil {
		return fmt.Errorf("visitor cannot be nil")
	}

	node, ok := g.getInternalNode(startID)
	if !ok {
		return fmt.Errorf("start node %s not found", startID)
	}

	visited := make(map[string]bool)
	return g.dfsHelper(node, visitor, visited)
}

// dfsHelper is the recursive helper for DFS.
func (g *Graph) dfsHelper(node *Node, visitor TraversalVisitor, visited map[string]bool) error {
	if visited[node.ID] {
		return nil
	}

	visited[node.ID] = true

	if err := visitor(node); err != nil {
		return err
	}

	g.mutex.RLock()
	edges := g.edgeIndex[node.ID]
	g.mutex.RUnlock()

	for _, edge := range edges {
		if err := g.dfsHelper(edge.Target, visitor, visited); err != nil {
			return err
		}
	}

	return nil
}

// BFS performs a breadth-first search starting from the given node.
// It calls the visitor function for each visited node in level order.
// Returns an error if the visitor returns an error or if the start node is not found.
func (g *Graph) BFS(startID string, visitor TraversalVisitor) error {
	if visitor == nil {
		return fmt.Errorf("visitor cannot be nil")
	}

	node, ok := g.getInternalNode(startID)
	if !ok {
		return fmt.Errorf("start node %s not found", startID)
	}

	visited := make(map[string]bool)
	queue := []*Node{node}
	visited[node.ID] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if err := visitor(current); err != nil {
			return err
		}

		g.mutex.RLock()
		edges := g.edgeIndex[current.ID]
		g.mutex.RUnlock()

		for _, edge := range edges {
			if !visited[edge.Target.ID] {
				visited[edge.Target.ID] = true
				queue = append(queue, edge.Target)
			}
		}
	}

	return nil
}

// ShortestPath finds the shortest path between two nodes using BFS.
// Returns the path as a slice of nodes, or nil if no path exists.
func (g *Graph) ShortestPath(fromID, toID string) ([]*Node, error) {
	fromNode, ok := g.getInternalNode(fromID)
	if !ok {
		return nil, fmt.Errorf("node %s not found", fromID)
	}

	toNode, ok := g.getInternalNode(toID)
	if !ok {
		return nil, fmt.Errorf("node %s not found", toID)
	}

	if fromID == toID {
		return []*Node{fromNode}, nil
	}

	path := g.shortestPath(fromNode, toNode)
	if path == nil {
		return nil, fmt.Errorf("no path found from %s to %s", fromID, toID)
	}

	return path, nil
}

// AllPaths finds all paths between two nodes.
// This can be expensive for large graphs with many cycles.
// Paths are limited to avoid exponential explosion.
func (g *Graph) AllPaths(fromID, toID string, maxDepth int) ([][]*Node, error) {
	fromNode, ok := g.getInternalNode(fromID)
	if !ok {
		return nil, fmt.Errorf("node %s not found", fromID)
	}

	if fromID == toID {
		return [][]*Node{{fromNode}}, nil
	}

	var paths [][]*Node
	visited := make(map[string]bool)
	g.findAllPathsHelper(fromNode, toID, []*Node{}, &paths, visited, maxDepth)

	return paths, nil
}

// findAllPathsHelper is the recursive helper for finding all paths.
func (g *Graph) findAllPathsHelper(current *Node, targetID string, path []*Node, paths *[][]*Node, visited map[string]bool, maxDepth int) {
	if maxDepth <= 0 {
		return
	}

	path = append(path, current)
	visited[current.ID] = true

	if current.ID == targetID {
		// Found a path; add a copy to results
		result := make([]*Node, len(path))
		copy(result, path)
		*paths = append(*paths, result)
	} else {
		g.mutex.RLock()
		edges := g.edgeIndex[current.ID]
		g.mutex.RUnlock()

		for _, edge := range edges {
			if !visited[edge.Target.ID] {
				g.findAllPathsHelper(edge.Target, targetID, path, paths, visited, maxDepth-1)
			}
		}
	}

	// Backtrack
	visited[current.ID] = false
}

// DistanceTo computes the shortest distance (number of edges) from one node to another.
// Returns -1 if no path exists.
func (g *Graph) DistanceTo(fromID, toID string) int {
	fromNode, ok := g.getInternalNode(fromID)
	if !ok {
		return -1
	}

	if fromID == toID {
		return 0
	}

	g.mutex.RLock()
	defer g.mutex.RUnlock()

	visited := make(map[string]bool)
	queue := []struct {
		node     *Node
		distance int
	}{
		{fromNode, 0},
	}
	visited[fromNode.ID] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for _, edge := range g.edgeIndex[current.node.ID] {
			if edge.Target.ID == toID {
				return current.distance + 1
			}

			if !visited[edge.Target.ID] {
				visited[edge.Target.ID] = true
				queue = append(queue, struct {
					node     *Node
					distance int
				}{edge.Target, current.distance + 1})
			}
		}
	}

	return -1 // No path found
}

// IncomingEdges returns all edges pointing to a given node.
func (g *Graph) IncomingEdges(nodeID string) []*Edge {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	var incomingEdges []*Edge
	for _, edges := range g.edgeIndex {
		for _, edge := range edges {
			if edge.Target.ID == nodeID {
				incomingEdges = append(incomingEdges, edge)
			}
		}
	}

	return incomingEdges
}

// Neighborhood returns all nodes directly connected to the given node
// (both incoming and outgoing edges).
func (g *Graph) Neighborhood(nodeID string) []*Node {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	neighbors := make(map[string]*Node)

	// Add outgoing neighbors
	if outEdges, ok := g.edgeIndex[nodeID]; ok {
		for _, edge := range outEdges {
			neighbors[edge.Target.ID] = edge.Target
		}
	}

	// Add incoming neighbors
	for _, edges := range g.edgeIndex {
		for _, edge := range edges {
			if edge.Target.ID == nodeID {
				neighbors[edge.Source.ID] = edge.Source
			}
		}
	}

	result := make([]*Node, 0, len(neighbors))
	for _, node := range neighbors {
		result = append(result, node)
	}

	return result
}

// FilterNodes returns nodes that satisfy the given predicate.
func (g *Graph) FilterNodes(predicate func(*Node) bool) []*Node {
	return g.Query(predicate)
}

// NodesByDistance returns a list of nodes sorted by their distance from a source node.
// Nodes are returned in order of increasing distance.
func (g *Graph) NodesByDistance(sourceID string) ([]*Node, []int, error) {
	sourceNode, ok := g.getInternalNode(sourceID)
	if !ok {
		return nil, nil, fmt.Errorf("source node %s not found", sourceID)
	}

	// BFS to compute distances
	visited := make(map[string]int) // node ID -> distance
	queue := []struct {
		node     *Node
		distance int
	}{
		{sourceNode, 0},
	}
	visited[sourceNode.ID] = 0

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		g.mutex.RLock()
		edges := g.edgeIndex[current.node.ID]
		g.mutex.RUnlock()

		for _, edge := range edges {
			if _, seen := visited[edge.Target.ID]; !seen {
				visited[edge.Target.ID] = current.distance + 1
				queue = append(queue, struct {
					node     *Node
					distance int
				}{edge.Target, current.distance + 1})
			}
		}
	}

	// Collect and sort by distance
	type pair struct {
		id string
		d  int
	}
	pairs := make([]pair, 0, len(visited))
	for id, d := range visited {
		pairs = append(pairs, pair{id: id, d: d})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].d == pairs[j].d {
			return pairs[i].id < pairs[j].id
		}
		return pairs[i].d < pairs[j].d
	})

	nodes := make([]*Node, 0, len(pairs))
	distances := make([]int, 0, len(pairs))
	for _, p := range pairs {
		if n, ok := g.getInternalNode(p.id); ok {
			nodes = append(nodes, cloneNode(n))
			distances = append(distances, p.d)
		}
	}

	return nodes, distances, nil
}
