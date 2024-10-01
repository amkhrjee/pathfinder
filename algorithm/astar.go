package algorithm

import (
	"math"
	"pfinder/grid"
	"slices"

	pq "gopkg.in/dnaeon/go-priorityqueue.v1"
)

func get_distance(x *grid.Box, y *grid.Box) float64 {
	return math.Sqrt(math.Pow(float64(x.Col-y.Col), 2) + math.Pow(float64(x.Row-y.Row), 2))
}

// A* Search
func AStar(g *grid.Grid, start *grid.Box, target *grid.Box) ([]*grid.Box, []*grid.Box) {
	root := start
	q := pq.New[*grid.Box, float64](pq.MinHeap)
	track := make([]*grid.Box, 0)
	final_path := make([]*grid.Box, 0)
	// inserting the initial nodes
	for _, n := range neighbors(g, root) {
		relative_cost := math.Abs(float64(n.Cost - root.Cost))
		euclidean_dist := get_distance(target, n)
		q.Put(n, relative_cost+euclidean_dist)
	}
	for !q.IsEmpty() {
		curr := q.Get()
		track = append(track, curr.Value)
		if curr.Value == target {
			break
		}

		for _, n := range neighbors(g, curr.Value) {
			if !slices.Contains(track, n) && n != start {
				euclidean_dist := get_distance(target, n)
				relative_cost := math.Abs(float64(n.Cost - curr.Value.Cost))
				total_cost := curr.Priority + relative_cost + euclidean_dist
				q.Put(n, total_cost)
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
