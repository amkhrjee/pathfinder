package algorithm

import (
	"math"
	"pfinder/grid"
	"slices"

	pq "gopkg.in/dnaeon/go-priorityqueue.v1"
)

// Local Beam Search
func LBeam(g *grid.Grid, start *grid.Box, target *grid.Box, k int) ([]*grid.Box, []*grid.Box) {
	// just to be on the safer side
	if k > 4 || k < 1 {
		k = 4
	}
	root := start
	q := pq.New[*grid.Box, float64](pq.MinHeap)
	track := make([]*grid.Box, 0)
	final_path := make([]*grid.Box, 0)
	temp := pq.New[*grid.Box, float64](pq.MinHeap)
	// inserting the initial nodes
	for _, n := range neighbors(g, root) {
		relative_cost := math.Abs(float64(n.Cost - root.Cost))
		euclidean_dist := get_distance(target, n)
		temp.Put(n, relative_cost+euclidean_dist)
	}

	for range k {
		if !temp.IsEmpty() {
			n := temp.Get()
			q.Put(n.Value, n.Priority)
		}
	}

	for !q.IsEmpty() {
		curr := q.Get()
		track = append(track, curr.Value)
		if curr.Value == target {
			break
		}

		temp := pq.New[*grid.Box, float64](pq.MinHeap)
		for _, n := range neighbors(g, curr.Value) {
			if !slices.Contains(track, n) && n != start {
				euclidean_dist := get_distance(target, n)
				relative_cost := math.Abs(float64(n.Cost - curr.Value.Cost))
				total_cost := curr.Priority + relative_cost + euclidean_dist
				temp.Put(n, total_cost)
				n.Parent = curr.Value
			}
		}
		for range k {
			if !temp.IsEmpty() {
				n := temp.Get()
				q.Put(n.Value, n.Priority)
			}
		}
	}
	// backtracking the path
	curr := target
	counter := 0
	for curr != nil {
		final_path = append(final_path, curr)
		curr = curr.Parent
		counter++
	}

	return track, final_path
}
