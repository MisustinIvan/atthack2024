package node

import (
	"container/heap"
	"errors"
	"iter"
	"math"
	"optitraffic/pqueue"
)

type Pathfinder struct {
	graph *Graph
}

func NewPathfinder(g *Graph) Pathfinder {
	return Pathfinder{
		graph: g,
	}
}

func (pf *Pathfinder) Path(from, to *Node) ([]*Node, error) {
	// Edge case: if the from and to nodes are the same
	if from == to {
		return []*Node{from}, nil
	}

	// Initialize distances and the priority queue
	distances := make(map[*Node]float64)
	cameFrom := make(map[*Node]*Node)
	for _, node := range pf.graph.Nodes {
		distances[node] = math.MaxFloat64
	}
	distances[from] = 0

	// Create a priority queue and add the starting node
	pq := make(pqueue.PriorityQueue[*Node], 0)
	heap.Init(&pq)
	heap.Push(&pq, pqueue.NewItem(from, 0, 0))

	// Dijkstra's algorithm
	for pq.Len() > 0 {
		currentItem := heap.Pop(&pq).(*pqueue.Item[*Node])
		currentNode := currentItem.Value()

		// If we reached the target node, stop processing
		if currentNode == to {
			break
		}

		// Relaxation: Update distances to neighbors
		for neighbor, params := range currentNode.Conns {
			newDist := distances[currentNode] + params.ConnectionWeight()
			if newDist < distances[neighbor] {
				distances[neighbor] = newDist
				cameFrom[neighbor] = currentNode
				heap.Push(&pq, pqueue.NewItem(neighbor, int(newDist), neighbor.Id))
			}
		}
	}

	// Reconstruct the path
	path := []*Node{}
	for current := to; current != nil; current = cameFrom[current] {
		path = append([]*Node{current}, path...)
	}

	// If we couldn't reach the destination node, return an error
	if len(path) == 0 || path[0] != from {
		return nil, errors.New("no path found")
	}

	return path, nil
}

func PathLength(nodes []*Node) float64 {
	res := 0.0

	for n1, n2 := range IntoPairIter(nodes) {
		res += n1.Conns[n2].Dist
	}

	return res
}

func IntoPairIter[T any](xs []T) iter.Seq2[T, T] {
	return func(yield func(T, T) bool) {
		for i := 0; i < len(xs)-1; i++ {
			a := xs[i]
			b := xs[i+1]
			if !yield(a, b) {
				return
			}
		}
	}
}
