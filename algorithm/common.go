package algorithm

import "pfinder/grid"

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