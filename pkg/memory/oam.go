package memory

import "github.com/colecrouter/gameboy-go/pkg/memory/vram/sprite"

type OAM struct {
	buffer [160]byte // 40 sprites, 4 bytes each
}

func (o *OAM) Read(addr uint16) byte {
	return o.buffer[addr]
}

func (o *OAM) Write(addr uint16, data byte) {
	o.buffer[addr] = data
}

func (o *OAM) ReadSprite(index byte) *sprite.Sprite {
	var arr [4]byte
	copy(arr[:], o.buffer[index*4:index*4+4])
	return sprite.NewSprite(arr)
}
