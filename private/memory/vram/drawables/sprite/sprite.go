package sprite

import (
	"github.com/colecrouter/gameboy-go/private/memory/vram"
)

type Sprite struct {
	y    uint8 // 0x0
	x    uint8 // 0x1
	tile uint8 // 0x2

	// 0x3
	Priority   bool  // Bit 7
	FlipY      bool  // Bit 6
	FlipX      bool  // Bit 5
	UseBank1   bool  // Bit 3
	DMGPalette uint8 // Bit 4
	CGBPalette uint8 // Bits 0-2

	vram       *vram.VRAM
	enable8x16 *bool
}

func (s *Sprite) Read(addr uint16) uint8 {
	switch addr {
	case 0:
		return s.y
	case 1:
		return s.x
	case 2:
		return s.tile
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
		if s.UseBank1 {
			b |= 1 << 3
		}
		b |= (s.DMGPalette & 0b0000_0001) << 4
		b |= s.CGBPalette & 0b0000_0111
		return b
	default:
		panic("Invalid address")
	}
}

func (s *Sprite) Write(addr uint16, data uint8) {
	switch addr {
	case 0:
		s.y = data
	case 1:
		s.x = data
	case 2:
		s.tile = data
	case 3:
		s.Priority = data&(1<<7) != 0
		s.FlipY = data&(1<<6) != 0
		s.FlipX = data&(1<<5) != 0
		s.DMGPalette = data & (1 << 4) >> 4
		s.UseBank1 = data&(1<<3) != 0
		s.CGBPalette = data & 0b0000_0111
	default:
		panic("Invalid address")
	}
}

func NewSprite(vram *vram.VRAM, data [4]byte, enable8x16 *bool) *Sprite {
	return &Sprite{
		y:          data[0],
		x:          data[1],
		tile:       data[2],
		Priority:   data[3]&(1<<7) != 0,
		FlipY:      data[3]&(1<<6) != 0,
		FlipX:      data[3]&(1<<5) != 0,
		DMGPalette: data[3] & (1 << 4) >> 4,
		UseBank1:   data[3]&(1<<3) != 0,
		CGBPalette: data[3] & 0b0000_0111,
		vram:       vram,
		enable8x16: enable8x16,
	}
}

func (s *Sprite) Pixels() []uint8 {
	if *s.enable8x16 {
		return append(s.vram.GetTile(int(s.tile)&0xFE).Pixels(), s.vram.GetTile(int(s.tile)|0x01).Pixels()...)
	}
	return s.vram.GetTile(int(s.tile)).Pixels()
}

// Y returns the Y coordinate of the sprite on the screen.
func (s *Sprite) Y() uint8 {
	return s.y - 16
}

// X returns the X coordinate of the sprite on the screen.
func (s *Sprite) X() uint8 {
	return s.x - 8
}
