package tile

import (
	"image"

	"github.com/colecrouter/gameboy-go/pkg/display/monochrome"
)

const TILE_SIZE = 8

type Tile struct {
	initialized bool
	Pixels      [TILE_SIZE * TILE_SIZE]uint8
}

func FromBytes(bytes [16]uint8) *Tile {
	t := &Tile{initialized: true}

	// Iterate per row
	// Each pair of bytes represents a row of 8 pixels
	for i := uint8(0); i < TILE_SIZE; i++ {
		// Iterate per column
		for j := uint8(0); j < TILE_SIZE; j++ {
			msb := (bytes[i*2] >> (7 - j)) & 1
			lsb := (bytes[i*2+1] >> (7 - j)) & 1
			t.Pixels[i*TILE_SIZE+j] = msb | (lsb << 1)
		}
	}

	return t
}

func (t *Tile) Image() *image.Paletted {
	img := image.NewPaletted(image.Rect(0, 0, TILE_SIZE, TILE_SIZE), monochrome.Palette)

	copy(img.Pix[:], t.Pixels[:])

	return img
}
