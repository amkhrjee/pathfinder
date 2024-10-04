package algorithm

import (
	"pfinder/grid"
)

func neighbors(g *grid.Grid, b *grid.Box) []*grid.Box {
	n := make([]*grid.Box, 0)
	// integer values
	bRow := int(b.Row)
	bCol := int(b.Col)
	// north neighbor
	if b.Row < grid.ROWS-1 {
		n = append(n, &g[bRow+1][bCol])
	}
	// south neighbor
	if b.Row > 0 {
		n = append(n, &g[bRow-1][bCol])
	}
	// right neighbor
	if b.Col < grid.COLS-1 {
		n = append(n, &g[bRow][bCol+1])
	}
	// left neighbor
	if b.Col > 0 {
		n = append(n, &g[bRow][bCol-1])
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
