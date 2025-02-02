package registers

// PPU modes
type PPUState uint8

const (
	HBlank PPUState = iota
	VBlank
	OAMScan
	Drawing
)

// LCDStatus represents the LCD status register
type LCDStatus struct {
	Status struct {
		LYCInterrupt   bool
		Mode2Interrupt bool
		Mode1Interrupt bool
		Mode0Interrupt bool
		LYCMatch       bool
		PPUMode        PPUState
	} // 0xFF41
	ScrollY     uint8 // 0xFF42
	ScrollX     uint8 // 0xFF43
	YCoordinate uint8 // 0xFF44
	LYCompare   uint8 // 0xFF45
	// DMA 0xFF46
	// TODO http://www.codeslinger.co.uk/pages/projects/gameboy/dma.html
	PaletteData Palette // 0xFF47
	PositionY   uint8   // 0xFF4A
	PositionX   uint8   // 0xFF4B
}

// Read returns the value of the LCD status register
func (l *LCDStatus) Read(addr uint16) uint8 {
	switch addr {
	case 0xFF41:
		val := uint8(0)
		if l.Status.LYCInterrupt {
			val |= 1 << 6
		}
		if l.Status.Mode2Interrupt {
			val |= 1 << 5
		}
		if l.Status.Mode1Interrupt {
			val |= 1 << 4
		}
		if l.Status.Mode0Interrupt {
			val |= 1 << 3
		}
		if l.Status.LYCMatch {
			val |= 1 << 2
		}
		val |= uint8(l.Status.PPUMode)
		return val
	case 0xFF42:
		return l.ScrollY
	case 0xFF43:
		return l.ScrollX
	case 0xFF44:
		return l.YCoordinate
	case 0xFF45:
		return l.LYCompare
	case 0xFF4A:
		return l.PositionY
	case 0xFF4B:
		return l.PositionX
	}

	return 0
}

// Write sets the value of the LCD status register
func (l *LCDStatus) Write(addr uint16, value uint8) {
	switch addr {
	case 0xFF41:
		l.Status.LYCInterrupt = value&(1<<6) > 0
		l.Status.Mode2Interrupt = value&(1<<5) > 0
		l.Status.Mode1Interrupt = value&(1<<4) > 0
		l.Status.Mode0Interrupt = value&(1<<3) > 0
		l.Status.LYCMatch = value&(1<<2) > 0
		l.Status.PPUMode = PPUState(value & 0x3)
	case 0xFF42:
		l.ScrollY = value
	case 0xFF43:
		l.ScrollX = value
	case 0xFF44:
		l.YCoordinate = value
	case 0xFF45:
		l.LYCompare = value
	case 0xFF4A:
		l.PositionY = value
	case 0xFF4B:
		l.PositionX = value
	}
}
