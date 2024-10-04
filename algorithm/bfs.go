package algorithm

import (
	"pfinder/grid"
	"slices"

	pq "gopkg.in/dnaeon/go-priorityqueue.v1"
)

// Breadth First Search
func Bfs(g *grid.Grid, start *grid.Box, target *grid.Box) ([]*grid.Box, []*grid.Box) {
	root := start
	q := pq.New[*grid.Box, float64](pq.MinHeap)
	track := make([]*grid.Box, 0)
	final_path := make([]*grid.Box, 0)
	// inserting the initial nodes
	for _, n := range neighbors(g, root) {
		depth := 1.
		q.Put(n, depth)
	}
	for !q.IsEmpty() {

		curr := q.Get()
		track = append(track, curr.Value)
		if curr.Value == target {
			break
		}

		for _, n := range neighbors(g, curr.Value) {
			if !slices.Contains(track, n) && n != start {
				depth := curr.Priority + 1
				q.Put(n, depth)
				n.Parent = curr.Value
			}

		}
	}
	// backtracking the path
	curr := target
	for curr != nil {
		final_path = append(final_path, curr)
		curr = curr.Parent
	}

	return track, final_path
}
