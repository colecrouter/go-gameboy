package registers

type Registers struct {
	JoypadState uint8          // 0x00
	SerialData  uint16         // 0x01-0x02
	Timer       uint32         // 0x04-0x07
	Interrupts  InterruptFlags // 0x0F
	Audio       uint32         // 0x10-0x26
	WavePattern uint16         // 0x30-0x3F
	LCDControl  LCDControl     // 0x40
	LCDStatus   LCDStatus      // 0x41-0x45
	// DMA 0x46
	// TODO http://www.codeslinger.co.uk/pages/projects/gameboy/dma.html
	PaletteData Palette // 0x47
	PositionY   uint8   // 0x4A
	PositionX   uint8   // 0x4B
	// VRAMBank1 bool
	DisableBootROM bool // 0x50
	// VRAMDMA uint8
}

/*
$FF00		DMG	Joypad input
$FF01	$FF02	DMG	Serial transfer
$FF04	$FF07	DMG	Timer and divider
$FF0F		DMG	Interrupts
$FF10	$FF26	DMG	Audio
$FF30	$FF3F	DMG	Wave pattern
$FF40	$FF4B	DMG	LCD Control, Status, Position, Scrolling, and Palettes
$FF4F		CGB	VRAM Bank Select
$FF50		DMG	Set to non-zero to disable boot ROM
$FF51	$FF55	CGB	VRAM DMA
$FF68	$FF6B	CGB	BG / OBJ Palettes
$FF70		CGB	WRAM Bank Select
*/

func (r *Registers) Read(addr uint16) uint8 {
	addr &= 0xFF

	switch addr {
	case 0x00:
		return r.JoypadState
	case 0x01, 0x02:
		offset := addr - 0x01
		return uint8(r.SerialData >> (8 * offset))
	case 0x04, 0x05, 0x06, 0x07:
		offset := addr - 0x04
		return uint8(r.Timer >> (8 * offset))
	case 0x0F:
		return r.Interrupts.Read(addr)
	case 0x10, 0x11, 0x12, 0x13, 0x14, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26:
		offset := addr - 0x10
		return uint8(r.Audio >> (8 * offset))
	case 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3A, 0x3B, 0x3C, 0x3D, 0x3E:
		offset := addr - 0x30
		return uint8(r.WavePattern >> (4 * offset))
	case 0x40:
		return r.LCDControl.Read(addr - 0x40)
	case 0x41, 0x42, 0x43, 0x44, 0x45:
		return r.LCDStatus.Read(addr - 0x41)
	case 0x47:
		return r.PaletteData.Read(0)
	case 0x4A:
		return r.PositionY
	case 0x4B:
		return r.PositionX
	case 0x50:
		if r.DisableBootROM {
			return 1
		}
		return 0
	}
	panic("Invalid register address")
}

func (r *Registers) Write(addr uint16, value uint8) {
	addr &= 0xFF

	switch addr {
	case 0x00:
		r.JoypadState = value
	case 0x01, 0x02:
		offset := addr - 0x01
		r.SerialData = uint16(value) << (8 * offset)
	case 0x04, 0x05, 0x06, 0x07:
		offset := addr - 0x04
		r.Timer = r.Timer | (uint32(value) << (8 * offset))
	case 0x0F:
		r.Interrupts.Write(0, value)
	case 0x10, 0x11, 0x12, 0x13, 0x14, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26:
		offset := addr - 0x10
		r.Audio = r.Audio | (uint32(value) << (8 * offset))
	case 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3A, 0x3B, 0x3C, 0x3D, 0x3E:
		offset := addr - 0x30
		r.WavePattern = r.WavePattern | (uint16(value) << (4 * offset))
	case 0x40:
		r.LCDControl.Write(0, value)
	case 0x41, 0x42, 0x43, 0x44, 0x45:
		r.LCDStatus.Write(addr-0x41, value)
	case 0x47:
		r.PaletteData.Write(0, value)
	case 0x4A:
		r.PositionY = value
	case 0x4B:
		r.PositionX = value
	case 0x50:
		r.DisableBootROM = value > 0
	default:
		panic("Invalid register address")
	}
}
