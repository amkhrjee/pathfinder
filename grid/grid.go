package grid

const WIDTH, HEIGHT = 800, 800
const BOX_DIM = 80

type Box struct {
	row       int32
	col       int32
	is_source bool
	is_target bool
}

type Grid [HEIGHT / BOX_DIM][WIDTH / BOX_DIM]Box
