package ppu

import (
	"image"

	"github.com/colecrouter/gameboy-go/pkg/display/monochrome"
	"github.com/colecrouter/gameboy-go/pkg/memory"
	"github.com/colecrouter/gameboy-go/pkg/memory/registers"
	"github.com/colecrouter/gameboy-go/pkg/memory/vram"
	"github.com/colecrouter/gameboy-go/pkg/memory/vram/tile"
	"github.com/colecrouter/gameboy-go/pkg/processor/cpu/lr35902"
)

type PPU struct {
	cpu              *lr35902.LR35902
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
func NewPPU(vram *vram.VRAM, oam *memory.OAM, registers *registers.Registers, cpu *lr35902.LR35902) *PPU {
	return &PPU{
		cpu:       cpu,
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
			p.cpu.ISR(lr35902.VBlankISR)
		}
	} else {
		switch p.lineCycleCounter {
		case 0:
			p.registers.LCDStatus.PPUMode = registers.OAMScan
			p.cpu.ISR(lr35902.LCDSTATISR)
		case oamScanCycles:
			p.registers.LCDStatus.PPUMode = registers.Drawing
			p.cpu.ISR(lr35902.LCDSTATISR)
		case oamScanCycles + pixelTransferCycles:
			p.registers.LCDStatus.PPUMode = registers.HBlank
			p.cpu.ISR(lr35902.LCDSTATISR)
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

	// Draw sprite layer
	spriteHeight := tile.TILE_SIZE
	if p.registers.LCDControl.Sprites8x16 {
		spriteHeight = tile.TILE_SIZE * 2
	}
	for i := 0; i < 40; i++ {
		sprite := p.oam.ReadSprite(i)

		// Skip drawing if the sprite is off the screen
		if sprite.X() >= visibleColumns || sprite.Y() >= visibleLines {
			continue
		}

		for x := 0; x < spriteHeight; x++ {
			for y := 0; y < tile.TILE_SIZE; y++ {
				// Get the pixel from the sprite
				pixel := p.registers.PaletteData.Match(sprite.ReadPixel(uint8(x), uint8(y)))

				matched := p.registers.PaletteData.Match(sprite.ReadPixel(uint8(x), uint8(y)))

				// Skip drawing if the pixel is transparent
				if matched == 0 {
					continue
				}

				// Draw the pixel
				p.image.Set(int(sprite.X())+x, int(sprite.Y())+y, monochrome.Palette[pixel])
			}
		}
	}

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
