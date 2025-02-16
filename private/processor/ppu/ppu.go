package ppu

import (
	"image"

	"github.com/colecrouter/gameboy-go/private/display/monochrome"
	"github.com/colecrouter/gameboy-go/private/memory"
	"github.com/colecrouter/gameboy-go/private/memory/registers"
	"github.com/colecrouter/gameboy-go/private/memory/vram"
	"github.com/colecrouter/gameboy-go/private/memory/vram/tile"
)

type PPU struct {
	interrupt        *registers.Interrupt
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
func NewPPU(vram *vram.VRAM, oam *memory.OAM, registers *registers.Registers, ie *registers.Interrupt) *PPU {
	return &PPU{
		interrupt: ie,
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

		if p.registers.LY == visibleLines {
			p.interrupt.VBlank = true
		}
	} else {
		switch p.lineCycleCounter {
		case 0:
			p.registers.LCDStatus.PPUMode = registers.OAMScan
			p.interrupt.LCD = true
		case oamScanCycles:
			p.registers.LCDStatus.PPUMode = registers.Drawing
			p.interrupt.LCD = true
		case oamScanCycles + pixelTransferCycles:
			p.registers.LCDStatus.PPUMode = registers.HBlank
			p.interrupt.LCD = true
		}
	}

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

// Clock updates the image produced by the PPU
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
	bgMapMode := vram.TileMapMode(p.registers.LCDControl.BackgroundUseSecondTileMap)
	for tileY := uint8(0); tileY < visibleLines/tile.TILE_SIZE; tileY++ {
		pixelY := (int(tileY)*tile.TILE_SIZE - int(p.registers.ScrollY) + 256) % 256

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

	// Draw window layer
	if p.registers.LCDControl.EnableWindow {
		fgMapMode := vram.TileMapMode(p.registers.LCDControl.WindowUseSecondTileMap)
		for tileY := uint8(0); tileY < visibleLines/tile.TILE_SIZE; tileY++ {
			pixelY := int(tileY)*tile.TILE_SIZE - int(p.registers.WindowY)

			if pixelY < 0 {
				continue
			}

			for tileX := uint8(0); tileX < visibleColumns/tile.TILE_SIZE; tileX++ {
				pixelX := int(tileX)*tile.TILE_SIZE - int(p.registers.WindowX) + 7

				if pixelX < 0 {
					continue
				}

				t := p.vram.GetMappedTile(tileY, tileX, fgMapMode, addressingMode)

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
	}

	// // Draw sprite layer
	// spriteHeight := tile.TILE_SIZE
	// if p.registers.LCDControl.Sprites8x16 {
	// 	spriteHeight = tile.TILE_SIZE * 2
	// }
	// for i := 0; i < 40; i++ {
	// 	sprite := p.oam.ReadSprite(i)

	// 	// Skip drawing if the sprite is off the screen
	// 	if sprite.X() >= visibleColumns || sprite.Y() >= visibleLines {
	// 		continue
	// 	}

	// 	for x := 0; x < spriteHeight; x++ {
	// 		for y := 0; y < tile.TILE_SIZE; y++ {
	// 			// Get the pixel from the sprite

	// 			pixel := sprite.ReadPixel(uint8(x), uint8(y))
	// 			matched := p.registers.PaletteData.Match(pixel)

	// 			// Skip drawing if the pixel is transparent
	// 			if matched == 0 {
	// 				continue
	// 			}

	// 			// Draw the pixel
	// 			p.image.Set(int(sprite.X())+x, int(sprite.Y())+y, monochrome.Palette[matched])
	// 		}
	// 	}
	// }

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
