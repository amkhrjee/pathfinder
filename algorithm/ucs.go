package algorithm

import (
	"fmt"
	"math"
	"pfinder/grid"
	"slices"

	pq "gopkg.in/dnaeon/go-priorityqueue.v1"
)

func neighbors(g *grid.Grid, b *grid.Box) []*grid.Box {
	n := make([]*grid.Box, 0)
	// north neighbor
	if b.Row < grid.ROWS-1 {
		n = append(n, &g[b.Row+1][b.Col])
	}
	// south neighbor
	if b.Row > 0 {
		n = append(n, &g[b.Row-1][b.Col])
	}
	// right neighbor
	if b.Col < grid.COLS-1 {
		n = append(n, &g[b.Row][b.Col+1])
	}
	// left neighbor
	if b.Col > 0 {
		n = append(n, &g[b.Row][b.Col-1])
	}
	return n
}

// Uniform Cost Search
func Ucs(g *grid.Grid, start *grid.Box, target *grid.Box) ([]*grid.Box, []*grid.Box) {
	root := start
	q := pq.New[*grid.Box, float64](pq.MinHeap)
	track := make([]*grid.Box, 0)
	final_path := make([]*grid.Box, 0)
	// inserting the initial nodes
	for _, n := range neighbors(g, root) {
		relative_cost := math.Abs(float64(n.Cost - root.Cost))
		q.Put(n, relative_cost)
	}
	for !q.IsEmpty() {
		curr := q.Get()
		track = append(track, curr.Value)
		if curr.Value == target {
			break
		}

		for _, n := range neighbors(g, curr.Value) {
			if !slices.Contains(track, n) && n != start {
				total_cost := curr.Priority + math.Abs(float64(n.Cost-curr.Value.Cost))
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

	fmt.Println("Length of final path: ", len(final_path))

	return track, final_path
}
