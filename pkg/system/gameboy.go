package system

import (
	"time"

	"github.com/colecrouter/gameboy-go/private/memory"
	"github.com/colecrouter/gameboy-go/private/memory/registers"
	"github.com/colecrouter/gameboy-go/private/memory/vram"
	"github.com/colecrouter/gameboy-go/private/processor/cpu/lr35902"
	"github.com/colecrouter/gameboy-go/private/processor/ppu"
	"github.com/colecrouter/gameboy-go/private/reader"
	"github.com/colecrouter/gameboy-go/private/reader/gamepak"
)

const CLOCK_SPEED = 4_194_304 // 4.194304 MHz
const DISPLAY_SPEED = 60

// const SPEED_MONITORING_FREQUENCY = 10

// https://gbdev.io/pandocs/Rendering.html
const FRAME_DURATION = time.Duration((float32(time.Second) * 1.0045) / DISPLAY_SPEED)
const TARGET_CYCLES_PER_FRAME = CLOCK_SPEED / DISPLAY_SPEED

type GameBoy struct {
	Bus             *memory.Bus
	IO              *registers.Registers
	CPU             *lr35902.LR35902
	PPU             *ppu.PPU
	VRAM            *vram.VRAM
	CartridgeReader reader.CartridgeReader
	IF              *registers.Interrupt
	IE              *registers.Interrupt

	done        chan struct{}
	totalCycles uint64 // added to track CPU cycles
	FastMode    bool
}

func NewGameBoy() *GameBoy {
	gb := &GameBoy{}
	gb.Bus = &memory.Bus{}
	gb.VRAM = &vram.VRAM{}
	gb.IF = &registers.Interrupt{}
	gb.IE = &registers.Interrupt{}
	gb.IO = registers.NewRegisters(gb.Bus, gb.IF)
	oamModule := memory.NewOAM(gb.VRAM, &gb.IO.LCDControl.Sprites8x16)
	gb.CPU = lr35902.NewLR35902(gb.Bus, gb.IO, gb.IE)
	gb.CartridgeReader = *reader.NewCartridgeReader(&gb.IO.DisableBootROM)

	gb.done = make(chan struct{}) // initialize done channel

	// gb.memoryBus.AddDevice(0x0000, 0x3FFF, &memory.Memory{Buffer: make([]byte, 0x4000)}) // ROM Bank 0
	// gb.memoryBus.AddDevice(0x4000, 0x7FFF, &memory.Memory{Buffer: make([]byte, 0x4000)}) // ROM Bank 1-xx aka mapper
	gb.Bus.AddDevice(0x0000, 0x7FFF, &gb.CartridgeReader)
	gb.Bus.AddDevice(0x8000, 0x9FFF, gb.VRAM)                                      // VRAM
	gb.Bus.AddDevice(0xA000, 0xBFFF, &memory.Memory{Buffer: make([]byte, 0x2000)}) // External RAM
	gb.Bus.AddDevice(0xC000, 0xCFFF, &memory.Memory{Buffer: make([]byte, 0x1000)}) // WRAM
	gb.Bus.AddDevice(0xD000, 0xDFFF, &memory.Memory{Buffer: make([]byte, 0x1000)}) // WRAM
	gb.Bus.AddDevice(0xE000, 0xFDFF, &memory.Memory{Buffer: make([]byte, 0x1E00)}) // ECHO RAM
	gb.Bus.AddDevice(0xFE00, 0xFE9F, oamModule)                                    // OAM
	// https://gbdev.io/pandocs/Memory_Map.html#fea0feff-range
	gb.Bus.AddDevice(0xFEA0, 0xFEFF, &memory.Memory{Buffer: make([]byte, 0x60)}) // Unusable Memory
	gb.Bus.AddDevice(0xFF00, 0xFF7F, gb.IO)                                      // I/O Registers
	gb.Bus.AddDevice(0xFF80, 0xFFFE, &memory.Memory{Buffer: make([]byte, 0x7F)}) // High RAM
	gb.Bus.AddDevice(0xFFFF, 0xFFFF, gb.IE)                                      // Interrupt Enable Register

	gb.PPU = ppu.NewPPU(gb.VRAM, oamModule, gb.IO, gb.IF)

	return gb
}

func (gb *GameBoy) Start(skip bool) {
	if skip {
		// Initialize registers to default DMG ROM values
		// https://gbdev.io/pandocs/Power_Up_Sequence.html#cpu-registers
		gb.CPU.Registers.PC = 0x0100

		gb.CPU.Registers.B = 0x00
		gb.CPU.Registers.A = 0x01
		gb.CPU.Registers.C = 0x13
		gb.CPU.Registers.D = 0x00
		gb.CPU.Registers.E = 0xD8
		gb.CPU.Registers.H = 0x01
		gb.CPU.Registers.L = 0x4D
		gb.CPU.Registers.SP = 0xFFFE

		gb.CPU.Flags.Write(0xB0)

		gb.IO.Write(0x00, 0xCF) // Joypad input
		gb.IO.Write(0x01, 0x00) // Serial transfer
		gb.IO.Write(0x02, 0x7E) // Serial transfer
		gb.IO.Write(0x04, 0xAB) // Timer and divider
		gb.IO.Write(0x05, 0x00) // Timer counter
		gb.IO.Write(0x06, 0x00) // Timer modulo
		gb.IO.Write(0x07, 0xF8) // Timer control
		gb.IO.Write(0x0F, 0xE1) // Interrupts

		// Audio
		gb.IO.Write(0x10, 0x80)
		gb.IO.Write(0x11, 0xBF)
		gb.IO.Write(0x12, 0xF3)
		gb.IO.Write(0x13, 0xFF)
		gb.IO.Write(0x14, 0xBF)
		gb.IO.Write(0x16, 0x3F)
		gb.IO.Write(0x17, 0x00)
		gb.IO.Write(0x18, 0xFF)
		gb.IO.Write(0x19, 0xBF)
		gb.IO.Write(0x1A, 0x7F)
		gb.IO.Write(0x1B, 0xFF)
		gb.IO.Write(0x1C, 0x9F)
		gb.IO.Write(0x1D, 0xFF)
		gb.IO.Write(0x1E, 0xBF)
		gb.IO.Write(0x20, 0xFF)
		gb.IO.Write(0x21, 0x00)
		gb.IO.Write(0x22, 0x00)
		gb.IO.Write(0x23, 0xBF)
		gb.IO.Write(0x24, 0x77)
		gb.IO.Write(0x25, 0xF3)
		gb.IO.Write(0x26, 0xF1)

		gb.IO.Write(0x40, 0x91) // LCD Control
		gb.IO.Write(0x41, 0x85) // LCD Status
		gb.IO.Write(0x42, 0x00) // Scroll Y
		gb.IO.Write(0x43, 0x00) // Scroll X
		gb.IO.Write(0x44, 0x00) // LY
		gb.IO.Write(0x45, 0x00) // LY Compare
		gb.IO.Write(0x46, 0xFF) // DMA
		gb.IO.Write(0x47, 0xFC) // Palette Data
		// gb.IO.Write(0x48, 0xFF) // Object Palette Data 1
		// gb.IO.Write(0x49, 0xFF) // Object Palette Data 2
		gb.IO.Write(0x4A, 0x00) // Window Y
		gb.IO.Write(0x4B, 0x00) // Window X
		// CGB only

		gb.IO.Write(0xFF, 0x00) // Interrupt Enable Register
	}

	for {
		frameStart := time.Now()

		select {
		case <-gb.done:
			return
		default:
			cycles := 0
			for cycles < TARGET_CYCLES_PER_FRAME {
				stepCycles := gb.CPU.Step()
				gb.totalCycles += uint64(stepCycles)
				cycles += stepCycles
				for i := 0; i < stepCycles; i++ {

					gb.PPU.SystemClock()
				}

				for i := 0; i < stepCycles/4; i++ {
					// Technicially this is an issue, because Timer.Clock() should happen before CPU.Step()
					gb.IO.Timer.MClock()
				}
			}
		}

		if !gb.FastMode {
			// Calculate remaining time for the frame.
			remaining := FRAME_DURATION - time.Since(frameStart)
			if remaining > 2*time.Millisecond {
				time.Sleep(remaining - 1*time.Millisecond)
			}
			// Busy wait for the final part of the frame.
			for time.Since(frameStart) < FRAME_DURATION {
			}
		}
	}
}

func (gb *GameBoy) Stop() {
	gb.done <- struct{}{}
}

func (gb *GameBoy) PC() uint16 {
	return gb.CPU.Registers.PC
}

func (gb *GameBoy) InsertCartridge(game *gamepak.GamePak) {
	gb.CartridgeReader.InsertCartridge(game)
}

func (gb *GameBoy) ConnectSerialDevice(d registers.SerialDevice) {
	gb.IO.Serial.Connect(d)
}

func (gb *GameBoy) Controller() *registers.JoyPad {
	return &gb.IO.JoypadState
}

func (gb *GameBoy) TotalCycles() uint64 {
	return gb.totalCycles
}
