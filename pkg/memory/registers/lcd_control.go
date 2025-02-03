package registers

type LCDControl struct {
	BackgroundEnabled bool
	SpritesEnabled    bool
	Sprite8x16        bool
	TileMap1          bool
	WindowEnabled     bool
	WindowTileMap1    bool
	DisplayEnabled    bool
}

func (l *LCDControl) Read(addr uint16) uint8 {
	if addr != 0 {
		panic("Invalid address")
	}

	var val uint8
	if l.BackgroundEnabled {
		val |= 1 << 0
	}
	if l.SpritesEnabled {
		val |= 1 << 1
	}
	if l.Sprite8x16 {
		val |= 1 << 2
	}
	if l.TileMap1 {
		val |= 1 << 3
	}
	if l.WindowEnabled {
		val |= 1 << 5
	}
	if l.WindowTileMap1 {
		val |= 1 << 6
	}
	if l.DisplayEnabled {
		val |= 1 << 7
	}
	return val
}

func (l *LCDControl) Write(addr uint16, value uint8) {
	if addr != 0 {
		panic("Invalid address")
	}

	l.BackgroundEnabled = value&(1<<0) > 0
	l.SpritesEnabled = value&(1<<1) > 0
	l.Sprite8x16 = value&(1<<2) > 0
	l.TileMap1 = value&(1<<3) > 0
	l.WindowEnabled = value&(1<<5) > 0
	l.WindowTileMap1 = value&(1<<6) > 0
	l.DisplayEnabled = value&(1<<7) > 0
}
