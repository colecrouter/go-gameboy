package system

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/colecrouter/gameboy-go/pkg/reader/gamepak"
)

func BenchmarkGameBoy_OneSecondClockCycles(b *testing.B) {
	f, err := os.ReadFile("../../tetris.gb")
	if err != nil {
		log.Fatalln(err.Error())
	}

	game := gamepak.NewGamePak(f)
	game.Title()

	gb := NewGameBoy()
	gb.InsertCartridge(game)

	// Record start time before running all benchmark iterations.
	start := time.Now()
	for i := 0; i < b.N; i++ {
		for cycle := 0; cycle < CLOCK_SPEED; cycle++ {
			gb.cpu.Clock()
		}
		// Removed per-iteration logging.
	}
	// Compute elapsed real time (for debugging if needed)
	elapsed := time.Since(start)
	// Each iteration simulates 1.0 second.
	simulatedTotalTime := float64(b.N)
	b.Logf("Simulated total GameBoy time: %.2f seconds, benchmark elapsed time: %v", simulatedTotalTime, elapsed)
	// Report custom metric if desired:
	b.ReportMetric(simulatedTotalTime, "simTime/s")
}
