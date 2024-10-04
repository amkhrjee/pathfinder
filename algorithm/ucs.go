package algorithm

import (
	"math"
	"pfinder/grid"

	pq "gopkg.in/dnaeon/go-priorityqueue.v1"
)

func Ucs(g *grid.Grid, start *grid.Box, target *grid.Box) ([]*grid.Box, []*grid.Box) {
	root := start
	q := pq.New[*grid.Box, float64](pq.MinHeap)
	track := make([]*grid.Box, 0)
	final_path := make([]*grid.Box, 0)
	// inserting the initial nodes
	for _, n := range neighbors(g, root) {
		relative_cost := math.Abs(float64(n.Cost - root.Cost))
		q.Put(n, relative_cost)
		n.Visited = true
	}
	for !q.IsEmpty() {
		curr := q.Get()
		track = append(track, curr.Value)
		if curr.Value == target {
			break
		}

		for _, n := range neighbors(g, curr.Value) {
			if !n.Visited {
				total_cost := curr.Priority + math.Abs(float64(n.Cost-curr.Value.Cost))
				q.Put(n, total_cost)
				n.Parent = curr.Value
				n.Visited = true
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
