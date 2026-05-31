package graph

// DegreeCentrality computes the degree centrality of a node.
// Degree centrality is the number of edges connected to a node.
// Normalized to [0, 1] by dividing by 2*(n-1) for directed graphs,
// where n is the total number of nodes.
func (g *Graph) DegreeCentrality(node *Node) float64 {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	inDegree := 0
	outDegree := len(g.edgeIndex[node.ID])

	// Count incoming edges
	for _, edges := range g.edgeIndex {
		for _, edge := range edges {
			if edge.Target.ID == node.ID {
				inDegree++
			}
		}
	}

	totalDegree := inDegree + outDegree
	if len(g.nodes) <= 1 {
		return 0
	}

	return float64(totalDegree) / float64(2*(len(g.nodes)-1))
}

// ClusteringCoefficient computes the clustering coefficient of a node.
// It measures the degree to which nodes cluster together.
// The clustering coefficient of a node is the ratio of edges between its neighbors
// to the maximum possible edges between them.
func (g *Graph) ClusteringCoefficient(node *Node) float64 {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	// Get all neighbors (both incoming and outgoing)
	neighbors := make(map[string]bool)

	// Add outgoing neighbors
	for _, edge := range g.edgeIndex[node.ID] {
		neighbors[edge.Target.ID] = true
	}

	// Add incoming neighbors
	for _, edges := range g.edgeIndex {
		for _, edge := range edges {
			if edge.Target.ID == node.ID {
				neighbors[edge.Source.ID] = true
			}
		}
	}

	numNeighbors := len(neighbors)
	if numNeighbors < 2 {
		return 0 // Can't form a triangle with fewer than 2 neighbors
	}

	// Count edges between neighbors (directed)
	edgesBetweenNeighbors := 0
	for neighborID := range neighbors {
		for _, edge := range g.edgeIndex[neighborID] {
			if other, exists := neighbors[edge.Target.ID]; exists && other {
				edgesBetweenNeighbors++
			}
		}
	}

	// Maximum possible edges between neighbors in a directed graph
	maxEdges := numNeighbors * (numNeighbors - 1)

	if maxEdges == 0 {
		return 0
	}

	return float64(edgesBetweenNeighbors) / float64(maxEdges)
}

// AverageDegreeCentrality computes the average degree centrality across all nodes.
func (g *Graph) AverageDegreeCentrality() float64 {
	nodes := g.GetAllNodes()
	if len(nodes) == 0 {
		return 0
	}

	sum := 0.0
	for _, node := range nodes {
		sum += g.DegreeCentrality(node)
	}

	return sum / float64(len(nodes))
}

// AverageClusteringCoefficient computes the average clustering coefficient across all nodes.
func (g *Graph) AverageClusteringCoefficient() float64 {
	nodes := g.GetAllNodes()
	if len(nodes) == 0 {
		return 0
	}

	sum := 0.0
	for _, node := range nodes {
		sum += g.ClusteringCoefficient(node)
	}

	return sum / float64(len(nodes))
}

// Density computes the density of the graph.
// Density is the ratio of actual edges to the maximum possible edges.
// For a directed graph: density = actual_edges / (n * (n - 1))
func (g *Graph) Density() float64 {
	g.mutex.RLock()
	defer g.mutex.RUnlock()

	numNodes := len(g.nodes)
	if numNodes <= 1 {
		return 0
	}

	maxEdges := numNodes * (numNodes - 1)
	actualEdges := len(g.edges)

	return float64(actualEdges) / float64(maxEdges)
}

// ConnectedComponents finds all connected components in the graph.
// Uses depth-first search to identify disconnected subgraphs.
func (g *Graph) ConnectedComponents() [][]Node {
	// Don't hold lock while calling GetAllNodes, which acquires its own lock
	nodes := g.GetAllNodes()

	visited := make(map[string]bool)
	var components [][]Node

	for _, node := range nodes {
		if !visited[node.ID] {
			component := g.dfsComponent(node, visited)
			components = append(components, component)
		}
	}

	return components
}

// dfsComponent is a helper function that performs DFS to find all nodes in a component.
func (g *Graph) dfsComponent(node *Node, visited map[string]bool) []Node {
	visited[node.ID] = true
	component := []Node{*node}

	// Visit all neighbors
	g.mutex.RLock()
	edges := g.edgeIndex[node.ID]
	g.mutex.RUnlock()

	for _, edge := range edges {
		if !visited[edge.Target.ID] {
			component = append(component, g.dfsComponent(edge.Target, visited)...)
		}
	}

	// Visit nodes that point to this node
	allEdges := g.GetAllEdges()
	for _, edge := range allEdges {
		if edge.Target.ID == node.ID && !visited[edge.Source.ID] {
			component = append(component, g.dfsComponent(edge.Source, visited)...)
		}
	}

	return component
}

// ComputeMetrics computes and returns the graph's structural metrics.
func (g *Graph) ComputeMetrics() GraphMetrics {
	return GraphMetrics{
		NodeCount:                    g.NodeCount(),
		EdgeCount:                    g.EdgeCount(),
		AverageNodeDegree:            g.AverageDegreeCentrality() * float64(g.NodeCount()-1),
		AverageClusteringCoefficient: g.AverageClusteringCoefficient(),
		Density:                      g.Density(),
		ConnectedComponentCount:      len(g.ConnectedComponents()),
	}
}

// BetweennessCentrality computes the betweenness centrality of a node.
// Betweenness centrality is the number of shortest paths between other nodes
// that pass through the given node. This is a simplified O(n^2) implementation
// suitable for small to medium graphs.
func (g *Graph) BetweennessCentrality(node *Node) float64 {
	// Don't hold lock while calling GetAllNodes, which acquires its own lock
	allNodes := g.GetAllNodes()

	if len(allNodes) <= 2 {
		return 0
	}

	count := 0.0

	// For each pair of other nodes
	for i, source := range allNodes {
		if source.ID == node.ID {
			continue
		}

		for j, target := range allNodes {
			if i >= j || target.ID == node.ID {
				continue
			}

			// Find shortest path from source to target
			path := g.shortestPath(source, target)
			if path != nil {
				// Check if our node is in the path (not counting endpoints)
				for k := 1; k < len(path)-1; k++ {
					if path[k].ID == node.ID {
						count++
						break
					}
				}
			}
		}
	}

	// Normalize by maximum possible betweenness
	numOtherNodes := len(allNodes) - 1
	maxBetweenness := float64(numOtherNodes * (numOtherNodes - 1) / 2)

	if maxBetweenness == 0 {
		return 0
	}

	return count / maxBetweenness
}

// shortestPath finds the shortest path between two nodes using BFS.
func (g *Graph) shortestPath(from, to *Node) []*Node {
	if from.ID == to.ID {
		return []*Node{from}
	}

	g.mutex.RLock()
	defer g.mutex.RUnlock()

	visited := make(map[string]bool)
	queue := [][]*Node{{from}}
	visited[from.ID] = true

	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:]

		current := path[len(path)-1]

		// Check all neighbors
		for _, edge := range g.edgeIndex[current.ID] {
			if edge.Target.ID == to.ID {
				// Found target
				return append(path, edge.Target)
			}

			if !visited[edge.Target.ID] {
				visited[edge.Target.ID] = true
				newPath := make([]*Node, len(path)+1)
				copy(newPath, path)
				newPath[len(newPath)-1] = edge.Target
				queue = append(queue, newPath)
			}
		}
	}

	return nil // No path found
}

// ClosestNodesByDistance finds the N closest nodes to a given node based on
// some distance metric. Uses degree-based distance as default.
func (g *Graph) ClosestNodesByDistance(node *Node, count int) []*Node {
	if node == nil {
		return nil
	}

	nodes, _, err := g.NodesByDistance(node.ID)
	if err != nil {
		return nil
	}

	if count <= 0 || count > len(nodes) {
		count = len(nodes)
	}

	return nodes[:count]
}
