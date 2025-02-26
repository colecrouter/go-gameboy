package system

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/colecrouter/gameboy-go/private/reader/gamepak"
)

// BenchmarkGameBoy_CycleAccurateOneSecond benchmarks the GameBoy's performance by running it for 1 second
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
	b.Logf("Total CPU cycles processed: %d, expected: ~%d", gb.totalTCycles, expectedCycles)
	b.ReportMetric(float64(gb.totalTCycles)/float64(expectedCycles), "speedFactor")
}

// TestGameBoy_BlarggCPUInstrs runs the Blargg CPU instruction tests
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

// TestGameBoy_BlarggInstrTiming runs the Blargg instruction timing test
func TestGameBoy_BlarggInstrTiming(t *testing.T) {
	RunBlarggTestRom(t, "../../tests/blargg/instr_timing/instr_timing.gb")
}
