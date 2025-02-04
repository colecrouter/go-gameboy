package ppu

import (
	"fmt" // added for debugging

	"github.com/colecrouter/gameboy-go/pkg/display"
	"github.com/colecrouter/gameboy-go/pkg/memory"
	"github.com/colecrouter/gameboy-go/pkg/memory/registers"
	"github.com/colecrouter/gameboy-go/pkg/memory/vram"
	"github.com/colecrouter/gameboy-go/pkg/memory/vram/sprite"
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
		// p.HBlank()
	case registers.VBlank:
		// p.VBlank()
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

	useSigned := !p.registers.LCDControl.UseSecondaryTileData

	// Draw background
	for pixelX := uint8(0); pixelX < visibleColumns; pixelX++ {
		useSecondaryMap := p.registers.LCDControl.BackgroundUseSecondaryTileMap
		scrolledY := p.registers.LY
		scrolledX := p.registers.ScrollX + pixelX

		tileX := scrolledX / 8
		tileY := scrolledY / 8

		// Multiply by 32 since the DMG tile map is 32 tiles wide, not 20.
		index := uint16(tileY)*32 + uint16(tileX)
		tile := p.vram.ReadMappedTileAt(index, useSecondaryMap, useSigned)
		tileColor := tile.Pixels[scrolledY%8][scrolledX%8]
		mapped := p.matchColorPalette(tileColor)
		scanline[pixelX] = mapped

		// Debug only for first pixel on first scanline.
		if p.registers.LY == 0 && pixelX == 0 {
			fmt.Printf("Debug: Background tile index %d, raw color %d, mapped %d\n", index, tileColor, mapped)
		}
	}

	// Draw window
	for pixelX := uint8(0); pixelX < visibleColumns; pixelX++ {
		useSecondaryMap := p.registers.LCDControl.WindowUseSecondTileMap
		positionedY := p.registers.PositionY + p.registers.LY
		positionedX := p.registers.PositionX + pixelX

		if positionedX >= visibleColumns {
			continue
		}

		tileX := positionedX / 8
		tileY := positionedY / 8

		// Use a row width of 32 as above.
		index := uint16(tileY)*32 + uint16(tileX)
		tile := p.vram.ReadMappedTileAt(index, useSecondaryMap, useSigned)
		tileColor := tile.Pixels[positionedY%8][positionedX%8]
		mapped := p.matchColorPalette(tileColor)
		scanline[pixelX] = mapped

		// Debug for first pixel (optional)
		if p.registers.LY == 0 && pixelX == 0 {
			fmt.Printf("Debug: Window tile index %d, raw color %d, mapped %d\n", index, tileColor, mapped)
		}
	}

	// Get sprites on the current scanline
	var sprites []*sprite.Sprite
	for i := 0; i < 40; i++ { // Max 40 sprites in OAM at once
		// sprite := p.oam.ReadSprite(uint8(i))
		// // Assuming sprite height is 8 pixels
		// if sprite.Y() <= p.registers.LY && sprite.Y()+8 > p.registers.LY {
		// 	sprites = append(sprites, sprite)
		// }
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
				spriteColor := s.ReadPixel(p.registers.LY-s.Y(), uint8(x-s.X()))

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

	// Debug print for first scanline (for example)
	if p.registers.LY == 0 {
		fmt.Printf("Debug: LY %d first 8 scanline pixels: %v\n", p.registers.LY, scanline[:8])
	}

	return scanline[:]
}

func (p *PPU) matchColorPalette(color uint8) uint8 {
	// fmt.Printf("Color: %v\n", p.registers.PaletteData)
	return p.registers.PaletteData.Colors[color]
	// return color
}
