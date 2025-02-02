package registers

type Registers struct {
	initialized bool

	JoypadState uint8          // 0xFF00
	SerialData  uint16         // 0xFF01-0xFF02
	Timer       uint32         // 0xFF04-0xFF07
	Interrupts  InterruptFlags // 0xFF0F
	Audio       uint32         // 0xFF10-0xFF26
	WavePattern uint16         // 0xFF30-0xFF3F
	LCDControl  LCDControl     // 0xFF40
	LCDStatus   LCDStatus      // 0xFF41-0xFF45
	// VRAMBank1 bool
	DisableBootROM bool // 0xFF50
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
	if !r.initialized {
		panic("Registers not initialized")
	}

	switch addr {
	case 0xFF00:
		return r.JoypadState
	case 0xFF01, 0xFF02:
		offset := addr - 0xFF01
		return uint8(r.SerialData >> (8 * offset))
	case 0xFF04, 0xFF05, 0xFF06, 0xFF07:
		offset := addr - 0xFF04
		return uint8(r.Timer >> (8 * offset))
	case 0xFF0F:
		return r.Interrupts.Read(addr)
	case 0xFF10, 0xFF11, 0xFF12, 0xFF13, 0xFF14, 0xFF16, 0xFF17, 0xFF18, 0xFF19, 0xFF1A, 0xFF1B, 0xFF1C, 0xFF1D, 0xFF1E, 0xFF20, 0xFF21, 0xFF22, 0xFF23, 0xFF24, 0xFF25, 0xFF26:
		offset := addr - 0xFF10
		return uint8(r.Audio >> (8 * offset))
	case 0xFF30, 0xFF31, 0xFF32, 0xFF33, 0xFF34, 0xFF35, 0xFF36, 0xFF37, 0xFF38, 0xFF39, 0xFF3A, 0xFF3B, 0xFF3C, 0xFF3D, 0xFF3E:
		offset := addr - 0xFF30
		return uint8(r.WavePattern >> (4 * offset))
	case 0xFF40:
		return r.LCDControl.Read(addr)
	case 0xFF41, 0xFF42, 0xFF43, 0xFF44, 0xFF45:
		return r.LCDStatus.Read(addr)
	case 0xFF50:
		if r.DisableBootROM {
			return 1
		}
		return 0
	}
	return 0
}

func (r *Registers) Write(addr uint16, value uint8) {
	if !r.initialized {
		panic("Registers not initialized")
	}

	switch addr {
	case 0xFF00:
		r.JoypadState = value
	case 0xFF01, 0xFF02:
		offset := addr - 0xFF01
		r.SerialData = uint16(value) << (8 * offset)
	case 0xFF04, 0xFF05, 0xFF06, 0xFF07:
		offset := addr - 0xFF04
		r.Timer = r.Timer | (uint32(value) << (8 * offset))
	case 0xFF0F:
		r.Interrupts.Write(addr, value)
	case 0xFF10, 0xFF11, 0xFF12, 0xFF13, 0xFF14, 0xFF16, 0xFF17, 0xFF18, 0xFF19, 0xFF1A, 0xFF1B, 0xFF1C, 0xFF1D, 0xFF1E, 0xFF20, 0xFF21, 0xFF22, 0xFF23, 0xFF24, 0xFF25, 0xFF26:
		offset := addr - 0xFF10
		r.Audio = r.Audio | (uint32(value) << (8 * offset))
	case 0xFF30, 0xFF31, 0xFF32, 0xFF33, 0xFF34, 0xFF35, 0xFF36, 0xFF37, 0xFF38, 0xFF39, 0xFF3A, 0xFF3B, 0xFF3C, 0xFF3D, 0xFF3E:
		offset := addr - 0xFF30
		r.WavePattern = r.WavePattern | (uint16(value) << (4 * offset))
	case 0xFF40:
		r.LCDControl.Write(addr, value)
	case 0xFF41, 0xFF42, 0xFF43, 0xFF44, 0xFF45:
		r.LCDStatus.Write(addr, value)
	case 0xFF50:
		r.DisableBootROM = value > 0
	}
}

func NewRegisters() *Registers {
	r := &Registers{}

	r.LCDStatus.PaletteData.Write(0xFC)

	r.initialized = true

	return r
}
