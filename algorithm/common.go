package algorithm

import (
	"math"
	"pfinder/grid"
)

func get_distance(x *grid.Box, y *grid.Box) float64 {
	return math.Sqrt(math.Pow(x.Col-y.Col, 2) + math.Pow(x.Row-y.Row, 2))
}

func neighbors(g *grid.Grid, b *grid.Box) []*grid.Box {
	n := make([]*grid.Box, 0)
	// integer values
	bRow := int(b.Row)
	bCol := int(b.Col)
	// north neighbor
	if b.Row < grid.ROWS-1 {
		box := &g[bRow+1][bCol]
		if !box.IsObstacle {
			n = append(n, box)
		}
	}
	// south neighbor
	if b.Row > 0 {
		box := &g[bRow-1][bCol]
		if !box.IsObstacle {
			n = append(n, box)
		}
	}
	// right neighbor
	if b.Col < grid.COLS-1 {
		box := &g[bRow][bCol+1]
		if !box.IsObstacle {
			n = append(n, box)
		}
	}
	// left neighbor
	if b.Col > 0 {
		box := &g[bRow][bCol-1]
		if !box.IsObstacle {
			n = append(n, box)
		}
	}
	// // north west
	// if b.Col > 0 && b.Row > 0 {
	// 	n = append(n, &g[b.Row-1][b.Col-1])
	// }
	// // north east
	// if b.Col < grid.COLS-1 && b.Row > 0 {
	// 	n = append(n, &g[b.Row-1][b.Col+1])
	// }
	// // south west
	// if b.Col > 0 && b.Row < grid.ROWS-1 {
	// 	n = append(n, &g[b.Row+1][b.Col-1])
	// }
	// // south east
	// if b.Col < grid.COLS && b.Row < grid.ROWS-1 {
	// 	n = append(n, &g[b.Row+1][b.Col+1])
	// }
	return n
}

const (
	_ int = iota
	ASTAR
	UCS
	BFS
	DFS
	LBEAM
)

const ALGORITHMS_COUNT = 6
