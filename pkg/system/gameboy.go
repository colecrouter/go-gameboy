package system

import (
	"time"

	"github.com/colecrouter/gameboy-go/pkg/cpu/lr35902"
	"github.com/colecrouter/gameboy-go/pkg/display/monochrome"
	"github.com/colecrouter/gameboy-go/pkg/memory"
	"github.com/colecrouter/gameboy-go/pkg/memory/vram"
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
	cpu       *lr35902.LR35902

	cpuTicker     *time.Ticker
	displayTicker *time.Ticker
	done          chan struct{}
}

func NewGameBoy() *GameBoy {
	gb := &GameBoy{}
	gb.memoryBus = &memory.Bus{}

	gb.cpu = lr35902.NewLR35902(gb.memoryBus)

	vramModule := &vram.VRAM{}

	// gb.memoryBus.AddDevice(0x0000, 0x3FFF, &memory.Memory{Buffer: make([]byte, 0x4000)}) // ROM Bank 0
	// gb.memoryBus.AddDevice(0x4000, 0x7FFF, &memory.Memory{Buffer: make([]byte, 0x4000)}) // ROM Bank 1-xx
	gb.memoryBus.AddDevice(0x8000, 0x9FFF, vramModule)                                   // VRAM
	gb.memoryBus.AddDevice(0xA000, 0xBFFF, &memory.Memory{Buffer: make([]byte, 0x2000)}) // External RAM
	gb.memoryBus.AddDevice(0xC000, 0xCFFF, &memory.Memory{Buffer: make([]byte, 0x1000)}) // WRAM
	gb.memoryBus.AddDevice(0xD000, 0xDFFF, &memory.Memory{Buffer: make([]byte, 0x1000)}) // WRAM
	gb.memoryBus.AddDevice(0xE000, 0xFDFF, &memory.Memory{Buffer: make([]byte, 0x1E00)}) // ECHO RAM
	gb.memoryBus.AddDevice(0xFE00, 0xFE9F, &memory.Memory{Buffer: make([]byte, 0xA0)})   // OAM
	// https://gbdev.io/pandocs/Memory_Map.html#fea0feff-range
	gb.memoryBus.AddDevice(0xFEA0, 0xFEFF, &memory.Memory{Buffer: make([]byte, 0x60)}) // Unusable Memory
	gb.memoryBus.AddDevice(0xFF00, 0xFF7F, &memory.Memory{Buffer: make([]byte, 0x80)}) // I/O Registers
	gb.memoryBus.AddDevice(0xFF80, 0xFFFE, &memory.Memory{Buffer: make([]byte, 0x7F)}) // High RAM
	gb.memoryBus.AddDevice(0xFFFF, 0xFFFF, &memory.Memory{Buffer: make([]byte, 0x1)})  // Interrupt Enable Register

	gb.display = monochrome.NewTerminalDisplay(vramModule)

	return gb
}

func (gb *GameBoy) Start() {
	gb.cpuTicker = time.NewTicker(CLOCK_DELAY)

	go func() {
		for {
			select {
			case <-gb.done:
				return
			case <-gb.cpuTicker.C:
				gb.cpu.Clock()
			}
		}
	}()

	gb.displayTicker = time.NewTicker(DISPLAY_DELAY)
	for {
		select {
		case <-gb.done:
			return
		case <-gb.displayTicker.C:
			gb.display.Clock()
		}
	}
}

func (gb *GameBoy) Stop() {
	gb.cpuTicker.Stop()
	gb.displayTicker.Stop()

	gb.done <- struct{}{}
}

func (gb *GameBoy) Reset() {
	gb.Stop()
	gb.cpu = lr35902.NewLR35902(gb.memoryBus)
}

func (gb *GameBoy) InsertCartridge(rom *gamepak.GamePak) {
	gb.memoryBus.AddDevice(0x0000, 0x7FFF, rom)
}
