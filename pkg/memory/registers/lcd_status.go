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
	} // 0x41
	ScrollY   uint8 // 0x42
	ScrollX   uint8 // 0x43
	LY        uint8 // 0x44
	LYCompare uint8 // 0x45
}

// Read returns the value of the LCD status register
func (l *LCDStatus) Read(addr uint16) uint8 {
	switch addr {
	case 0x00:
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
	case 0x01:
		return l.ScrollY
	case 0x02:
		return l.ScrollX
	case 0x03:
		return l.LY
	case 0x04:
		return l.LYCompare
	}

	panic("Invalid address")
}

// Write sets the value of the LCD status register
func (l *LCDStatus) Write(addr uint16, value uint8) {
	switch addr {
	case 0x00:
		l.Status.LYCInterrupt = value&(1<<6) > 0
		l.Status.Mode2Interrupt = value&(1<<5) > 0
		l.Status.Mode1Interrupt = value&(1<<4) > 0
		l.Status.Mode0Interrupt = value&(1<<3) > 0
		l.Status.LYCMatch = value&(1<<2) > 0
		l.Status.PPUMode = PPUState(value & 0x3)
	case 0x01:
		l.ScrollY = value
	case 0x02:
		l.ScrollX = value
	case 0x03:
		l.LY = value
	case 0x04:
		l.LYCompare = value
	default:
		panic("Invalid address")
	}
}
