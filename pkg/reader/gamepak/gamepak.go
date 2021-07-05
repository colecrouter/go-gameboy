package gamepak

type GamePak struct {
	buffer []byte
}

func (g *GamePak) Read(addr uint16) uint8 {
	return g.buffer[addr]
}

func (g *GamePak) Write(addr uint16, data uint8) {
	// TODO: Implement proper write protection
	g.buffer[addr] = data
}

func NewGamePak(b []byte) *GamePak {
	return &GamePak{buffer: b}
}
