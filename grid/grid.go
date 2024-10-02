package grid

import rl "github.com/gen2brain/raylib-go/raylib"

const WIDTH, HEIGHT = 800, 800
const BOX_DIM = 40
const PADDING = 4

type Box struct {
	Row      int32
	Col      int32
	Cost     int
	IsSource bool
	IsTarget bool
	Parent   *Box
}

const ROWS = HEIGHT / BOX_DIM
const COLS = WIDTH / BOX_DIM

type Grid [ROWS][COLS]Box

var Colors = []rl.Color{
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
