package vram

import (
	"github.com/colecrouter/gameboy-go/pkg/memory/vram/tile"
)

type VRAM struct {
	tileData [0x1800]uint8 // 0x8000-0x9FFF
	tileMap0 [0x400]uint8  // 0x9800-0x9BFF
	tileMap1 [0x400]uint8  // 0x9C00-0x9FFF

	tiles [384]*tile.Tile
}

func (v *VRAM) Read(addr uint16) uint8 {
	if addr < 0x1800 {
		return v.tileData[addr]
	} else if addr < 0x1C00 {
		return v.tileMap0[addr-0x1800]
	} else if addr < 0x2000 {
		return v.tileMap1[addr-0x1C00]
	} else {
		panic("Invalid address")
	}
}

func (v *VRAM) Write(addr uint16, data uint8) {
	if addr < 0x1800 {
		v.tileData[addr] = data

		// Figure out which tile this is and update it
		index := addr / 16

		var tileBytes [16]uint8
		copy(tileBytes[:], v.tileData[index*16:index*16+16])
		v.tiles[index] = tile.FromBytes(tileBytes)
	} else if addr < 0x1C00 {
		v.tileMap0[addr-0x1800] = data
	} else if addr < 0x2000 {
		v.tileMap1[addr-0x1C00] = data
	} else {
		panic("Invalid address")
	}
}

/*
	Addressing modes:

	- $8000 addressing mode:
		- Tile data is stored in the range 0x8000-0x97FF
		- Base is 0x8000, range is 0–255
		- 0-127 are block 1, 128-255 are block 2
	- $8800 addressing mode:
		- Tile data is stored in the range 0x8800-0x97FF
		- Base is 0x9000, range is -128-127
		- -128–(-1) are block 1, 0–127 are block 2

	Switching between addressing modes is done by setting the LCD Control register bit 4.
	Except for sprites, which always use $8000 addressing mode.
*/

type TileAddressingMode bool

const (
	Mode8000 TileAddressingMode = false
	Mode8800 TileAddressingMode = true
)

type TileMapMode bool

const (
	Map0 TileMapMode = false
	Map1 TileMapMode = true
)

// ReadMappedTileAt reads a tile from the VRAM at the given pixel coordinates.
func (v *VRAM) ReadMappedTileAt(x, y uint8, mapMode TileMapMode, addressingMode TileAddressingMode) *tile.Tile {
	tileX := x / 8
	tileY := y / 8
	mapIndex := uint16(tileY)*32 + uint16(tileX) // assuming a 32-tile wide tilemap; adjust if needed

	currentMap := &v.tileMap0
	if mapMode {
		currentMap = &v.tileMap1
	}

	tileIndex := currentMap[mapIndex]

	var effectiveIndex int
	switch addressingMode {
	case Mode8000:
		effectiveIndex = int(tileIndex)
	case Mode8800:
		effectiveIndex = int(int8(tileIndex)) + 128
	default:
		panic("Invalid addressing mode")
	}

	if effectiveIndex == 270 {
		println("woo")
	}

	if effectiveIndex != 0 {
		// fmt.Println("Should be drawing tile", effectiveIndex)
	}

	return v.tiles[effectiveIndex]
}
