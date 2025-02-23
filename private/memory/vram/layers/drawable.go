package layers

import (
	"image"

	"github.com/colecrouter/gameboy-go/private/memory/vram/drawables/tile"
)

type Layer interface {
	Image() image.Image
}

func drawTile(dst *image.Paletted, pixels []uint8, startY, startX int) {
	bounds := dst.Bounds()

	width := tile.TILE_SIZE
	height := tile.TILE_SIZE

	// Fast path: if the entire block is within bounds, copy row by row.
	if startX >= bounds.Min.X && startX+width <= bounds.Max.X &&
		startY >= bounds.Min.Y && startY+height <= bounds.Max.Y {
		for y := 0; y < height; y++ {
			destIdx := dst.PixOffset(startX, startY+y)
			srcStart := y * width
			copy(dst.Pix[destIdx:destIdx+width], pixels[srcStart:srcStart+width])
		}
		return
	}

	// Otherwise, copy pixel by pixel.
	for y := 0; y < height; y++ {
		destY := startY + y
		if destY < bounds.Min.Y || destY >= bounds.Max.Y {
			continue
		}
		for x := 0; x < width; x++ {
			destX := startX + x
			if destX < bounds.Min.X || destX >= bounds.Max.X {
				continue
			}
			idx := dst.PixOffset(destX, destY)
			dst.Pix[idx] = pixels[y*width+x]
		}
	}
}

// drawSprite copies a sprite block (which can be 8x8 or 8x16) while skipping pixels matching the
// transparent color index. The block is provided as a flat []uint8 and its dimensions via height.
func drawSprite(dst *image.Paletted, pixels []uint8, startY, startX int, doubleHeight bool, transparent uint8) {
	bounds := dst.Bounds()

	width := tile.TILE_SIZE
	height := tile.TILE_SIZE
	if doubleHeight {
		height *= 2
	}

	for y := 0; y < height; y++ {
		destY := startY + y
		if destY < bounds.Min.Y || destY >= bounds.Max.Y {
			continue
		}
		for x := 0; x < width; x++ {
			destX := startX + x
			if destX < bounds.Min.X || destX >= bounds.Max.X {
				continue
			}
			pix := pixels[y*width+x]
			if pix == transparent {
				continue
			}
			idx := dst.PixOffset(destX, destY)
			dst.Pix[idx] = pix
		}
	}
}
