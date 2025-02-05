package monochrome

import (
	"image"
	"image/color"

	"github.com/colecrouter/gameboy-go/pkg/display"
)

type Display struct {
	initialised bool
	buffer      [display.HEIGHT][display.WIDTH]uint8 // 0-3 for color (0 = white, 3 = black)

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

	for y := 0; y < display.HEIGHT; y++ {
		for x := 0; x < display.WIDTH; x++ {
			var color uint8
			switch d.buffer[y][x] {
			case 0:
				color = 255
			case 1:
				color = 170
			case 2:
				color = 85
			case 3:
				color = 0
			}
			d.image.Set(x, y, Palette[color])
		}
	}
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
		d.buffer[row][x] = color
	}
}

func NewDisplay() *Display {
	d := &Display{initialised: true}
	d.image = image.NewPaletted(image.Rect(0, 0, display.WIDTH, display.HEIGHT), Palette)

	return d
}
