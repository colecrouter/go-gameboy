package display

import "image"

const WIDTH = 160
const HEIGHT = 144

type Display interface {
	Image() image.Image
	Clock()
}
