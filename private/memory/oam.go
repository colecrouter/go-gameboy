package memory

import (
	"github.com/colecrouter/gameboy-go/private/memory/vram"
	"github.com/colecrouter/gameboy-go/private/memory/vram/sprite"
)

type OAM struct {
	Sprites [40]sprite.Sprite

	buffer      [160]byte // 40 sprites, 4 bytes each
	vram        *vram.VRAM
	initialized bool
}

func (o *OAM) Read(addr uint16) byte {
	if !o.initialized {
		panic("OAM not initialized")
	}

	return o.buffer[addr]
}

func (o *OAM) Write(addr uint16, data byte) {
	if !o.initialized {
		panic("OAM not initialized")
	}

	o.buffer[addr] = data

	// Update sprite
	index := addr / 4
	o.Sprites[index] = *sprite.NewSprite(o.vram, [4]byte{
		o.buffer[index*4],
		o.buffer[index*4+1],
		o.buffer[index*4+2],
		o.buffer[index*4+3],
	})
}

func (o *OAM) ReadSprite(index int) *sprite.Sprite {
	var arr [4]byte
	copy(arr[:], o.buffer[index*4:index*4+4])
	return sprite.NewSprite(o.vram, arr)
}

func NewOAM(vram *vram.VRAM) *OAM {
	return &OAM{
		vram:        vram,
		initialized: true,
	}
}
