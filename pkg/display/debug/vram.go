package debug

import (
	"image"
	"image/color"

	"github.com/colecrouter/gameboy-go/pkg/memory/vram"
)

type TileDebug struct {
	vram    *vram.VRAM
	img     *image.Paletted
	palette *color.Palette
}

func NewTileDebug(v *vram.VRAM, p *color.Palette) *TileDebug {
	return &TileDebug{vram: v, palette: p}
}

func (t *TileDebug) Image() image.Image {
	return t.img
}

func (t *TileDebug) Clock() {
	t.img = image.NewPaletted(image.Rect(0, 0, 8, 8), *t.palette)

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

}
