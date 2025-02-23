package layers

import (
	"image"

	"github.com/colecrouter/gameboy-go/private/display/monochrome"
	"github.com/colecrouter/gameboy-go/private/memory"
	"github.com/colecrouter/gameboy-go/private/memory/registers"
	"github.com/colecrouter/gameboy-go/private/memory/vram"
	"github.com/colecrouter/gameboy-go/private/memory/vram/drawables/tile"
)

// transformPixels applies horizontal and/or vertical flip to a flat pixel array.
// `pixels` should be in row-major order with dimensions (width x height).
func transformPixels(pixels []uint8, width, height int, flipX, flipY bool) []uint8 {
	transformed := make([]uint8, len(pixels))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			srcX, srcY := x, y
			if flipX {
				srcX = width - 1 - x
			}
			if flipY {
				srcY = height - 1 - y
			}
			transformed[y*width+x] = pixels[srcY*width+srcX]
		}
	}
	return transformed
}

// SpriteLayer renders the sprite layer.
type SpriteLayer struct {
	oam         *memory.OAM
	vram        *vram.VRAM
	registers   *registers.Registers
	visibleSize image.Rectangle
}

// NewSpriteLayer creates a new SpriteLayer.
func NewSpriteLayer(o *memory.OAM, v *vram.VRAM, regs *registers.Registers, bounds image.Rectangle) *SpriteLayer {
	return &SpriteLayer{
		oam:         o,
		vram:        v,
		registers:   regs,
		visibleSize: bounds,
	}
}

// Image renders all sprites into an image.
func (s *SpriteLayer) Image() image.Image {
	img := image.NewPaletted(s.visibleSize, monochrome.Palette)
	width := tile.TILE_SIZE
	height := tile.TILE_SIZE

	// Use 8x16 mode if set in LCDControl.
	use8x16 := s.registers.LCDControl.Sprites8x16
	if use8x16 {
		height = tile.TILE_SIZE * 2
	}

	// Loop through all 40 sprites.
	for i := 0; i < 40; i++ {
		spr := s.oam.ReadSprite(i)
		if spr == nil {
			continue
		}

		// Adjust for hardware offsets (sprites are offset by 16 in Y and 8 in X).
		posX := int(spr.X())
		posY := int(spr.Y())
		// Skip sprites off-screen.
		if posX >= s.visibleSize.Dx() || posY >= s.visibleSize.Dy() {
			continue
		}

		pixels := spr.Pixels()

		// Apply sprite flip transformations if needed.
		if spr.FlipX || spr.FlipY {
			pixels = transformPixels(pixels, width, height, spr.FlipX, spr.FlipY)
		}

		// Map pixel values through the sprite palette.
		palette := s.registers.ObjectPalletes[spr.DMGPalette]
		for i := range pixels {
			pixels[i] = palette.Match(pixels[i])
		}

		// Draw the sprite onto the frame.
		// Assuming that transparent pixels have a value of 0.
		drawSprite(img, pixels, posY, posX, use8x16, 0)
	}
	return img
}
