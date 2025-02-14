package system

import (
	"log"
	"os"
	"runtime/pprof"
	"strings"
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

	go gb.Start()

	for {
		// Wait until CPU passes the boot ROM (PC > 0xFF).
		if gb.CPU.PC() > 0xFF {
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
	go gb.Start()

	// let emulator run for approximately 1 second
	time.Sleep(1 * time.Second)
	gb.Stop()

	expectedCycles := CLOCK_SPEED
	b.Logf("Total CPU cycles processed: %d, expected: ~%d", gb.totalCycles, expectedCycles)
	b.ReportMetric(float64(gb.totalCycles)/float64(expectedCycles), "speedFactor")
}

// TestBlarggOutput runs the cpu_instrs test ROMs from Blargg's test suite.
func TestGameBoy_BlarggCPUInstrs(t *testing.T) {
	// Get all roms in ./tests/blargg/cpu_instrs/individual/*.gb

	dir := "../../tests/blargg/cpu_instrs/individual/"
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		f := file // capture file variable
		t.Run(f.Name(), func(t *testing.T) {
			t.Parallel()

			// Replace with the actual ROM path when available.
			romData, err := os.ReadFile("../../tests/blargg/cpu_instrs/individual/" + f.Name())
			if err != nil {
				t.Fatal(err)
			}
			game := gamepak.NewGamePak(romData)
			// Use the test setup that connects the test serial device.
			gb, testDevice := SetupBlarggTestSystem()
			// Insert cartridge and any setup needed.
			gb.CartridgeReader.InsertCartridge(game)

			go gb.Start()
			defer gb.Stop()

			// Use a ticker for periodic checks and a timeout channel.
			ticker := time.NewTicker(1 * time.Second)
			defer ticker.Stop()

			for range ticker.C {
				output := string(testDevice.output)
				if strings.Contains(output, "Failed") {
					t.Fatal("Test failed")
				} else if strings.Contains(output, "Passed") {
					t.Log("Test passed")
					return
				}
			}
		})
	}
}
