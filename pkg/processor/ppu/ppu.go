package ppu

import (
	"github.com/colecrouter/gameboy-go/pkg/display"
	"github.com/colecrouter/gameboy-go/pkg/memory"
	"github.com/colecrouter/gameboy-go/pkg/memory/vram"
	"github.com/colecrouter/gameboy-go/pkg/memory/vram/sprite"
)

// PPU modes
type PPUState uint8

const (
	HBlank PPUState = iota
	VBlank
	OAMScan
	Drawing
)

type PPU struct {
	vram             *vram.VRAM
	oam              *memory.OAM
	display          display.Display
	Mode             PPUState
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
)

// NewPPU creates a new PPU instance
func NewPPU(vram *vram.VRAM, oam *memory.OAM, display display.Display) *PPU {
	return &PPU{
		vram:    vram,
		oam:     oam,
		display: display,
		Mode:    HBlank,
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
		p.Mode = OAMScan
	case oamScanCycles:
		p.Mode = Drawing
	case oamScanCycles + pixelTransferCycles:
		p.Mode = HBlank
	}
	if p.lineCounter >= visibleLines {
		p.Mode = VBlank
	}

	// Handle mode-specific operations
	switch p.Mode {
	case OAMScan:
		// p.OAMScan()
	case Drawing:
		// p.DrawScanline()
		line := p.getScanline()
		p.display.DrawScanline(p.lineCounter, line)
	case HBlank:
		// p.HBlank()
	case VBlank:
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
	// Draw the scanline
	var scanline [160]byte

	// Background from VRAM
	tileY := p.lineCounter / 8
	for tileX := uint8(0); tileX < 20; tileX++ {
		for pixelX := uint8(0); pixelX < 8; pixelX++ {
			tileIndex := tileY*32 + tileX
			tile := p.vram.ReadTile(tileIndex)
			tileColor := tile.ReadPixel(p.lineCounter%8, pixelX)
			scanline[tileX*8+pixelX] = tileColor
		}
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
				// Draw sprit
				s.ReadPixel(p.lineCounter-s.Y(), uint8(x-s.X()))
			}
		}
	}

	return scanline[:]
}
