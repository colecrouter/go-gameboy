package registers

import (
	"fmt"

	"github.com/colecrouter/gameboy-go/pkg/memory"
	"github.com/colecrouter/gameboy-go/pkg/reader"
)

type Registers struct {
	initialized bool
	oam         *memory.OAM
	cartridge   *reader.CartridgeReader

	JoypadState        uint8          // 0xFF00
	SerialData         uint16         // 0xFF01-0xFF02
	Timer              Timer          // 0xFF04-0xFF07
	Interrupts         InterruptFlags // 0xFF0F
	Audio              uint32         // 0xFF10-0xFF26
	WavePattern        uint16         // 0xFF30-0xFF3F
	LCDControl         LCDControl     // 0xFF40
	LCDStatus          LCDStatus      // 0xFF41-0xFF45
	ScrollY            uint8          // 0xFF42
	ScrollX            uint8          // 0xFF43
	LY                 uint8          // 0xFF44
	LYCompare          uint8          // 0xFF45
	DMA                uint8          // 0xFF46
	PaletteData        Palette        // 0xFF47
	ObjectPaletteData1 Palette        // 0xFF48
	ObjectPaletteData2 Palette        // 0xFF49
	WindowY            uint8          // 0xFF4A
	WindowX            uint8          // 0xFF4B

	// CGB only
	// TODO
	VRAMBank1      bool     // 0xFF4F
	DisableBootROM bool     // 0xFF50
	VRAMDMA        [5]uint8 // 0xFF51-0xFF55
	WRAMBank       uint8    // 0xFF70
	GBCPaletteData [8]uint8 // 0xFF68-0xFF6B
	WRAMBank1      bool     // 0xFF70

	// Interrupts
	InterruptRequest Interrupt // 0xFFF0
	InterruptEnable  Interrupt // 0xFFFF

	// ???
	rest [0x8A]uint8
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

func NewRegisters(oam *memory.OAM, cartridge *reader.CartridgeReader) *Registers {
	return &Registers{
		oam:         oam,
		cartridge:   cartridge,
		initialized: true,
	}
}

func (r *Registers) Read(addr uint16) uint8 {
	if !r.initialized {
		panic("Registers not initialized")
	}

	switch addr {
	case 0x00:
		return r.JoypadState
	case 0x01, 0x02:
		return 0
	case 0x04, 0x05, 0x06, 0x07:
		offset := addr - 0x04
		return r.Timer.Read(offset)
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
	case 0x41:
		return r.LCDStatus.Read(0)
	case 0x42:
		return r.ScrollY
	case 0x43:
		return r.ScrollX
	case 0x44:
		return r.LY
	case 0x45:
		return r.LYCompare
	case 0x46:
		return r.DMA
	case 0x47:
		return r.PaletteData.Read(0)
	case 0x48:
		return r.ObjectPaletteData1.Read(0)
	case 0x49:
		return r.ObjectPaletteData2.Read(0)
	case 0x4A:
		return r.WindowY
	case 0x4B:
		return r.WindowX
	case 0x50:
		if r.DisableBootROM {
			return 1
		}
		return 0
	case 0x4F:
		if r.VRAMBank1 {
			return 1
		}
		return 0
	case 0x68, 0x69, 0x6A, 0x6B:
		offset := addr - 0x68
		return r.GBCPaletteData[offset]
	case 0x51, 0x52, 0x53, 0x54, 0x55:
		offset := addr - 0x51
		return r.VRAMDMA[offset]
	case 0x70:
		if r.WRAMBank1 {
			return 1
		}
		return 0
	case 0xF0:
		return r.InterruptRequest.Read()
	case 0xFF:
		return r.InterruptEnable.Read()
	default:
		return 0
	}
}

func (r *Registers) Write(addr uint16, value uint8) {
	if !r.initialized {
		panic("Registers not initialized")
	}

	switch addr {
	case 0x00:
		r.JoypadState = value
	case 0x01, 0x02:
		offset := addr - 0x01
		r.SerialData = uint16(value) << (8 * offset)
	case 0x04, 0x05, 0x06, 0x07:
		offset := addr - 0x04
		r.Timer.Write(offset, value)
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
	case 0x41:
		r.LCDStatus.Write(0, value)
	case 0x42:
		r.ScrollY = value
	case 0x43:
		r.ScrollX = value
	case 0x44:
		panic("LY is read-only")
	case 0x45:
		r.LYCompare = value
	case 0x46:
		r.dma(value)
	case 0x47:
		r.PaletteData.Write(0, value)
	case 0x48:
		r.ObjectPaletteData1.Write(0, value)
	case 0x49:
		r.ObjectPaletteData2.Write(0, value)
	case 0x4A:
		r.WindowY = value
	case 0x4B:
		r.WindowX = value
	case 0x50:
		r.DisableBootROM = value > 0
	case 0x4F:
		r.VRAMBank1 = value > 0
	case 0x68, 0x69, 0x6A, 0x6B:
		offset := addr - 0x68
		r.GBCPaletteData[offset] = value
	case 0x51, 0x52, 0x53, 0x54, 0x55:
		offset := addr - 0x51
		r.VRAMDMA[offset] = value
	case 0x70:
		r.WRAMBank1 = value > 0
	case 0xF0:
		r.InterruptRequest.Write(value)
	case 0xFF:
		r.InterruptEnable.Write(value)
	default:
		if addr >= 0x71 && addr <= 0xFF {
			r.rest[addr-0x71] = value
			return
		}
		panic("Invalid register address")
	}
}

func (r *Registers) dma(addr uint8) {
	fmt.Printf("DMA transfer from 0x%02X00\n", addr)
	source := uint16(addr) << 8 // Source address is * 100
	for i := 0; i < 0xA0; i++ {
		// Add 0x4000 to cancel the bus’ subtraction.
		r.oam.Write(uint16(i), r.cartridge.Read(source+0x4000+uint16(i)))
	}
}
