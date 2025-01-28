package vram

import "github.com/colecrouter/gameboy-go/pkg/memory/vram/tile"

type VRAM struct {
	data [0x2000]uint8 // 8KB
}

func (v *VRAM) Read(addr uint16) uint8 {
	return v.data[addr]
}

func (v *VRAM) Write(addr uint16, data uint8) {
	v.data[addr] = data
}

// ReadTile reads a tile from VRAM. The memory bank contains 384 tiles ()
func (v *VRAM) ReadTile(index uint8) tile.Tile {
	tile := tile.Tile{}
	for i := 0; i < 16; i++ {
		tile.Bytes[i] = v.Read(uint16(index)*16 + uint16(i))
	}

	return tile
}
