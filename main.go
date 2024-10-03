package main

import (
	"math/rand"
	"pfinder/algorithm"
	"pfinder/grid"

	rgui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func makeGrid() *grid.Grid {
	g := grid.Grid{}
	for i, row := range g {
		for j := range row {
			g[i][j] = grid.Box{
				Row:      int32(i),
				Col:      int32(j),
				IsSource: false,
				IsTarget: false,
				Cost:     rand.Intn(10-1) + 1,
				Parent:   nil,
			}
		}
	}
	return &g
}

func button(index int) rl.Rectangle {
	return rl.Rectangle{X: float32(10*(index+1)) + float32(index*250), Y: 815, Width: 250, Height: 70}
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

	timer := float32(0.)
	trackIndex := 0

	rgui.LoadStyle("./style_bluish.rgs")
	rgui.SetStyle(rgui.DEFAULT, rgui.TEXT_SIZE, 20)
	rl.SetWindowIcon(*rl.LoadImage("./assets/windowicon.png"))
	is_astar := true
	algo_name := "A* Search"

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		rgui.Label(
			rl.Rectangle{
				X:      float32(50*(1+1)) + float32(1*250),
				Y:      815,
				Width:  100,
				Height: 70,
			},
			"@amkhrjee",
		)

		timer += rl.GetFrameTime()

		if rgui.Button(button(0), "#211#Reset") || rl.IsKeyPressed(rl.KeyR) {
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
			timer = float32(0.)
			trackIndex = 0
			for i, row := range g {
				for j := range row {
					g[i][j].Parent = nil
				}
			}
		}
		if rgui.Button(button(2), "#62#"+algo_name) {
			is_astar = !is_astar
			if is_astar {
				algo_name = "A* Search"
			} else {
				algo_name = "Djikstra UCS"
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
					if source_set && !target_set {
						selected.IsTarget = true
						target_set = true
						target = selected
					}
				}

			}
		}

		if source_set && target_set && track == nil {
			if is_astar {
				track, final_path = algorithm.AStar(g, source, target)
			} else {
				track, final_path = algorithm.Ucs(g, source, target)
			}
		}

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
				box := track[trackIndex]
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
					X:      float32(t.Col * grid.BOX_DIM),
					Y:      float32(t.Row * grid.BOX_DIM),
					Width:  float32(grid.BOX_DIM),
					Height: float32(grid.BOX_DIM)}
				rl.DrawRectangleRec(r, rl.Pink)
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

		rl.EndDrawing()
	}
}
