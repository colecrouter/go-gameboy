package ppu

import (
	"fmt"
	"testing"
	"time"

	"github.com/colecrouter/gameboy-go/pkg/display/monochrome"
	"github.com/colecrouter/gameboy-go/pkg/memory"
	"github.com/colecrouter/gameboy-go/pkg/memory/registers"
	"github.com/colecrouter/gameboy-go/pkg/memory/vram"
)

// Use a dummy tile that returns a valid color index (e.g. 1) for all pixels.
var dummyTileData = []uint8{
	// Each row: low=0xFF (all bits 1) and high=0x00 (all bits 0) => pixel color = 1 for every pixel.
	0xFF, 0x00, // row 0
	0xFF, 0x00, // row 1
	0xFF, 0x00, // row 2
	0xFF, 0x00, // row 3
	0xFF, 0x00, // row 4
	0xFF, 0x00, // row 5
	0xFF, 0x00, // row 6
	0xFF, 0x00, // row 7
}

func TestPPU(t *testing.T) {
	// Set up VRAM, OAM, registers, and display
	vramModule := &vram.VRAM{}
	oamModule := &memory.OAM{}
	regs := &registers.Registers{}
	// Initialize registers with known values
	regs.LY = 0
	regs.ScrollX = 0
	regs.ScrollY = 0
	regs.PositionX = 0 // window position if used
	regs.PositionY = 0
	// Set tile data mode to unsigned (dummy tile data is loaded accordingly)
	regs.LCDControl.UseSecondaryTileData = true
	regs.LCDControl.BackgroundUseSecondaryTileMap = false
	regs.LCDControl.WindowUseSecondTileMap = false
	// Initialize a simple grayscale palette mapping (index 0 -> black, index 1 -> dark gray, etc.)
	regs.PaletteData.Colors = [4]uint8{0, 85, 170, 255}

	display := monochrome.NewTerminalDisplay()
	ppuUnit := NewPPU(vramModule, oamModule, display, regs)

	// Load dummyTileData into the first tile (16 bytes per tile)
	for i, b := range dummyTileData {
		vramModule.Write(uint16(i), b)
	}

	// Map the first 32 entries of tileMap0 to point to the first tile (tile index 0)
	for i := 0; i < 32; i++ {
		vramModule.Write(0x1800+uint16(i), 0)
	}

	// Clock the PPU enough times to build a frame
	for i := 0; i < 10000; i++ {
		ppuUnit.Clock()
	}

	// Print out one scanline explicitly for debug purposes.
	scanline := ppuUnit.getScanline()
	fmt.Printf("Debug (test): First 8 pixels of scanline: %v\n", scanline[:8])

	// Small pause to allow the display to update
	time.Sleep(500 * time.Millisecond)
	fmt.Println("Test complete. Check debug logs and display output.")

	display.Clock()
}
