package main

import (
	"fmt"
	"math"
	"math/rand"
	"slices"

	rl "github.com/gen2brain/raylib-go/raylib"
	pq "gopkg.in/dnaeon/go-priorityqueue.v1"
)

const WIDTH, HEIGHT = 800, 800
const BOX_DIM = 40

type Box struct {
	row             int32
	col             int32
	cost            int
	is_source       bool
	is_target       bool
	cumulative_cost float64
	parent          *Box
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
				row:             int32(i),
				col:             int32(j),
				is_source:       false,
				is_target:       false,
				cost:            rand.Intn(10-1) + 1,
				cumulative_cost: math.Inf(0),
				parent:          nil,
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
func ucs(grid *Grid, start *Box, target *Box) ([]*Box, []*Box) {
	root := start
	q := pq.New[*Box, float64](pq.MinHeap)
	track := make([]*Box, 0)
	final_path := make([]*Box, 0)
	// inserting the initial nodes
	for _, n := range neighbors(grid, root) {
		relative_cost := math.Abs(float64(n.cost - root.cost))
		q.Put(n, relative_cost)
		n.cumulative_cost = relative_cost
	}
	for !q.IsEmpty() {
		curr := q.Get()
		track = append(track, curr.Value)
		if curr.Value == target {
			break
		}

		for _, n := range neighbors(grid, curr.Value) {
			if !slices.Contains(track, n) && n != start {
				total_cost := curr.Priority + math.Abs(float64(n.cost-curr.Value.cost))
				q.Put(n, total_cost)
				n.cumulative_cost = total_cost
				n.parent = curr.Value
			}
		}
	}

	fmt.Printf("Cost of start: %1.f\n", start.cumulative_cost)
	fmt.Printf("Cost of target: %1.f\n", target.cumulative_cost)

	// backtracking the path
	curr := target
	for curr != nil {
		final_path = append(final_path, curr)
		curr = curr.parent
	}

	fmt.Println("Length of final path: ", len(final_path))

	return track, final_path
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
	var final_path []*Box = nil

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
			source.cumulative_cost = 0
		}

		if rl.IsMouseButtonPressed(rl.MouseButtonRight) {
			m := rl.GetMousePosition()
			selected := &grid[int(m.Y/BOX_DIM)][int(m.X/BOX_DIM)]
			selected.is_target = !selected.is_target
			selected.cumulative_cost = 0
			target_set = true
			target = selected
		}

		if source_set && target_set && track == nil {
			track, final_path = ucs(grid, source, target)
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
			if timer >= 0.1 {
				box := track[trackIndex]
				r := rl.Rectangle{
					X:      float32(box.col * BOX_DIM),
					Y:      float32(box.row * BOX_DIM),
					Width:  float32(BOX_DIM),
					Height: float32(BOX_DIM)}
				rl.DrawRectangleRec(r, rl.Green)
				trackIndex++

				timer = 0.0
			}
			for _, t := range track[:trackIndex] {
				r := rl.Rectangle{
					X:      float32(t.col * BOX_DIM),
					Y:      float32(t.row * BOX_DIM),
					Width:  float32(BOX_DIM),
					Height: float32(BOX_DIM)}
				rl.DrawRectangleRec(r, rl.Green)
			}
		}
		if track != nil && trackIndex >= len(track) {

			for _, f := range final_path {
				r := rl.Rectangle{
					X:      float32(f.col * BOX_DIM),
					Y:      float32(f.row * BOX_DIM),
					Width:  float32(BOX_DIM),
					Height: float32(BOX_DIM)}
				rl.DrawRectangleRec(r, rl.Red)
			}

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
