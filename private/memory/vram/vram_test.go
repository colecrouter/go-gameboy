package vram

import (
	"reflect"
	"testing"

	"github.com/colecrouter/gameboy-go/private/memory/vram/tile"
)

// Existing test for reading mapped tiles.
func TestReadMappedTileAt(t *testing.T) {
	v := &VRAM{}

	// Test for Mode8000:
	// Use tile coordinates directly, e.g. (10,0) => map index = 10.
	expectedBytes1 := [16]uint8{}
	for i := 0; i < 16; i++ {
		expectedBytes1[i] = uint8(10 + i)
	}
	expectedTile1 := tile.FromBytes(expectedBytes1)

	v.tiles[10] = expectedTile1
	// Set the tile map entry to index 10.
	v.tileMap0[10] = 10

	// Changed: pass tile coordinates directly instead of pixel values.
	resultTile1 := v.GetMappedTile(0, 10, false, Mode8000)
	if !reflect.DeepEqual(resultTile1, expectedTile1) {
		t.Errorf("Mode8000: expected tile %+v, got %+v", expectedTile1, resultTile1)
	}

	// Test for Mode8800:
	// Use tile coordinates such that map index = 0.
	// To get effective index 130, we set tileMap0[0] = 2 (since int8(2)+128 == 130).
	v.tileMap0[0] = 2
	expectedBytes2 := [16]uint8{}
	for i := 0; i < 16; i++ {
		expectedBytes2[i] = uint8(130 + i)
	}
	expectedTile2 := tile.FromBytes(expectedBytes2)

	v.tiles[130] = expectedTile2

	// Changed: pass tile coordinates directly.
	resultTile2 := v.GetMappedTile(0, 0, false, Mode8800)
	if !reflect.DeepEqual(resultTile2, expectedTile2) {
		t.Errorf("Mode8800: expected tile %+v, got %+v", expectedTile2, resultTile2)
	}
}

// Global test data to use for rendering tests.
var (
	lines  = [16]uint8{0x3c, 0x7e, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x7e, 0x5e, 0x7e, 0x0a, 0x7c, 0x56, 0x38, 0x7c}
	buffer = [8][8]uint8{
		{0, 2, 3, 3, 3, 3, 2, 0},
		{0, 3, 0, 0, 0, 0, 3, 0},
		{0, 3, 0, 0, 0, 0, 3, 0},
		{0, 3, 0, 0, 0, 0, 3, 0},
		{0, 3, 1, 3, 3, 3, 3, 0},
		{0, 1, 1, 1, 3, 1, 3, 0},
		{0, 3, 1, 3, 1, 3, 2, 0},
		{0, 2, 3, 3, 3, 2, 0, 0},
	}
)

// TestTileRendering now checks both addressing modes.
func TestTileRenderingModes(t *testing.T) {
	t.Run("Mode8000", func(t *testing.T) {
		v := &VRAM{}

		// Write the test tile data to VRAM for tile index 0.
		// Assuming a sequential write starting at address 0 writes into tile 0.
		for i, b := range lines {
			v.Write(uint16(i), b)
		}

		// For Mode8000 the mapping is direct, so tileMap0[0] is 0 by default.
		ti := v.GetMappedTile(0, 0, false, Mode8000)

		// Compare each pixel.
		for y := 0; y < 8; y++ {
			for x := 0; x < 8; x++ {
				got := ti.Pixels[y*tile.TILE_SIZE+x]
				expected := buffer[y][x]
				if got != expected {
					t.Errorf("Mode8000 pixel (%d,%d): expected %d, got %d", x, y, expected, got)
				}
			}
		}
	})

	t.Run("Mode8800", func(t *testing.T) {
		v := &VRAM{}
		// For Mode8800, the effective tile index is determined by interpreting the
		// tile map value as a signed int and adding 128.
		// For this test, we want to use a specific tile (say, tile index 128), so we do:
		//   int8(0) + 128 = 128.
		v.tileMap0[0] = 0

		// Write the test tile data into VRAM for tile index 128.
		baseAddr := 128 * 16 // each tile is 16 bytes in tileData.
		for i, b := range lines {
			v.Write(uint16(baseAddr+i), b)
		}
		// Set the tile in the tiles cache, so GetMappedTile can find it.
		expectedTile := tile.FromBytes(lines)
		v.tiles[128] = expectedTile

		ti := v.GetMappedTile(0, 0, false, Mode8800)

		// Compare each pixel.
		for y := 0; y < 8; y++ {
			for x := 0; x < 8; x++ {
				got := ti.Pixels[y*tile.TILE_SIZE+x]
				expected := buffer[y][x]
				if got != expected {
					t.Errorf("Mode8800 pixel (%d,%d): expected %d, got %d", x, y, expected, got)
				}
			}
		}
	})
}
