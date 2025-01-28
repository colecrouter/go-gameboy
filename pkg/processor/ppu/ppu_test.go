package ppu

import (
	"testing"
)

func TestDrawTestPattern(t *testing.T) {
	// vram := &vram.VRAM{}
	// oam := &memory.OAM{}
	// display := monochrome.NewTerminalDisplay()
	// ppu := NewPPU(vram, oam, display)

	// for tileIndex := uint16(0); tileIndex < 32*32; tileIndex++ {
	// 	for pixelIndex := uint; pixelIndex < 8*8; pixelIndex++ {

	// 	}
	// 	tile := vram.ReadTile(uint16(tileIndex))
	// 	for y := 0; y < 8; y++ {
	// 		for x := 0; x < 8; x++ {
	// 			tile.WritePixel(uint8(y), uint8(x), uint8((tileIndex+int(y)+int(x))%4))
	// 		}
	// 	}
	// }

	// // Simulate enough clock cycles to draw the test pattern
	// for i := 0; i < TotalCyclesPerLine*TotalLinesPerFrame; i++ {
	// 	ppu.Clock()
	// }

	// // TODO: Add assertions to verify the test pattern
}
