package system

import (
	"time"

	"github.com/colecrouter/gameboy-go/pkg/memory"
	"github.com/colecrouter/gameboy-go/pkg/memory/registers"
	"github.com/colecrouter/gameboy-go/pkg/memory/vram"
	"github.com/colecrouter/gameboy-go/pkg/processor/cpu/lr35902"
	"github.com/colecrouter/gameboy-go/pkg/processor/ppu"
	"github.com/colecrouter/gameboy-go/pkg/reader/gamepak"
)

const CLOCK_SPEED = 4_194_304 // 4.194304 MHz
const DISPLAY_SPEED = 60

// const SPEED_MONITORING_FREQUENCY = 10

// https://gbdev.io/pandocs/Rendering.html
const FRAME_DURATION = time.Duration((float32(time.Second) * 1.0045) / DISPLAY_SPEED)
const TARGET_CYCLES_PER_FRAME = CLOCK_SPEED / DISPLAY_SPEED

type GameBoy struct {
	Bus  *memory.Bus
	IO   *registers.Registers
	CPU  *lr35902.LR35902
	PPU  *ppu.PPU
	VRAM *vram.VRAM

	done        chan struct{}
	totalCycles int64 // added to track CPU cycles
}

func NewGameBoy() *GameBoy {
	gb := &GameBoy{}
	gb.Bus = &memory.Bus{}
	gb.VRAM = &vram.VRAM{}
	oamModule := &memory.OAM{}
	gb.IO = &registers.Registers{}
	gb.CPU = lr35902.NewLR35902(gb.Bus, gb.IO)

	gb.done = make(chan struct{}) // initialize done channel

	// gb.memoryBus.AddDevice(0x0000, 0x3FFF, &memory.Memory{Buffer: make([]byte, 0x4000)}) // ROM Bank 0
	// gb.memoryBus.AddDevice(0x4000, 0x7FFF, &memory.Memory{Buffer: make([]byte, 0x4000)}) // ROM Bank 1-xx aka mapper
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
	gb.Bus.AddDevice(0xFFFF, 0xFFFF, &memory.Memory{Buffer: make([]byte, 0x1)})  // Interrupt Enable Register

	gb.PPU = ppu.NewPPU(gb.VRAM, oamModule, gb.IO)

	return gb
}

func (gb *GameBoy) Start() {
	for {
		frameStart := time.Now()

		select {
		case <-gb.done:
			return
		default:
			cycles := 0
			for cycles < TARGET_CYCLES_PER_FRAME {
				stepCycles := gb.CPU.Step()
				gb.totalCycles += int64(stepCycles)
				cycles += stepCycles / 2
				// Run the PPU clock for each CPU cycle executed.
				for i := 0; i < stepCycles; i++ {
					gb.PPU.SystemClock()
				}
			}
		}

		// Check how long the frame took to process.
		frameElapsed := time.Since(frameStart)
		if frameElapsed < FRAME_DURATION {
			time.Sleep(FRAME_DURATION - frameElapsed)
		}
	}
}

func (gb *GameBoy) Stop() {
	gb.done <- struct{}{}
}

func (gb *GameBoy) InsertCartridge(rom *gamepak.GamePak) {
	gb.Bus.AddDevice(0x0000, 0x7FFF, rom)
}
