package layers

import (
	"image"

	"github.com/colecrouter/gameboy-go/private/display/monochrome"
	"github.com/colecrouter/gameboy-go/private/memory/registers"
	"github.com/colecrouter/gameboy-go/private/memory/vram"
	"github.com/colecrouter/gameboy-go/private/memory/vram/drawables/tile"
)

// BGLayer represents the background layer.
type BGLayer struct {
	vram        *vram.VRAM
	registers   *registers.Registers
	visibleSize image.Rectangle
}

// NewBGLayer creates a new BGLayer.
func NewBGLayer(v *vram.VRAM, regs *registers.Registers, bounds image.Rectangle) *BGLayer {
	return &BGLayer{
		vram:        v,
		registers:   regs,
		visibleSize: bounds,
	}
}

// Image renders the background for the current registers/VRAM state.
func (bg *BGLayer) Image() image.Image {
	img := image.NewPaletted(bg.visibleSize, monochrome.Palette)
	tSize := tile.TILE_SIZE

	// Determine addressing mode and tile map mode from registers.
	addressingMode := vram.Mode8000
	if !bg.registers.LCDControl.Use8000Method {
		addressingMode = vram.Mode8800
	}
	bgMapMode := vram.TileMapMode(bg.registers.LCDControl.BackgroundUseSecondTileMap)

	// Determine how many tiles to draw. We add one extra tile to account for scrolling.
	cols := bg.visibleSize.Dx()/tSize + 1
	rows := bg.visibleSize.Dy()/tSize + 1

	scrollX := int(bg.registers.ScrollX)
	scrollY := int(bg.registers.ScrollY)

	// Calculate starting tile indexes.
	startTileX := scrollX / tSize
	startTileY := scrollY / tSize

	// For each visible tile, fetch the mapped tile and draw it.
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			// Wrap around if the tilemap is fixed-size (typically 32x32).
			tileX := (startTileX + col) % 32
			tileY := (startTileY + row) % 32

			curTile := bg.vram.GetMappedTile(uint8(tileY), uint8(tileX), bgMapMode, addressingMode)
			if curTile == nil {
				continue
			}

			// Compute the destination position.
			// The (scrollX % tSize) gives the offset within a tile.
			destX := col*tSize - (scrollX % tSize)
			destY := row*tSize - (scrollY % tSize)

			pixels := curTile.Pixels() // Convert the array to a slice.

			// Map pixel values through the background palette.
			for i := range pixels {
				pixels[i] = bg.registers.TilePalette.Match(pixels[i])
			}

			// Draw the tile using our generic opaque blit.
			drawTile(img, pixels, destY, destX)
		}
	}

	return img
}
