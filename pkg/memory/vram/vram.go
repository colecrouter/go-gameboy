package vram

import (
	"github.com/colecrouter/gameboy-go/pkg/memory/vram/tile"
)

type VRAM struct {
	tileData [0x1800]uint8 // 0x8000-0x97FF
	tileMap0 [0x400]uint8  // 0x9800-0x9BFF
	tileMap1 [0x400]uint8  // 0x9C00-0x9FFF

	tiles [384]tile.Tile
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
		v.tiles[index] = tile.Tile{}
		for i := 0; i < 16; i++ {
			v.tiles[index].Bytes[i] = v.Read(index*16 + uint16(i))
		}
	} else if addr < 0x1C00 {
		v.tileMap0[addr-0x1800] = data
	} else if addr < 0x2000 {
		v.tileMap1[addr-0x1C00] = data
	} else {
		panic("Invalid address")
	}
}

// ReadTile reads a tile from VRAM. The memory bank contains 384 tiles ()
func (v *VRAM) ReadTile(index uint8) *tile.Tile {
	return &v.tiles[index]
}

func (v *VRAM) ReadMappedTileAt(index uint16, useSecondaryMap, useSigned bool) *tile.Tile {
	currentMap := &v.tileMap0
	if useSecondaryMap {
		currentMap = &v.tileMap1
	}

	// Compute the index from tile coordinates.
	tileIndex := currentMap[index]
	// println("tile:", index, "index:", tileIndex)

	if useSigned {
		tileIndex = uint8(int16(int8(tileIndex)) + 128)
	}

	return v.ReadTile(tileIndex)
}
