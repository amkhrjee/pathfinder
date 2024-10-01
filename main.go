package main

import (
	"math"
	"math/rand"
	"slices"

	rl "github.com/gen2brain/raylib-go/raylib"
	pq "gopkg.in/dnaeon/go-priorityqueue.v1"
)

const WIDTH, HEIGHT = 800, 800
const BOX_DIM = 80

type Box struct {
	row       int32
	col       int32
	cost      int
	is_source bool
	is_target bool
}

const ROWS = HEIGHT / BOX_DIM
const COLS = WIDTH / BOX_DIM

var colors = []rl.Color{
	{245, 230, 255, 127},
	{220, 190, 255, 127},
	{195, 150, 255, 127},
	{170, 110, 255, 127},
	{145, 70, 255, 127},
	{120, 30, 255, 127},
	{95, 0, 230, 127},
	{70, 0, 180, 127},
	{45, 0, 130, 127},
	{20, 0, 80, 127},
}

type Grid [ROWS][COLS]Box

func makeGrid() *Grid {
	g := Grid{}
	for i, row := range g {
		for j := range row {
			g[i][j] = Box{
				row:       int32(i),
				col:       int32(j),
				is_source: false,
				is_target: false,
				cost:      rand.Intn(10-1) + 1,
			}
		}
	}

	return &g
}

func neighbors(g *Grid, b *Box) []*Box {
	n := make([]*Box, 0)
	// north neighbor
	if b.row < ROWS-1 {
		n = append(n, &g[b.row+1][b.col])
	}
	// south neighbor
	if b.row > 0 {
		n = append(n, &g[b.row-1][b.col])
	}
	// right neighbor
	if b.col < COLS-1 {
		n = append(n, &g[b.row][b.col+1])
	}
	// left neighbor
	if b.col > 0 {
		n = append(n, &g[b.row][b.col-1])
	}
	return n
}

// Uniform Cost Search
func ucs(grid *Grid, start *Box, target *Box) []*Box {
	root := start
	q := pq.New[*Box, float64](pq.MinHeap)
	track := make([]*Box, 0)
	// inserting the initial nodes
	for _, n := range neighbors(grid, root) {
		q.Put(n, math.Abs(float64(n.cost-root.cost)))
	}
	for !q.IsEmpty() {
		curr := q.Get()
		track = append(track, curr.Value)
		if curr.Value == target {
			break
		}

		for _, n := range neighbors(grid, curr.Value) {
			if !slices.Contains(track, n) {
				q.Put(n, curr.Priority+math.Abs(float64(n.cost-curr.Value.cost)))
			}
		}
	}
	return track
}

func main() {
	const WindowTitle = "Path Finder"

	grid := makeGrid()

	rl.InitWindow(WIDTH, HEIGHT, WindowTitle)
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	source_set := false
	target_set := false
	var source *Box = nil
	var target *Box = nil
	var track []*Box = nil

	timer := float32(0.)
	trackIndex := 0

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		timer += rl.GetFrameTime()

		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			m := rl.GetMousePosition()
			selected := &grid[int(m.Y/BOX_DIM)][int(m.X/BOX_DIM)]
			selected.is_source = !selected.is_source
			source_set = true
			source = selected
		}

		if rl.IsMouseButtonPressed(rl.MouseButtonRight) {
			m := rl.GetMousePosition()
			selected := &grid[int(m.Y/BOX_DIM)][int(m.X/BOX_DIM)]
			selected.is_target = !selected.is_target
			target_set = true
			target = selected
		}

		if source_set && target_set && track == nil {
			track = ucs(grid, source, target)
		}

		for _, row := range grid {
			for _, box := range row {
				r := rl.Rectangle{
					X:      float32(box.col * BOX_DIM),
					Y:      float32(box.row * BOX_DIM),
					Width:  float32(BOX_DIM),
					Height: float32(BOX_DIM)}

				if box.is_source {
					rl.DrawRectangleLinesEx(r, 10.0, rl.Blue)
				}
				if box.is_target {
					rl.DrawRectangleLinesEx(r, 10.0, rl.Red)
				}
				rl.DrawRectangleRec(r, colors[box.cost-1])
			}
		}

		if track != nil && trackIndex < len(track) {
			if timer >= 0.2 {
				box := track[trackIndex]
				r := rl.Rectangle{
					X:      float32(box.col * BOX_DIM),
					Y:      float32(box.row * BOX_DIM),
					Width:  float32(BOX_DIM),
					Height: float32(BOX_DIM)}
				rl.DrawRectangleRec(r, rl.Black)
				trackIndex++
				timer = 0.0
			}
		}

		if track != nil && trackIndex >= len(track) {
			r := rl.Rectangle{
				X:      float32(track[len(track)-1].col * BOX_DIM),
				Y:      float32(track[len(track)-1].row * BOX_DIM),
				Width:  float32(BOX_DIM),
				Height: float32(BOX_DIM)}
			rl.DrawRectangleRec(r, rl.Pink)
		}

		rl.EndDrawing()
	}
}
