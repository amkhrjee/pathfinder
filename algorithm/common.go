package algorithm

import (
	"math"
	"pfinder/grid"
)

func get_distance(x *grid.Box, y *grid.Box) float64 {
	// Good idea to scale down the distance
	// to be in the range of relative cost between the tiles
	// Thanks to this reddit comment: https://www.reddit.com/r/computerscience/comments/1fw9q09/comment/lqeqgdd/?utm_source=share&utm_medium=web3x&utm_name=web3xcss&utm_term=1&utm_content=share_button
	// Scaling with formula from: https://stackoverflow.com/a/5295202/12404524
	actual_distance := math.Sqrt(math.Pow(x.Col-y.Col, 2) + math.Pow(x.Row-y.Row, 2))
	scaled_distance := 10 * (actual_distance - 0) / 26.8
	return scaled_distance
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
