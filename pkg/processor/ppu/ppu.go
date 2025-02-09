package ppu

import (
	"image"

	"github.com/colecrouter/gameboy-go/pkg/display/monochrome"
	"github.com/colecrouter/gameboy-go/pkg/memory"
	"github.com/colecrouter/gameboy-go/pkg/memory/registers"
	"github.com/colecrouter/gameboy-go/pkg/memory/vram"
	"github.com/colecrouter/gameboy-go/pkg/memory/vram/tile"
)

type PPU struct {
	vram             *vram.VRAM
	oam              *memory.OAM
	registers        *registers.Registers
	lineCycleCounter uint16
	image            *image.Paletted
}

const (
	// Row timings
	oamScanCycles       = 80
	pixelTransferCycles = 172
	hBlankCycles        = 204
	TotalCyclesPerLine  = oamScanCycles + pixelTransferCycles + hBlankCycles

	// Column timings
	visibleLines       = 144
	vBlankLines        = 10
	TotalLinesPerFrame = visibleLines + vBlankLines

	// Helpers
	visibleColumns = 160
)

// NewPPU creates a new PPU instance
func NewPPU(vram *vram.VRAM, oam *memory.OAM, registers *registers.Registers) *PPU {
	return &PPU{
		vram:      vram,
		oam:       oam,
		registers: registers,
		image:     image.NewPaletted(image.Rect(0, 0, visibleColumns, visibleLines), monochrome.Palette),
	}
}

/*
-------┌──────────┐-------
oam    │ transfer │ hblank
80 c   │   172 c  │ 204 c
       │ x 144 l  │
       │          │
-------└──────────┘-------
         vblank
         10 l
*/

// SystemClock emulates a clock cycle on the PPU
func (p *PPU) SystemClock() {
	// Transitions modes
	// See above diagram for reference

	if p.registers.LY >= visibleLines {
		p.registers.LCDStatus.PPUMode = registers.VBlank
	}

	switch p.lineCycleCounter {
	case 0:
		p.registers.LCDStatus.PPUMode = registers.OAMScan
	case oamScanCycles:
		p.registers.LCDStatus.PPUMode = registers.Drawing
	case oamScanCycles + pixelTransferCycles:
		p.registers.LCDStatus.PPUMode = registers.HBlank
	}

	// Handle mode-specific operations & interrupts
	switch p.registers.LCDStatus.PPUMode {
	case registers.OAMScan:
		// p.OAMScan()
		// Is this necessary? Could be an optimization in the future
	case registers.Drawing:
		// Nothing to do here
		// TODO elaborate
	case registers.HBlank:
		// TODO HBlank interrupt
	case registers.VBlank:
		// TODO VBlank interrupt
	}

	// Update the display

	// Handle counter incrementation
	p.lineCycleCounter++
	if p.lineCycleCounter == TotalCyclesPerLine {
		p.registers.LY++
		p.lineCycleCounter = 0
	}
	if p.registers.LY == TotalLinesPerFrame {
		p.registers.LY = 0
	}
}

// DisplayClock updates the image produced by the PPU
func (p *PPU) DisplayClock() {
	addressingMode := vram.Mode8000
	if !p.registers.LCDControl.Use8000Method {
		addressingMode = vram.Mode8800
	}

	// Reset the image
	for i := range p.image.Pix {
		p.image.Pix[i] = 0
	}

	// Draw background layer
	bgMapMode := vram.TileMapMode(p.registers.LCDControl.BackgroundUseSecondaryTileMap)
	for tileY := uint8(0); tileY < visibleLines/tile.TILE_SIZE; tileY++ {
		pixelY := int(tileY)*tile.TILE_SIZE - int(p.registers.ScrollY)

		for tileX := uint8(0); tileX < visibleColumns/tile.TILE_SIZE; tileX++ {
			pixelX := int(tileX)*tile.TILE_SIZE - int(p.registers.ScrollX)
			t := p.vram.GetMappedTile(tileY, tileX, bgMapMode, addressingMode)

			// Issue maybe?
			if t == nil {
				continue
			}

			// Apply color palette
			// Convert to 2D array
			var mapped [tile.TILE_SIZE][tile.TILE_SIZE]uint8
			for y := 0; y < tile.TILE_SIZE; y++ {
				for x := 0; x < tile.TILE_SIZE; x++ {
					mapped[y][x] = p.registers.PaletteData.Match(t.Pixels[y*tile.TILE_SIZE+x])
				}
			}

			// Draw the tile
			p.safeDraw(mapped[:], int(pixelY), int(pixelX))
		}
	}

	// TODO window layer

	// TODO sprite layer

}

func (p *PPU) Image() image.Image {
	return p.image
}

func (p *PPU) safeDraw(pixels [][tile.TILE_SIZE]uint8, y, x int) {
	// Since we're drawing by row, we can have rows that are partially off the screen
	// Clip the row so it doesn't exceed the bounds of the image
	minX := max(0, 0-x)
	maxX := min(tile.TILE_SIZE, (visibleColumns-x)-1)

	for i, row := range pixels {
		// Skip drawing if the row is off the screen
		if y+i < 0 || y+i >= visibleLines {
			continue
		}

		// Copy directly to the image buffer
		copy(p.image.Pix[(y+i)*visibleColumns+x:], row[minX:maxX])
	}
}
