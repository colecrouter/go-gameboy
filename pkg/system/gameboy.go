package system

import (
	"time"

	"github.com/colecrouter/gameboy-go/private/memory"
	"github.com/colecrouter/gameboy-go/private/memory/io"
	"github.com/colecrouter/gameboy-go/private/memory/vram"
	"github.com/colecrouter/gameboy-go/private/processor/cpu/lr35902"
	"github.com/colecrouter/gameboy-go/private/processor/ppu"
	"github.com/colecrouter/gameboy-go/private/reader"
	"github.com/colecrouter/gameboy-go/private/reader/gamepak"
	"github.com/colecrouter/gameboy-go/private/system"
)

const CLOCK_SPEED = 4_194_304 // 4.194304 MHz
const DISPLAY_SPEED = 60

// const SPEED_MONITORING_FREQUENCY = 10

// https://gbdev.io/pandocs/Rendering.html
const FRAME_DURATION = time.Duration((float32(time.Second) * 1.0045) / DISPLAY_SPEED)
const TARGET_CYCLES_PER_FRAME = CLOCK_SPEED / DISPLAY_SPEED

type GameBoy struct {
	Bus             *memory.Bus
	IO              *io.Registers
	CPU             *lr35902.LR35902
	PPU             *ppu.PPU
	VRAM            *vram.VRAM
	CartridgeReader reader.CartridgeReader
	IF              *io.Interrupt
	IE              *io.Interrupt

	done         chan struct{}
	totalTCycles uint64
	broadcaster  *system.Broadcaster
	FastMode     bool
}

func NewGameBoy() *GameBoy {
	gb := &GameBoy{}

	gb.broadcaster = system.NewBroadcaster()

	gb.Bus = &memory.Bus{}
	gb.VRAM = &vram.VRAM{}
	gb.IF = &io.Interrupt{}
	gb.IE = &io.Interrupt{}
	gb.IO = io.NewRegisters(gb.broadcaster, gb.Bus, gb.IF)
	oamModule := memory.NewOAM(gb.VRAM, &gb.IO.LCDControl.Sprites8x16)
	gb.CPU = lr35902.NewLR35902(gb.broadcaster, gb.Bus, gb.IO, gb.IE)
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

	gb.PPU = ppu.NewPPU(gb.broadcaster, gb.VRAM, oamModule, gb.IO, gb.IF)

	return gb
}

func (gb *GameBoy) Start(skip bool) {
	if skip {
		reg := gb.CPU.Registers()

		// Initialize registers to default DMG ROM values
		// https://gbdev.io/pandocs/Power_Up_Sequence.html#cpu-registers
		reg.PC = 0x0100

		reg.B = 0x00
		reg.A = 0x01
		reg.C = 0x13
		reg.D = 0x00
		reg.E = 0xD8
		reg.H = 0x01
		reg.L = 0x4D
		reg.SP = 0xFFFE

		fl := gb.CPU.Flags()
		fl.Write(0xB0)

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

	go gb.CPU.Run(gb.done)
	go gb.PPU.Run(gb.done)
	go gb.IO.Timer.Run(gb.done)

	for {
		frameStart := time.Now()

		select {
		case <-gb.done:
			return
		default:
			for range TARGET_CYCLES_PER_FRAME {
				for i := range 4 {
					if i == 0 {
						gb.broadcaster.BroadcastM()
					}

					gb.broadcaster.BroadcastT()
				}
				gb.totalTCycles += 4
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
	close(gb.done)
}

func (gb *GameBoy) PC() uint16 {
	return gb.CPU.Registers().PC
}

func (gb *GameBoy) InsertCartridge(game *gamepak.GamePak) {
	gb.CartridgeReader.InsertCartridge(game)
}

func (gb *GameBoy) ConnectSerialDevice(d io.SerialDevice) {
	gb.IO.Serial.Connect(d)
}

func (gb *GameBoy) Controller() *io.JoyPad {
	return &gb.IO.JoypadState
}

func (gb *GameBoy) TotalCycles() uint64 {
	return gb.totalTCycles
}
