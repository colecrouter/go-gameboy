package lcd

import (
	"image"

	"github.com/colecrouter/gameboy-go/pkg/display"
	"github.com/colecrouter/gameboy-go/pkg/processor/ppu"
)

const WIDTH = 160
const HEIGHT = 144

type Display struct {
	initialised bool
	ppu         *ppu.PPU
	config      display.Config
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

	return d.ppu.Image()
}

func (d *Display) Config() *display.Config {
	return &d.config
}

func NewDisplay(ppu *ppu.PPU) *Display {
	d := &Display{initialised: true}
	d.ppu = ppu
	d.config = display.Config{Title: "Display"}

	return d
}
