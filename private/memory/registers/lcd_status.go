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
	LYCInterrupt   bool
	Mode2Interrupt bool
	Mode1Interrupt bool
	Mode0Interrupt bool
	LYCMatch       bool
	PPUMode        PPUState
} // 0x41

// Read returns the value of the LCD status register
func (l *LCDStatus) Read(addr uint16) uint8 {
	if addr != 0 {
		panic("Invalid address")
	}

	switch addr {
	case 0x00:
		val := uint8(0)
		if l.LYCInterrupt {
			val |= 1 << 6
		}
		if l.Mode2Interrupt {
			val |= 1 << 5
		}
		if l.Mode1Interrupt {
			val |= 1 << 4
		}
		if l.Mode0Interrupt {
			val |= 1 << 3
		}
		if l.LYCMatch {
			val |= 1 << 2
		}
		val |= uint8(l.PPUMode)
		return val
	}

	panic("Invalid address")
}

// Write sets the value of the LCD status register
func (l *LCDStatus) Write(addr uint16, value uint8) {
	if addr != 0 {
		panic("Invalid address")
	}

	l.LYCInterrupt = value&(1<<6) > 0
	l.Mode2Interrupt = value&(1<<5) > 0
	l.Mode1Interrupt = value&(1<<4) > 0
	l.Mode0Interrupt = value&(1<<3) > 0
	l.LYCMatch = value&(1<<2) > 0
	l.PPUMode = PPUState(value & 0x3)
}
