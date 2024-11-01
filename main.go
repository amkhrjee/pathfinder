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
	r := rand.New(rand.NewSource(7))
	g := grid.Grid{}
	for i, row := range g {
		for j := range row {
			g[i][j] = grid.Box{
				Row:        float64(i),
				Col:        float64(j),
				IsSource:   false,
				IsTarget:   false,
				Cost:       float64(r.Intn(10-1) + 1),
				Parent:     nil,
				Visited:    false,
				IsObstacle: false,
			}
		}
	}
	return &g
}

const (
	LIGHT int32 = 0
	DARK  int32 = 1
)

func button(index int) rl.Rectangle {
	return rl.Rectangle{X: float32(12*(index+1)) + float32(index*250), Y: 815, Width: 250, Height: 70}
}

func rectangle(x int, y int) rl.Rectangle {
	xPadding := float32(10)
	yOffset := float32(40)
	return rl.Rectangle{
		X:      200 + float32(x)*200 + xPadding,
		Y:      200 + float32(y+1)*yOffset + 20,
		Width:  150,
		Height: 20,
	}
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

	isSettingsOpen := false

	timer := float32(0.)
	trackIndex := 0

	rgui.LoadStyle("./style_bluish.rgs")
	rgui.SetStyle(rgui.DEFAULT, rgui.TEXT_SIZE, 20)
	rl.SetWindowIcon(*rl.LoadImage("./assets/windowicon.png"))

	curr_algo := algorithm.ASTAR
	algo_name := "A* Search"

	kValue := int32(2)
	kvalueEditMode := false

	speedValue := int32(1)
	speedEditMode := false

	theme := LIGHT

	background_color := rl.RayWhite

	infoLabel := rl.Rectangle{
		X:      0,
		Y:      750,
		Width:  150,
		Height: 50,
	}

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		if theme == LIGHT {
			background_color = rl.RayWhite
		} else {
			background_color = rl.Black
		}
		rl.ClearBackground(background_color)

		timer += rl.GetFrameTime()

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
				if box.IsObstacle {
					if theme == DARK {
						rl.DrawRectangleRec(r, rl.LightGray)
					} else {
						rl.DrawRectangleRec(r, rl.DarkGray)
					}
				} else {
					rl.DrawRectangleRec(r, grid.Colors[int(box.Cost-1)])
				}
			}
		}

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
					g[i][j].Visited = false
				}
			}
		}

		if rgui.Button(button(1), "#142# Settings") {
			isSettingsOpen = !isSettingsOpen
		}

		if rgui.Button(button(2), "#97# "+algo_name) {
			curr_algo = (curr_algo + 1) % algorithm.ALGORITHMS_COUNT
			if curr_algo == 0 {
				curr_algo = 1
			}
			switch curr_algo {
			case algorithm.ASTAR:
				algo_name = "A* Search"
			case algorithm.UCS:
				algo_name = "Uniform Cost"
			case algorithm.BFS:
				algo_name = "Breadth First"
			case algorithm.DFS:
				algo_name = "Depth First"
			case algorithm.LBEAM:
				algo_name = "Local Beam"
			}
		}

		if !isSettingsOpen {

			if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
				m := rl.GetMousePosition()
				if m.Y < 800 {
					selected := &g[int(m.Y/grid.BOX_DIM)][int(m.X/grid.BOX_DIM)]
					// setting the source
					if !target_set && !source_set {
						selected.IsSource = true
						selected.Visited = true
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

			if rl.IsMouseButtonPressed((rl.MouseButtonRight)) {
				m := rl.GetMousePosition()
				if m.Y < 800 {
					selected := &g[int(m.Y/grid.BOX_DIM)][int(m.X/grid.BOX_DIM)]
					// setting the source
					selected.IsObstacle = !selected.IsObstacle
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
					track, final_path = algorithm.LBeam(g, source, target, int(kValue))
				}
			}

			if track != nil && trackIndex < len(track) {
				if timer >= float32(0.05)/float32(speedValue) {
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

				last := track[len(track)-1]
				r := rl.Rectangle{
					X:      float32(last.Col*grid.BOX_DIM + grid.PADDING),
					Y:      float32(last.Row*grid.BOX_DIM + grid.PADDING),
					Width:  float32(grid.BOX_DIM - 2*grid.PADDING),
					Height: float32(grid.BOX_DIM - 2*grid.PADDING)}
				rl.DrawRectangleRec(r, rl.Green)
				if last != target {
					rl.DrawText(
						"X",
						int32(last.Col*grid.BOX_DIM+grid.PADDING+5),
						int32(last.Row*grid.BOX_DIM+grid.PADDING),
						35,
						rl.Red)
				}
			}

			if track != nil {

				rl.DrawRectangleRec(infoLabel, color.RGBA{background_color.R, background_color.G, background_color.B, 200})
				rgui.Label(infoLabel,
					fmt.Sprintf("  Visited: %d", len(visited)),
				)
			}
		}

		if isSettingsOpen {
			r := rl.Rectangle{
				X:      200,
				Y:      200,
				Width:  400,
				Height: 400,
			}
			isClosePressed := rgui.WindowBox(r, "#142# Settings")
			if isClosePressed {
				isSettingsOpen = false
			}
			rgui.Label(rectangle(0, 0), "Local Beam K")
			res := rgui.Spinner(rectangle(1, 0), "", &kValue, 1, 4, kvalueEditMode)
			if res < 1 || res > 4 {
				kvalueEditMode = !kvalueEditMode
			}

			rgui.Label(rectangle(0, 1), "Speedup (x)")
			speed := rgui.Spinner(rectangle(1, 1), "", &speedValue, 1, 4, speedEditMode)
			if speed < 1 || speed > 4 {
				speedEditMode = !speedEditMode
			}

			n := rectangle(0, 2)
			theme = rgui.ToggleGroup(rl.Rectangle{
				X:      n.X + 40,
				Y:      n.Y,
				Height: 40,
				Width:  150,
			}, "Light Mode;Dark Mode", theme)

			beg := 4
			f := rectangle(0, beg+1)
			five := rectangle(0, beg+2)
			s := rectangle(0, beg+3)
			rgui.Label(rl.Rectangle{
				X:      f.X + 110,
				Y:      f.Y,
				Width:  200,
				Height: 20,
			}, "Pathfinder v0.1.2")
			rgui.Label(rl.Rectangle{
				X:      five.X + 120,
				Y:      five.Y,
				Width:  200,
				Height: 20,
			}, "by @amkhrjee")
			rgui.Label(rl.Rectangle{
				X:      s.X + 60,
				Y:      s.Y,
				Width:  300,
				Height: 20,
			}, "made at Tezpur University")

			// Closing the window
			if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
				m := rl.GetMousePosition()
				if (m.X > 600 || m.X < 200) || (m.Y > 600 || m.Y < 200) {
					isSettingsOpen = false
				}
			}

		}

		rl.EndDrawing()
	}
}
