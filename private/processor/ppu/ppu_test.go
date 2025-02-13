package ppu

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/colecrouter/gameboy-go/private/memory"
	"github.com/colecrouter/gameboy-go/private/memory/registers"
	"github.com/colecrouter/gameboy-go/private/memory/vram"
	"github.com/colecrouter/gameboy-go/private/processor/cpu/lr35902"
)

// Use a dummy tile that returns a valid color index (e.g. 1) for all pixels.
var dummyTileData = [16]uint8{0x3c, 0x7e, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x7e, 0x5e, 0x7e, 0x0a, 0x7c, 0x56, 0x38, 0x7c}

func TestPPU(t *testing.T) {
	// Set up VRAM, OAM, registers, and display
	vramModule := &vram.VRAM{}
	oamModule := &memory.OAM{}
	regs := &registers.Registers{}
	// Set the palette to a simple 4-color palette
	regs.PaletteData.Set([4]uint8{0, 1, 2, 3})
	cpu := lr35902.NewLR35902(&memory.Bus{}, regs)

	ppuUnit := NewPPU(vramModule, oamModule, regs, cpu.ISR)

	ppuUnit.registers.LCDControl.Use8000Method = true

	// Load dummyTileData into the first tile (16 bytes per tile)
	for i, b := range dummyTileData {
		vramModule.Write(uint16(i), b)
	}

	// Map the first 32 entries of tileMap0 to point to the first tile (tile index 0)
	for i := 0; i < 32; i++ {
		vramModule.Write(0x1800+uint16(i), 0)
	}

	ppuUnit.DisplayClock()

	// Uncomment to print the rendered image to the console
	// fmt.Println(renderer.RenderANSI(ppuUnit.image))

	scanline := ppuUnit.image.Pix[:16]

	expected := []uint8{0, 2, 3, 3, 3, 3, 2, 0, 0, 2, 3, 3, 3, 3, 2, 0}

	if !reflect.DeepEqual(scanline, expected) {
		t.Errorf("Expected scanline\n\t%v\n—got\n\t%v", expected, scanline)
	}

	// Small pause to allow the display to update
	time.Sleep(500 * time.Millisecond)
	fmt.Println("Test complete. Check debug logs and display output.")
}
