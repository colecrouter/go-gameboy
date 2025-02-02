package ppu

import (
	"github.com/colecrouter/gameboy-go/pkg/display"
	"github.com/colecrouter/gameboy-go/pkg/memory"
	"github.com/colecrouter/gameboy-go/pkg/memory/registers"
	"github.com/colecrouter/gameboy-go/pkg/memory/vram"
	"github.com/colecrouter/gameboy-go/pkg/memory/vram/sprite"
	"github.com/colecrouter/gameboy-go/pkg/memory/vram/tile"
)

type PPU struct {
	vram             *vram.VRAM
	oam              *memory.OAM
	display          display.Display
	registers        *registers.Registers
	lineCounter      uint8
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
		p.registers.LCDStatus.Status.PPUMode = registers.OAMScan
	case oamScanCycles:
		p.registers.LCDStatus.Status.PPUMode = registers.Drawing
	case oamScanCycles + pixelTransferCycles:
		p.registers.LCDStatus.Status.PPUMode = registers.HBlank
	}
	if p.lineCounter >= visibleLines {
		p.registers.LCDStatus.Status.PPUMode = registers.VBlank
	}

	// Handle mode-specific operations
	switch p.registers.LCDStatus.Status.PPUMode {
	case registers.OAMScan:
		// p.OAMScan()
	case registers.Drawing:
		// p.DrawScanline()
		line := p.getScanline()
		p.display.DrawScanline(p.lineCounter, line)
	case registers.HBlank:
		// p.HBlank()
	case registers.VBlank:
		// p.VBlank()
	}

	// Update the display

	// Handle counter incrementation
	p.lineCycleCounter++
	if p.lineCycleCounter == TotalCyclesPerLine {
		p.lineCounter++
		p.lineCycleCounter = 0
	}
	if p.lineCounter == TotalLinesPerFrame {
		p.lineCounter = 0
	}
}

func (p *PPU) getScanline() []uint8 {
	horizontalTiles := uint8(visibleColumns / tile.TILE_SIZE)

	var scanline [160]byte

	// Draw background
	for pixelX := uint8(0); pixelX < visibleColumns; pixelX++ {
		scrolledY := (p.registers.LCDStatus.ScrollY + p.lineCounter)
		scrolledX := (p.registers.LCDStatus.ScrollX + pixelX)

		tile := p.vram.ReadMappedTile(0, uint8(scrolledY/8)*horizontalTiles+uint8(scrolledX/8))
		tileColor := tile.ReadPixel(scrolledY%8, scrolledX%8)

		if tileColor == 0 {
			continue
		}

		scanline[pixelX] = p.matchColorPalette(tileColor)
	}

	// Draw window
	for pixelX := uint8(0); pixelX < visibleColumns; pixelX++ {
		positionedY := p.registers.LCDStatus.PositionY + p.lineCounter
		positionedX := p.registers.LCDStatus.PositionX + pixelX

		// If window pixel is outside display bounds, skip drawing it
		if positionedX >= visibleColumns {
			continue
		}

		tile := p.vram.ReadMappedTile(0, uint8(positionedY/8)*horizontalTiles+uint8(positionedX/8))
		tileColor := tile.ReadPixel(positionedY%8, positionedX%8)

		if tileColor == 0 {
			continue
		}

		scanline[pixelX] = p.matchColorPalette(tileColor)
	}

	// Get sprites on the current scanline
	var sprites []*sprite.Sprite
	for i := 0; i < 40; i++ { // Max 40 sprites in OAM at once
		sprite := p.oam.ReadSprite(uint8(i))
		// Assuming sprite height is 8 pixels
		if sprite.Y() <= p.lineCounter && sprite.Y()+8 > p.lineCounter {
			sprites = append(sprites, sprite)
		}
	}

	// Max of 10 sprites per scanline
	if len(sprites) > 10 {
		sprites = sprites[:10]
	}

	// Draw sprites
	for x := uint8(0); x < 160; x++ {
		for _, s := range sprites {
			if (s.X() <= uint8(x)) && (s.X()+8 > uint8(x)) {
				// Draw sprite
				spriteColor := s.ReadPixel(p.lineCounter-s.Y(), uint8(x-s.X()))

				// Check for sprite priority
				// TODO

				// Check for sprite color 0
				if spriteColor == 0 {
					continue
				}

				// Draw
				scanline[x] = p.matchColorPalette(spriteColor)
			}
		}
	}

	return scanline[:]
}

func (p *PPU) matchColorPalette(color uint8) uint8 {
	return p.registers.LCDStatus.PaletteData.Colors[color]
}
