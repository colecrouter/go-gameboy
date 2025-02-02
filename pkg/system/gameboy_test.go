package system

import (
	"log"
	"os"
	"runtime/pprof"
	"testing"
	"time"

	"github.com/colecrouter/gameboy-go/pkg/reader/gamepak"
)

func BenchmarkGameBoy_OneSecondClockCycles(b *testing.B) {
	// Start CPU profiling
	profFile, err := os.Create("cpu.prof")
	if err != nil {
		b.Fatal(err)
	}
	err = pprof.StartCPUProfile(profFile)
	if err != nil {
		b.Fatal(err)
	}
	defer func() {
		pprof.StopCPUProfile()
		profFile.Close()
	}()

	f, err := os.ReadFile("../../tetris.gb")
	if err != nil {
		log.Fatalln(err.Error())
	}

	game := gamepak.NewGamePak(f)
	game.Title()

	gb := NewGameBoy()
	// // Enable fast mode to bypass real-time delays during benchmarking
	// gb.fastMode = true
	gb.InsertCartridge(game)

	// run GameBoy in a goroutine so that we can stop it after 1 second
	go gb.Start()

	// let emulator run for approximately 1 second
	time.Sleep(1 * time.Second)
	gb.Stop()

	expectedCycles := CLOCK_SPEED
	b.Logf("Total CPU cycles processed: %d, expected: ~%d", gb.totalCycles, expectedCycles)
	b.ReportMetric(float64(gb.totalCycles)/float64(expectedCycles), "speedFactor")
}
