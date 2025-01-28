package gamepak

import (
	"log"
	"os"
)

type GamePak struct {
	initialized bool
	buffer      []byte
}

func (g *GamePak) Read(addr uint16) uint8 {
	if !g.initialized {
		panic("GamePak not initialized")
	}
	return g.buffer[addr]
}

func (g *GamePak) Write(addr uint16, data uint8) {
	if !g.initialized {
		panic("GamePak not initialized")
	}
	// TODO: Implement proper write protection
	g.buffer[addr] = data
}

func NewGamePak(b []byte) *GamePak {
	// Load bootrom
	boot, err := os.ReadFile("dmg_boot.bin")
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Copy bootrom to first 0x100 bytes of buffer
	copy(b[:0x100], boot)

	gp := &GamePak{buffer: b, initialized: true}
	return gp
}
