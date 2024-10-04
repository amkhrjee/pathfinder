package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"pfinder/algorithm"
	"pfinder/grid"
	"slices"

	rgui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func makeGrid() *grid.Grid {
	r := rand.New(rand.NewSource(42))
	g := grid.Grid{}
	for i, row := range g {
		for j := range row {
			g[i][j] = grid.Box{
				Row:      int32(i),
				Col:      int32(j),
				IsSource: false,
				IsTarget: false,
				Cost:     r.Intn(10-1) + 1,
				Parent:   nil,
			}
		}
	}
	return &g
}

func button(index int) rl.Rectangle {
	return rl.Rectangle{X: float32(12*(index+1)) + float32(index*250), Y: 815, Width: 250, Height: 70}
}

func main() {
	const WindowTitle = "Path Finder"

	g := makeGrid()

	rl.InitWindow(grid.WIDTH, grid.HEIGHT+100, WindowTitle)
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	source_set := false
	target_set := false
	var source *grid.Box = nil
	var target *grid.Box = nil
	var track []*grid.Box = nil
	var final_path []*grid.Box = nil
	var visited []*grid.Box = nil

	timer := float32(0.)
	trackIndex := 0

	rgui.LoadStyle("./style_bluish.rgs")
	rgui.SetStyle(rgui.DEFAULT, rgui.TEXT_SIZE, 20)
	rl.SetWindowIcon(*rl.LoadImage("./assets/windowicon.png"))

	curr_algo := algorithm.ASTAR
	algo_name := "A* Search"

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		timer += rl.GetFrameTime()

		if rgui.Button(button(0), "#211# Reset") || rl.IsKeyPressed(rl.KeyR) {
			if source != nil {
				source.IsSource = false
				source = nil
			}
			if target != nil {
				target.IsTarget = false
				target = nil
			}
			source_set = false
			target_set = false
			track = nil
			final_path = nil
			visited = nil
			timer = float32(0.)
			trackIndex = 0

			for i, row := range g {
				for j := range row {
					g[i][j].Parent = nil
				}
			}
		}

		rgui.Button(button(1), "#142# Settings")

		if rgui.Button(button(2), "#97# "+algo_name) {
			curr_algo = (curr_algo + 1) % algorithm.ALGORITHMS_COUNT
			if curr_algo == 0 {
				curr_algo = 1
			}
			switch curr_algo {
			case algorithm.ASTAR:
				algo_name = "A* Search"
			case algorithm.UCS:
				algo_name = "Unified Cost"
			case algorithm.BFS:
				algo_name = "Breadth First"
			case algorithm.DFS:
				algo_name = "Depth First"
			case algorithm.LBEAM:
				algo_name = "Local Beam"
			}
		}

		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			m := rl.GetMousePosition()
			if m.Y < 800 {
				selected := &g[int(m.Y/grid.BOX_DIM)][int(m.X/grid.BOX_DIM)]
				// setting the source
				if !target_set && !source_set {
					selected.IsSource = true
					source_set = true
					source = selected
				} else {
					if source_set && !target_set && selected != source {
						selected.IsTarget = true
						target_set = true
						target = selected
					}
				}

			}
		}

		if source_set && target_set && track == nil {
			switch curr_algo {
			case algorithm.ASTAR:
				track, final_path = algorithm.AStar(g, source, target)
			case algorithm.UCS:
				track, final_path = algorithm.Ucs(g, source, target)
			case algorithm.BFS:
				track, final_path = algorithm.Bfs(g, source, target)
			case algorithm.DFS:
				track, final_path = algorithm.Dfs(g, source, target)
			case algorithm.LBEAM:
				track, final_path = algorithm.LBeam(g, source, target, 2)
			}
			fmt.Println("Ran: ", algo_name)
		}

		// Draws the grid
		for _, row := range g {
			for _, box := range row {
				r := rl.Rectangle{
					X:      float32(box.Col*grid.BOX_DIM + grid.PADDING),
					Y:      float32(box.Row*grid.BOX_DIM + grid.PADDING),
					Width:  float32(grid.BOX_DIM - 2*grid.PADDING),
					Height: float32(grid.BOX_DIM) - 2*grid.PADDING}

				if box.IsSource {
					rl.DrawRectangleLinesEx(r, 10.0, rl.Blue)
				}
				if box.IsTarget {
					rl.DrawRectangleLinesEx(r, 10.0, rl.Red)
				}
				rl.DrawRectangleRec(r, grid.Colors[box.Cost-1])
			}
		}

		if track != nil && trackIndex < len(track) {
			if timer >= 0.05 {
				// if timer >= 0.02 {
				box := track[trackIndex]
				if !slices.Contains(visited, box) {
					visited = append(visited, box)
				}
				r := rl.Rectangle{
					X:      float32(box.Col * grid.BOX_DIM),
					Y:      float32(box.Row * grid.BOX_DIM),
					Width:  float32(grid.BOX_DIM),
					Height: float32(grid.BOX_DIM)}
				rl.DrawRectangleRec(r, rl.Blue)
				trackIndex++

				timer = 0.0
			}
			for _, t := range track[:trackIndex-1] {
				r := rl.Rectangle{
					X:      float32(t.Col*grid.BOX_DIM + grid.PADDING),
					Y:      float32(t.Row*grid.BOX_DIM + grid.PADDING),
					Width:  float32(grid.BOX_DIM - 2*grid.PADDING),
					Height: float32(grid.BOX_DIM - 2*grid.PADDING)}
				rl.DrawRectangleRec(r, rl.Pink)

				if t == source {
					rl.DrawRectangleRec(r, rl.Red)
				}

				if slices.Contains(final_path, t) {
					x := float32(t.Col*grid.BOX_DIM + grid.PADDING)
					y := float32(t.Row*grid.BOX_DIM + grid.PADDING)
					r := rl.Rectangle{
						X:      x,
						Y:      y,
						Width:  float32(grid.BOX_DIM - 2*grid.PADDING),
						Height: float32(grid.BOX_DIM - 2*grid.PADDING)}
					rl.DrawRectangleRec(r, rl.Red)

				}
			}
		}
		if track != nil && trackIndex >= len(track) {

			for _, f := range final_path {
				r := rl.Rectangle{
					X:      float32(f.Col*grid.BOX_DIM + grid.PADDING),
					Y:      float32(f.Row*grid.BOX_DIM + grid.PADDING),
					Width:  float32(grid.BOX_DIM - 2*grid.PADDING),
					Height: float32(grid.BOX_DIM) - 2*grid.PADDING}
				rl.DrawRectangleRec(r, rl.Red)
			}

			r := rl.Rectangle{
				X:      float32(track[len(track)-1].Col*grid.BOX_DIM + grid.PADDING),
				Y:      float32(track[len(track)-1].Row*grid.BOX_DIM + grid.PADDING),
				Width:  float32(grid.BOX_DIM - 2*grid.PADDING),
				Height: float32(grid.BOX_DIM - 2*grid.PADDING)}
			rl.DrawRectangleRec(r, rl.Green)
		}

		if track != nil {
			l := rl.Rectangle{
				X:      0,
				Y:      750,
				Width:  150,
				Height: 50,
			}
			rl.DrawRectangleRec(l, color.RGBA{rl.RayWhite.R, rl.RayWhite.G, rl.RayWhite.B, 200})
			rgui.Label(l,
				fmt.Sprintf(" Visited: %d", len(visited)),
			)
		}

		rl.EndDrawing()
	}
}
