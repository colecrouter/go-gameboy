package vram

import (
	"reflect"
	"testing"

	"github.com/colecrouter/gameboy-go/pkg/memory"
	"github.com/colecrouter/gameboy-go/pkg/memory/vram/tile"
)

func TestReadMappedTileAt(t *testing.T) {
	v := &VRAM{}

	// Test for Mode8000:
	// Use pixel coordinates such that (x/8, y/8) equals (10,0) => map index = 10.
	expectedBytes1 := [16]uint8{}
	for i := 0; i < 16; i++ {
		expectedBytes1[i] = uint8(10 + i)
	}
	expectedTile1 := tile.FromBytes(expectedBytes1)

	v.tiles[10] = expectedTile1
	// Set the tile map entry to index 10.
	v.tileMap0[10] = 10

	resultTile1 := v.GetMappedTile(10*8, 0, false, Mode8000)
	if !reflect.DeepEqual(resultTile1, expectedTile1) {
		t.Errorf("Mode8000: expected tile %+v, got %+v", expectedTile1, resultTile1)
	}

	// Test for Mode8800:
	// Use pixel coordinates such that (x/8, y/8) equals (0,0) => map index = 0.
	// To get effective index 130, we set tileMap0[0] = 2, since int8(2)+128 == 130.
	v.tileMap0[0] = 2
	expectedBytes2 := [16]uint8{}
	for i := 0; i < 16; i++ {
		expectedBytes2[i] = uint8(130 + i)
	}
	expectedTile2 := tile.FromBytes(expectedBytes2)

	v.tiles[130] = expectedTile2

	resultTile2 := v.GetMappedTile(0, 0, false, Mode8800)
	if !reflect.DeepEqual(resultTile2, expectedTile2) {
		t.Errorf("Mode8800: expected tile %+v, got %+v", expectedTile2, resultTile2)
	}
}

var lines = [16]uint8{0x3c, 0x7e, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x7e, 0x5e, 0x7e, 0x0a, 0x7c, 0x56, 0x38, 0x7c}
var buffer = [8][8]uint8{
	{0, 2, 3, 3, 3, 3, 2, 0},
	{0, 3, 0, 0, 0, 0, 3, 0},
	{0, 3, 0, 0, 0, 0, 3, 0},
	{0, 3, 0, 0, 0, 0, 3, 0},
	{0, 3, 1, 3, 3, 3, 3, 0},
	{0, 1, 1, 1, 3, 1, 3, 0},
	{0, 3, 1, 3, 1, 3, 2, 0},
	{0, 2, 3, 3, 3, 2, 0, 0},
}

func TestTileRendering(t *testing.T) {
	v := &VRAM{}

	// Load the tile data into VRAM at tile 0
	for i, b := range lines {
		v.Write(uint16(i), b)
	}

	// Read the tile data back out and compare
	// By default, everything should be 0, so getting a blankly mapped tile should return the tile at index 0.
	ti := v.GetMappedTile(0, 0, false, Mode8000)

	// Check that the tile data matches the expected data
	for y, row := range buffer {
		for x, color := range row {
			if ti.Pixels[y*tile.TILE_SIZE+x] != color {
				t.Errorf("Expected %d, got %d", color, ti.Pixels[y*tile.TILE_SIZE+x])
			}
		}
	}

	if t.Failed() {
		/// Print tile data
		memory.PrintMemory(v, 0x0, 0x1F)

		// Print tile map data
		memory.PrintMemory(v, 0x1800, 0x181F)
	}
}
