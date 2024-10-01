package main

import (
	"fmt"
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

	message := "Visualise search algorithms in a fun way!"

	showMessage := true

	rgui.LoadStyle("./style_bluish.rgs")
	rgui.SetStyle(rgui.DEFAULT, rgui.TEXT_SIZE, 20)
	// rgui.SetStyle(rgui.icn, rgui., 20)
	rl.SetWindowIcon(*rl.LoadImage("./assets/windowicon.png"))
	algo_name := "A* Search"

	message_box := rl.Rectangle{X: 150, Y: 250, Height: 250, Width: 500}
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		timer += rl.GetFrameTime()

		rgui.Button(button(0), "#72#Reset")
		rgui.Button(button(1), "#90#Obstacles")
		rgui.Button(button(2), fmt.Sprintf("#62#%s", algo_name))

		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) && !showMessage {
			m := rl.GetMousePosition()
			selected := &g[int(m.Y/grid.BOX_DIM)][int(m.X/grid.BOX_DIM)]
			selected.IsSource = !selected.IsSource
			source_set = true
			source = selected
		}

		if rl.IsMouseButtonPressed(rl.MouseButtonRight) && !showMessage {
			m := rl.GetMousePosition()
			selected := &g[int(m.Y/grid.BOX_DIM)][int(m.X/grid.BOX_DIM)]
			selected.IsTarget = !selected.IsTarget
			target_set = true
			target = selected
		}

		if source_set && target_set && track == nil {
			track, final_path = algorithm.AStar(g, source, target)
		}

		for _, row := range g {
			for _, box := range row {
				r := rl.Rectangle{
					X:      float32(box.Col * grid.BOX_DIM),
					Y:      float32(box.Row * grid.BOX_DIM),
					Width:  float32(grid.BOX_DIM),
					Height: float32(grid.BOX_DIM)}

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
				rl.DrawRectangleRec(r, rl.Green)
			}
		}
		if track != nil && trackIndex >= len(track) {

			for _, f := range final_path {
				r := rl.Rectangle{
					X:      float32(f.Col * grid.BOX_DIM),
					Y:      float32(f.Row * grid.BOX_DIM),
					Width:  float32(grid.BOX_DIM),
					Height: float32(grid.BOX_DIM)}
				rl.DrawRectangleRec(r, rl.Red)
			}

			r := rl.Rectangle{
				X:      float32(track[len(track)-1].Col * grid.BOX_DIM),
				Y:      float32(track[len(track)-1].Row * grid.BOX_DIM),
				Width:  float32(grid.BOX_DIM),
				Height: float32(grid.BOX_DIM)}
			rl.DrawRectangleRec(r, rl.Pink)

		}

		if showMessage {
			rl.DrawRectangleRec(message_box, rl.RayWhite)
			res := rgui.MessageBox(message_box, "#191#Pathfinder v0.1.1 October 2024", message, "Cool!")
			if res >= 0 {
				showMessage = false
			}
		}

		rl.EndDrawing()
	}
}
