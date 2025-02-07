package lcd

import (
	"image"

	"github.com/colecrouter/gameboy-go/pkg/display"
	"github.com/colecrouter/gameboy-go/pkg/display/monochrome"
)

const WIDTH = 160
const HEIGHT = 144

type Display struct {
	initialised bool

	image  *image.Paletted
	config display.Config
}

// Clock updates the image on the display
func (d *Display) Clock() {
	if !d.initialised {
		panic("Display not initialised")
	}

	// Do nothing
}

func (d *Display) Image() image.Image {
	if !d.initialised {
		panic("Display not initialised")
	}

	return d.image
}

func (d *Display) DrawScanline(row uint8, line []uint8) {
	if !d.initialised {
		panic("Display not initialised")
	}

	for x := 0; x < WIDTH; x++ {
		color := line[x]
		d.image.Set(int(x), int(row), monochrome.Palette[color])
	}
}

func (d *Display) Config() *display.Config {
	return &d.config
}

func NewDisplay() *Display {
	d := &Display{initialised: true}
	d.image = image.NewPaletted(image.Rect(0, 0, WIDTH, HEIGHT), monochrome.Palette)
	d.config = display.Config{Title: "Display"}

	return d
}
