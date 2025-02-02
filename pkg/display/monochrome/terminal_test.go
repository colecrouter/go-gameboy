package monochrome

import (
	"testing"

	"github.com/colecrouter/gameboy-go/pkg/memory"
	"github.com/colecrouter/gameboy-go/pkg/memory/registers"
	"github.com/colecrouter/gameboy-go/pkg/memory/vram"
	"github.com/colecrouter/gameboy-go/pkg/processor/ppu"
)

var lines = [16]uint8{0x3c, 0x7e, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x7e, 0x5e, 0x7e, 0x0a, 0x7c, 0x56, 0x38, 0x7c}
var buffer = [16][16]uint8{
	{0, 2, 3, 3, 3, 3, 2, 0},
	{0, 3, 0, 0, 0, 0, 3, 0},
	{0, 3, 0, 0, 0, 0, 3, 0},
	{0, 3, 0, 0, 0, 0, 3, 0},
	{0, 3, 1, 3, 3, 3, 3, 0},
	{0, 1, 1, 1, 3, 1, 3, 0},
	{0, 3, 1, 3, 1, 3, 2, 0},
	{0, 2, 3, 3, 3, 2, 0, 0},
}

func TestTileColor(t *testing.T) {
	// Create new VRAM instance
	vram := &vram.VRAM{}
	for i, line := range lines {
		vram.Write(uint16(i), line)
	}

	// Create new TerminalDisplay instance
	terminal := NewTerminalDisplay()

	// Create new PPU instance
	oam := &memory.OAM{}
	reg := &registers.Registers{}
	ppuInstance := ppu.NewPPU(vram, oam, terminal, reg)

	for i := 0; i < ppu.TotalCyclesPerLine*ppu.TotalLinesPerFrame; i++ {
		ppuInstance.Clock()
	}
	ppuInstance.Clock()

	// Draw
	terminal.Clock()

	// Terminal is not actually 16x16, so need to compensate
	for x, row := range buffer {
		for y, color := range row {
			if terminal.buffer[x][y] != color {
				t.Errorf("Expected %d, got %d", color, terminal.buffer[x][y])
			}
		}
	}
}
