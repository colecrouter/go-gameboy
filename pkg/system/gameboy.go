package system

import (
	"time"

	"github.com/colecrouter/gameboy-go/pkg/display/monochrome"
	"github.com/colecrouter/gameboy-go/pkg/memory"
	"github.com/colecrouter/gameboy-go/pkg/memory/registers"
	"github.com/colecrouter/gameboy-go/pkg/memory/vram"
	"github.com/colecrouter/gameboy-go/pkg/processor/cpu/lr35902"
	"github.com/colecrouter/gameboy-go/pkg/processor/ppu"
	"github.com/colecrouter/gameboy-go/pkg/reader/gamepak"
)

const CLOCK_SPEED = 4_194_304 // 4.194304 MHz
const CLOCK_DELAY = time.Second / CLOCK_SPEED
const DISPLAY_SPEED = 60 // 60 Hz
// https://gbdev.io/pandocs/Rendering.html
const DISPLAY_DELAY = time.Duration((float32(time.Second) * 1.0045) / DISPLAY_SPEED)

type GameBoy struct {
	display   *monochrome.TerminalDisplay
	memoryBus *memory.Bus
	io        *registers.Registers
	cpu       *lr35902.LR35902
	ppu       *ppu.PPU

	cpuTicker     *time.Ticker
	ppuTicker     *time.Ticker
	displayTicker *time.Ticker
	done          chan struct{}
	totalCycles   int64 // added to track CPU cycles
	// fastMode      bool  // new flag for fast execution (bypasses tickers)
}

func NewGameBoy() *GameBoy {
	gb := &GameBoy{}
	gb.memoryBus = &memory.Bus{}
	vramModule := &vram.VRAM{}
	oamModule := &memory.OAM{}
	gb.io = &registers.Registers{}
	gb.cpu = lr35902.NewLR35902(gb.memoryBus, gb.io)

	gb.cpuTicker = time.NewTicker(CLOCK_DELAY)
	gb.ppuTicker = time.NewTicker(CLOCK_DELAY)
	gb.displayTicker = time.NewTicker(DISPLAY_DELAY)
	gb.done = make(chan struct{}) // initialize done channel

	// gb.memoryBus.AddDevice(0x0000, 0x3FFF, &memory.Memory{Buffer: make([]byte, 0x4000)}) // ROM Bank 0
	// gb.memoryBus.AddDevice(0x4000, 0x7FFF, &memory.Memory{Buffer: make([]byte, 0x4000)}) // ROM Bank 1-xx aka mapper
	gb.memoryBus.AddDevice(0x8000, 0x9FFF, vramModule)                                   // VRAM
	gb.memoryBus.AddDevice(0xA000, 0xBFFF, &memory.Memory{Buffer: make([]byte, 0x2000)}) // External RAM
	gb.memoryBus.AddDevice(0xC000, 0xCFFF, &memory.Memory{Buffer: make([]byte, 0x1000)}) // WRAM
	gb.memoryBus.AddDevice(0xD000, 0xDFFF, &memory.Memory{Buffer: make([]byte, 0x1000)}) // WRAM
	gb.memoryBus.AddDevice(0xE000, 0xFDFF, &memory.Memory{Buffer: make([]byte, 0x1E00)}) // ECHO RAM
	gb.memoryBus.AddDevice(0xFE00, 0xFE9F, oamModule)                                    // OAM
	// https://gbdev.io/pandocs/Memory_Map.html#fea0feff-range
	gb.memoryBus.AddDevice(0xFEA0, 0xFEFF, &memory.Memory{Buffer: make([]byte, 0x60)}) // Unusable Memory
	gb.memoryBus.AddDevice(0xFF00, 0xFF7F, gb.io)                                      // I/O Registers
	gb.memoryBus.AddDevice(0xFF80, 0xFFFE, &memory.Memory{Buffer: make([]byte, 0x7F)}) // High RAM
	gb.memoryBus.AddDevice(0xFFFF, 0xFFFF, &memory.Memory{Buffer: make([]byte, 0x1)})  // Interrupt Enable Register

	gb.display = monochrome.NewTerminalDisplay()
	gb.ppu = ppu.NewPPU(vramModule, oamModule, gb.display, gb.io)

	return gb
}

func (gb *GameBoy) Start() {
	frameDuration := DISPLAY_DELAY   // ~1/60 second
	targetCycles := CLOCK_SPEED / 60 // cycles per frame
	for {
		select {
		case <-gb.done:
			return
		default:
			frameStart := time.Now()
			frameCycles := 0
			// Execute CPU instructions until reaching target cycles for this frame.
			for frameCycles < targetCycles {
				cycles := gb.cpu.Step()
				frameCycles += cycles
				gb.totalCycles += int64(cycles)
				for i := 0; i < cycles; i++ {
					gb.ppu.Clock()
				}
			}

			// Update the display
			go gb.display.Clock()
			// Sleep for the remainder of the frame, if any.
			elapsed := time.Since(frameStart)
			if remainder := frameDuration - elapsed; remainder > 0 {
				time.Sleep(remainder)
			}

			// Print VRAM tile range
			// gb.memoryBus.PrintMemory(0x8000, 0x9FFF)
		}
	}
	// }
}

func (gb *GameBoy) Stop() {
	gb.done <- struct{}{}

	gb.cpuTicker.Stop()
	gb.ppuTicker.Stop()
	gb.displayTicker.Stop()
}

func (gb *GameBoy) InsertCartridge(rom *gamepak.GamePak) {
	gb.memoryBus.AddDevice(0x0000, 0x7FFF, rom)
}
