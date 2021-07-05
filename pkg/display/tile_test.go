package display

import "testing"

var lines = [16]uint8{0x3c, 0x7e, 0x7e, 0x7e, 0x7e, 0x7e, 0x7e, 0x7e, 0x5e, 0x7e, 0x0a, 0x7c, 0x56, 0x38, 0x7c}
var colors = []uint8{0b0, 0b10, 0b11, 0b11, 0b11, 0b11, 0b10, 0b00}

func TestTileColor(t *testing.T) {
	tile := Tile{Bytes: lines}

	for i, color := range colors {
		if tile.ReadPixel(0, uint8(i)) != color {
			t.Errorf("Expected %d, got %d", color, tile.ReadPixel(0, uint8(i)))
		}
	}
}
