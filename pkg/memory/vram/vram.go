package vram

import "github.com/colecrouter/gameboy-go/pkg/memory/vram/tile"

type VRAM struct {
	tileData [0x1800]uint8 // 6KB
	tileMap0 [0x400]uint8  // 1KB
	tileMap1 [0x400]uint8  // 1KB
}

func (v *VRAM) Read(addr uint16) uint8 {
	if addr < 0x1800 {
		return v.tileData[addr]
	} else if addr < 0x1C00 {
		return v.tileMap0[addr-0x1800]
	} else {
		return v.tileMap1[addr-0x1C00]
	}
}

func (v *VRAM) Write(addr uint16, data uint8) {
	if addr < 0x1800 {
		v.tileData[addr] = data
	} else if addr < 0x1C00 {
		v.tileMap0[addr-0x1800] = data
	} else {
		v.tileMap1[addr-0x1C00] = data
	}
}

// readTile reads a tile from VRAM. The memory bank contains 384 tiles ()
func (v *VRAM) readTile(index uint8) tile.Tile {
	tile := tile.Tile{}
	for i := 0; i < 16; i++ {
		tile.Bytes[i] = v.Read(uint16(index)*16 + uint16(i))
	}

	return tile
}

// func (v *VRAM) ReadTileMap(index uint8) [32][32]tile.Tile {
// 	currentMap := v.tileMap0
// 	if index == 1 {
// 		currentMap = v.tileMap1
// 	} else if index != 0 {
// 		panic("Invalid tile map index")
// 	}

// 	tileMap := [32][32]tile.Tile{}
// 	for i := 0; i < 32; i++ {
// 		for j := 0; j < 32; j++ {
// 			tileMap[i][j] = v.ReadTile(currentMap[i*32+j])
// 		}
// 	}

// 	return tileMap
// }

func (v *VRAM) ReadMappedTile(tilemap uint8, index uint8) tile.Tile {
	currentMap := v.tileMap0
	if tilemap == 1 {
		currentMap = v.tileMap1
	} else if tilemap != 0 {
		panic("Invalid tile map index")
	}

	return v.readTile(currentMap[index])
}
