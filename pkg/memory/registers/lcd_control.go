package registers

// 0xFF40 - LCDC - LCD Control (R/W)

type LCDControl struct {
	EnableLCD                     bool // Bit 7 - LCD Display Enable             (0=Off, 1=On)
	WindowUseSecondTileMap        bool // Bit 6 - Window Tile Map Display Select (0=9800-9BFF, 1=9C00-9FFF)
	EnableWindow                  bool // Bit 5 - Window Display Enable          (0=Off, 1=On)
	Use8000Method                 bool // Bit 4 - BG & Window Tile Data Select   (0=8800-97FF, 1=8000-8FFF)
	BackgroundUseSecondaryTileMap bool // Bit 3 - BG Tile Map Display Select     (0=9800-9BFF, 1=9C00-9FFF)
	Sprites8x16                   bool // Bit 2 - OBJ (Sprite) Size              (0=8x8, 1=8x16)
	EnableSprites                 bool // Bit 1 - OBJ (Sprite) Display Enable    (0=Off, 1=On)
	EnableBackgroundAndWindow     bool // Bit 0 - BG & Window Display Priority   (0=Off, 1=On)
}

func (l *LCDControl) Read(addr uint16) uint8 {
	if addr != 0 {
		panic("Invalid address")
	}

	var val uint8
	if l.EnableLCD {
		val |= 1 << 7
	}
	if l.WindowUseSecondTileMap {
		val |= 1 << 6
	}
	if l.EnableWindow {
		val |= 1 << 5
	}
	if l.Use8000Method {
		val |= 1 << 4
	}
	if l.BackgroundUseSecondaryTileMap {
		val |= 1 << 3
	}
	if l.Sprites8x16 {
		val |= 1 << 2
	}
	if l.EnableSprites {
		val |= 1 << 1
	}
	if l.EnableBackgroundAndWindow {
		val |= 1 << 0
	}

	return val
}

func (l *LCDControl) Write(addr uint16, value uint8) {
	if addr != 0 {
		panic("Invalid address")
	}

	l.EnableLCD = value&(1<<7) != 0
	l.WindowUseSecondTileMap = value&(1<<6) != 0
	l.EnableWindow = value&(1<<5) != 0
	l.Use8000Method = value&(1<<4) != 0
	l.BackgroundUseSecondaryTileMap = value&(1<<3) != 0
	l.Sprites8x16 = value&(1<<2) != 0
	l.EnableSprites = value&(1<<1) != 0
	l.EnableBackgroundAndWindow = value&(1<<0) != 0
}
