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

type GameBoy struct {
	display   monochrome.TerminalDisplay
	cpu       *lr35902.LR35902
	memoryBus *memory.Bus

	cpuTicker     *time.Ticker
	displayTicker *time.Ticker
}

func NewGameBoy() GameBoy {
	gb := GameBoy{}

	gb.cpu = lr35902.NewLR35902(nil)

	// gb.memoryBus.AddDevice(0x0000, 0x3FFF, &memory.Memory{Buffer: make([]byte, 0x4000)}) // ROM Bank 0
	// gb.memoryBus.AddDevice(0x4000, 0x7FFF, &memory.Memory{Buffer: make([]byte, 0x4000)}) // ROM Bank 1-xx
	gb.memoryBus.AddDevice(0x8000, 0x9FFF, &vram.VRAM{})                                 // VRAM
	gb.memoryBus.AddDevice(0xA000, 0xBFFF, &memory.Memory{Buffer: make([]byte, 0x2000)}) // External RAM
	gb.memoryBus.AddDevice(0xC000, 0xCFFF, &memory.Memory{Buffer: make([]byte, 0x1000)}) // WRAM
	gb.memoryBus.AddDevice(0xD000, 0xDFFF, &memory.Memory{Buffer: make([]byte, 0x1000)}) // WRAM
	gb.memoryBus.AddDevice(0xE000, 0xFDFF, &memory.Memory{Buffer: make([]byte, 0x1E00)}) // ECHO RAM
	gb.memoryBus.AddDevice(0xFE00, 0xFE9F, &memory.Memory{Buffer: make([]byte, 0xA0)})   // OAM
	gb.memoryBus.AddDevice(0xFEA0, 0xFEFF, nil)                                          // Unusable Memory
	gb.memoryBus.AddDevice(0xFF00, 0xFF7F, &memory.Memory{Buffer: make([]byte, 0x80)})   // I/O Registers
	gb.memoryBus.AddDevice(0xFF80, 0xFFFE, &memory.Memory{Buffer: make([]byte, 0x7F)})   // High RAM
	gb.memoryBus.AddDevice(0xFFFF, 0xFFFF, nil)                                          // Interrupt Enable Register

	return gb
}

func (gb *GameBoy) Start() {
	gb.cpuTicker = time.NewTicker(CLOCK_DELAY)
	for range gb.cpuTicker.C {
		gb.cpu.Clock()
	}
}

func (gb *GameBoy) Stop() {
	gb.cpuTicker.Stop()
}

func (gb *GameBoy) Step() {
	gb.cpu.Clock()
}

func (gb *GameBoy) Reset() {
	gb.Stop()
	gb.cpu = lr35902.NewLR35902(gb.memoryBus)
}

func (gb *GameBoy) InsertCartridge(rom *gamepak.GamePak) {
	gb.memoryBus.AddDevice(0x0000, 0x7FFF, rom)
}
