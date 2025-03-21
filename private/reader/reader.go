package reader

import (
	bootroms "github.com/colecrouter/gameboy-go/private/memory/roms"
	"github.com/colecrouter/gameboy-go/private/reader/gamepak"
)

type CartridgeReader struct {
	disableBootRom *bool
	cartridge      *gamepak.GamePak
}

func NewCartridgeReader(disableBootRom *bool) *CartridgeReader {
	return &CartridgeReader{
		disableBootRom: disableBootRom,
	}
}

func (cr *CartridgeReader) InsertCartridge(game *gamepak.GamePak) {
	cr.cartridge = game
}

func (cr *CartridgeReader) Cartridge() *gamepak.GamePak {
	return cr.cartridge
}

func (cr *CartridgeReader) Read(addr uint16) uint8 {
	if cr.disableBootRom == nil {
		panic("Boot ROM disable flag not set")
	}

	if (!*cr.disableBootRom) && addr < 0x100 {
		return bootroms.DMG_BOOT[addr]
	}
	return cr.cartridge.Read(addr)
}

func (cr *CartridgeReader) Write(addr uint16, val uint8) {
	// TODO implement cartridge write protection
	// cr.cartridge.Write(addr, val)
	// panic("Cartridge write protection not implemented")
}
