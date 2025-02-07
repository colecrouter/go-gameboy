package debug

import (
	"image"
	"image/color"

	"github.com/colecrouter/gameboy-go/pkg/display"
	"github.com/colecrouter/gameboy-go/pkg/memory/vram"
)

const GRID_WIDTH = 192
const GRID_HEIGHT = 128

type TileMenu struct {
	vram    *vram.VRAM
	img     *image.Paletted
	palette *color.Palette
	config  display.Config
}

func NewTileDebug(v *vram.VRAM, p *color.Palette) *TileMenu {
	return &TileMenu{
		vram:    v,
		palette: p,
		config:  display.Config{Title: "Tile Viewer"},
	}
}

func (t *TileMenu) Image() image.Image {
	return t.img
}

func (t *TileMenu) Clock() {
	t.img = image.NewPaletted(image.Rect(0, 0, GRID_WIDTH, GRID_HEIGHT), *t.palette)

	// Construct a grid of tiles
	// Each tile is 8x8 pixels
	for tileY := 0; tileY < 16; tileY++ {
		for tileX := 0; tileX < 24; tileX++ {
			tileIndex := tileY*16 + tileX
			tile := t.vram.GetTile(tileIndex)
			if tile == nil {
				continue
			}
			tileImage := tile.Image()
			if tileImage == nil {
				continue
			}

			for y := 0; y < 8; y++ {
				for x := 0; x < 8; x++ {
					t.img.Set(x+tileX*8, y+tileY*8, tileImage.At(x, y))
				}
			}
		}
	}

	// Draw gridlines between tiles
	gridColor := color.Black
	// Draw horizontal gridlines every 8 pixels
	for y := 8; y < GRID_HEIGHT; y += 8 {
		for x := 0; x < GRID_WIDTH; x++ {
			t.img.Set(x, y, gridColor)
		}
	}
	// Draw vertical gridlines every 8 pixels
	for x := 8; x < GRID_WIDTH; x += 8 {
		for y := 0; y < GRID_HEIGHT; y++ {
			t.img.Set(x, y, gridColor)
		}
	}

}

func (t *TileMenu) Config() *display.Config {
	return &t.config
}
