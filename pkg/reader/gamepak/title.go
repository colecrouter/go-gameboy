package gamepak

import "strings"

func (gp *GamePak) Title() string {
	// Read the title of the game from the cartridge
	// The title is located at 0x0134-0x0143
	// ASCII encoded, 16 bytes, right padded with 0x00
	title := gp.buffer[0x0134:0x0143]
	return strings.TrimSpace(string(title))
}

func (gp *GamePak) ManufacturerCode() string {
	// Read the manufacturer code from the cartridge
	// The manufacturer code is located at 0x013F-0x0142
	// ASCII encoded, 4 bytes
	code := gp.buffer[0x013F:0x0142]
	return string(code)
}
