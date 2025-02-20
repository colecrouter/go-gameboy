package system

import (
	"log"
	"os"
	"runtime/pprof"
	"testing"
	"time"

	"github.com/colecrouter/gameboy-go/private/reader/gamepak"
)

func BenchmarkGameBoy_BootROMPerformance(b *testing.B) {
	f, err := os.ReadFile("../../tests/blargg/cpu_instrs/01-special.gb")
	if err != nil {
		log.Fatalln(err.Error())
	}

	game := gamepak.NewGamePak(f)
	gb := NewGameBoy()
	gb.FastMode = true
	gb.CartridgeReader.InsertCartridge(game)

	profFile, err := os.Create("cpu.prof")
	if err != nil {
		b.Fatal(err)
	}
	err = pprof.StartCPUProfile(profFile)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	go gb.Start(true)

	for {
		// Wait until CPU passes the boot ROM (PC > 0xFF).
		if gb.CPU.Registers.PC > 0xFF {
			gb.Stop()
			break
		}
	}
	b.Logf("Boot ROM completed. Cycles processed: %d", gb.totalCycles)

	pprof.StopCPUProfile()
	profFile.Close()
}

// Renamed existing benchmark to clarify its purpose
func BenchmarkGameBoy_CycleAccurateOneSecond(b *testing.B) {
	f, err := os.ReadFile("../../tests/blargg/cpu_instrs/cpu_instrs.gb")
	if err != nil {
		log.Fatalln(err.Error())
	}

	game := gamepak.NewGamePak(f)

	gb := NewGameBoy()
	// Don't enable FastMode for this benchmark
	// gb.fastMode = true
	gb.CartridgeReader.InsertCartridge(game)

	b.ResetTimer()

	// run GameBoy in a goroutine so that we can stop it after 1 second
	go gb.Start(false)

	// let emulator run for approximately 1 second
	time.Sleep(1 * time.Second)
	gb.Stop()

	expectedCycles := CLOCK_SPEED
	b.Logf("Total CPU cycles processed: %d, expected: ~%d", gb.totalCycles, expectedCycles)
	b.ReportMetric(float64(gb.totalCycles)/float64(expectedCycles), "speedFactor")
}

// Updated test for Blargg CPUInstrs using the reusable runner.
func TestGameBoy_BlarggCPUInstrs(t *testing.T) {
	dir := "../../tests/blargg/cpu_instrs/individual/"
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		romPath := dir + file.Name()
		t.Run(file.Name(), func(t *testing.T) {
			t.Parallel()
			RunBlarggTestRom(t, romPath)
		})
	}
}

// New test for instr_timing ROM
func TestGameBoy_BlarggInstrTiming(t *testing.T) {
	RunBlarggTestRom(t, "../../tests/blargg/instr_timing/instr_timing.gb")
}
