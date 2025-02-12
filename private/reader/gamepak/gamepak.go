package gamepak

type GamePak struct {
	initialized bool
	buffer      []byte
}

func (g *GamePak) Read(addr uint16) uint8 {
	if !g.initialized {
		panic("GamePak not initialized")
	}
	val := g.buffer[addr]
	return val
}

func (g *GamePak) Write(addr uint16, data uint8) {
	if !g.initialized {
		panic("GamePak not initialized")
	}
	// TODO: Implement proper write protection
	g.buffer[addr] = data
}

func NewGamePak(b []byte) *GamePak {
	gp := &GamePak{buffer: b, initialized: true}
	return gp
}
