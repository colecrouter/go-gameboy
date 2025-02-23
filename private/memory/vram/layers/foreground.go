package layers

import (
	"image"

	"github.com/colecrouter/gameboy-go/private/display/monochrome"
	"github.com/colecrouter/gameboy-go/private/memory/registers"
	"github.com/colecrouter/gameboy-go/private/memory/vram"
	"github.com/colecrouter/gameboy-go/private/memory/vram/drawables/tile"
)

// WindowLayer represents the window (foreground) layer.
type WindowLayer struct {
	vram        *vram.VRAM
	registers   *registers.Registers
	visibleSize image.Rectangle
}

// NewWindowLayer creates a new WindowLayer.
func NewWindowLayer(v *vram.VRAM, regs *registers.Registers, bounds image.Rectangle) *WindowLayer {
	return &WindowLayer{
		vram:        v,
		registers:   regs,
		visibleSize: bounds,
	}
}

// Image renders the window layer according to the window registers.
func (w *WindowLayer) Image() image.Image {
	img := image.NewPaletted(w.visibleSize, monochrome.Palette)
	tSize := tile.TILE_SIZE

	// Determine addressing mode, similar to the BG layer.
	addressingMode := vram.Mode8000
	if !w.registers.LCDControl.Use8000Method {
		addressingMode = vram.Mode8800
	}
	// Use the window tilemap based on the LCD control register.
	winMapMode := vram.TileMapMode(w.registers.LCDControl.WindowUseSecondTileMap)

	// Window position: WindowY is the scanline where the window begins
	// WindowX is offset by 7 (as per hardware spec).
	winY := int(w.registers.WindowY)
	winX := int(w.registers.WindowX) - 7

	// Calculate how many tiles we need to draw to cover the window area.
	cols := w.visibleSize.Dx()/tSize + 1
	rows := w.visibleSize.Dy()/tSize + 1

	// For each window tile, fetch and draw.
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			// The window tilemap is addressed directly (row & col are relative to the window).
			curTile := w.vram.GetMappedTile(uint8(row), uint8(col), winMapMode, addressingMode)
			if curTile == nil {
				continue
			}

			// Compute destination position in the window image.
			destX := col*tSize + winX
			destY := row*tSize + winY

			// If the tile is completely off-screen, skip it.
			if destX >= w.visibleSize.Max.X || destY >= w.visibleSize.Max.Y {
				continue
			}

			pixels := curTile.Pixels()

			// Map pixel values through the background palette.
			for i := range pixels {
				pixels[i] = w.registers.TilePalette.Match(pixels[i])
			}

			// Draw the tile using our generic opaque blit.
			drawTile(img, pixels, destY, destX)
		}
	}

	return img
}
