package monochrome

import (
	"image"
	"image/color"

	"github.com/colecrouter/gameboy-go/pkg/display"
)

type Display struct {
	initialised bool

	image *image.Paletted
}

var Palette = []color.Color{
	color.RGBA{255, 255, 255, 255},
	color.RGBA{170, 170, 170, 255},
	color.RGBA{85, 85, 85, 255},
	color.Black,
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

	for x := 0; x < display.WIDTH; x++ {
		color := line[x]
		d.image.Set(int(x), int(row), Palette[color])
	}
}

func NewDisplay() *Display {
	d := &Display{initialised: true}
	d.image = image.NewPaletted(image.Rect(0, 0, display.WIDTH, display.HEIGHT), Palette)

	return d
}
