package sprite

import (
	"github.com/colecrouter/gameboy-go/private/memory/vram"
	"github.com/colecrouter/gameboy-go/private/memory/vram/tile"
)

type Sprite struct {
	Y          uint8
	X          uint8
	Tile       uint8
	Priority   bool
	FlipY      bool
	FlipX      bool
	DmgPalette bool
	Bank       bool
	CGBPalette uint8

	vram *vram.VRAM
}

func (s *Sprite) Read(addr uint16) uint8 {
	switch addr {
	case 0:
		return s.Y
	case 1:
		return s.X
	case 2:
		return s.Tile
	case 3:
		var b byte
		if s.Priority {
			b |= 1 << 7
		}
		if s.FlipY {
			b |= 1 << 6
		}
		if s.FlipX {
			b |= 1 << 5
		}
		if s.DmgPalette {
			b |= 1 << 4
		}
		if s.Bank {
			b |= 1 << 3
		}
		b |= s.CGBPalette & 0b0000_0111
		return b
	default:
		panic("Invalid address")
	}
}

func (s *Sprite) Write(addr uint16, data uint8) {
	switch addr {
	case 0:
		s.Y = data
	case 1:
		s.X = data
	case 2:
		s.Tile = data
	case 3:
		s.Priority = data&(1<<7) != 0
		s.FlipY = data&(1<<6) != 0
		s.FlipX = data&(1<<5) != 0
		s.DmgPalette = data&(1<<4) != 0
		s.Bank = data&(1<<3) != 0
		s.CGBPalette = data & 0b0000_0111
	default:
		panic("Invalid address")
	}
}

func (s *Sprite) GetTile() *tile.Tile {
	return s.vram.GetTile(int(s.Tile))
}

func NewSprite(vram *vram.VRAM, data [4]byte) *Sprite {
	return &Sprite{
		Y:          data[0],
		X:          data[1],
		Tile:       data[2],
		Priority:   data[3]&(1<<7) != 0,
		FlipY:      data[3]&(1<<6) != 0,
		FlipX:      data[3]&(1<<5) != 0,
		DmgPalette: data[3]&(1<<4) != 0,
		Bank:       data[3]&(1<<3) != 0,
		CGBPalette: data[3] & 0b0000_0111,
		vram:       vram,
	}
}
