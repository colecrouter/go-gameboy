package ppu

import (
	// added for debugging

	"github.com/colecrouter/gameboy-go/pkg/display"
	"github.com/colecrouter/gameboy-go/pkg/memory"
	"github.com/colecrouter/gameboy-go/pkg/memory/registers"
	"github.com/colecrouter/gameboy-go/pkg/memory/vram"
)

type PPU struct {
	vram             *vram.VRAM
	oam              *memory.OAM
	display          display.Display
	registers        *registers.Registers
	lineCycleCounter uint16
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
func NewPPU(vram *vram.VRAM, oam *memory.OAM, display display.Display, registers *registers.Registers) *PPU {
	return &PPU{
		vram:      vram,
		oam:       oam,
		display:   display,
		registers: registers,
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

// Clock emulates a clock cycle on the PPU
func (p *PPU) Clock() {
	// Transitions modes
	// See above diagram for reference
	switch p.lineCycleCounter {
	case 0:
		p.registers.LCDStatus.PPUMode = registers.OAMScan
	case oamScanCycles:
		p.registers.LCDStatus.PPUMode = registers.Drawing
	case oamScanCycles + pixelTransferCycles:
		p.registers.LCDStatus.PPUMode = registers.HBlank
	}
	if p.registers.LY >= visibleLines {
		p.registers.LCDStatus.PPUMode = registers.VBlank
	}

	// Handle mode-specific operations & interrupts
	switch p.registers.LCDStatus.PPUMode {
	case registers.OAMScan:
		// p.OAMScan()
		// Is this necessary? Could be an optimization in the future
	case registers.Drawing:
		line := p.getScanline()
		p.display.DrawScanline(p.registers.LY, line)
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

func (p *PPU) getScanline() []uint8 {
	var scanline [160]byte

	addressingMode := vram.Mode8000
	if !p.registers.LCDControl.Use8000Method {
		addressingMode = vram.Mode8800
	}

	// Draw background layer
	bgMapMode := vram.TileMapMode(p.registers.LCDControl.BackgroundUseSecondaryTileMap)
	for pixelX := uint8(0); pixelX < visibleColumns; pixelX++ {
		scrolledY := p.registers.LY + p.registers.ScrollY
		scrolledX := p.registers.ScrollX + pixelX

		tile := p.vram.ReadMappedTileAt(scrolledX, scrolledY, bgMapMode, addressingMode)
		if tile == nil {
			continue
		}

		// Allow for scrolling within the tile
		tileColor := tile.Pixels[scrolledY%8][scrolledX%8]

		// Map the color to the display palette
		mapped := p.matchColorPalette(tileColor)

		// Draw the pixel
		scanline[pixelX] = mapped
	}

	// TODO window layer

	// TODO sprite layer

	return scanline[:]
}

func (p *PPU) matchColorPalette(color uint8) uint8 {
	return p.registers.PaletteData.Match(color)
}
