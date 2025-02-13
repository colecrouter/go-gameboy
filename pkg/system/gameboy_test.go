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

// TestBlarggOutput loads the CPU instructions ROM and checks serial output.
func TestBlarggOutput(t *testing.T) {
	// Replace with the actual ROM path when available.
	romData, err := os.ReadFile("../../tests/blargg/cpu_instrs/individual/07-jr,jp,call,ret,rst.gb")
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

	// Wait for output to end with either "Failed\n" or "Passed\n".
	// Or wait for a timeout.

	var output string
	for {
		select {
		case <-time.Tick(1 * time.Second):
			output = string(testDevice.output)
			if strings.HasSuffix(output, "Failed\n") {
				t.Fatal("Test failed")
			} else if strings.HasSuffix(output, "Passed\n") {
				t.Log("Test passed")
				break
			}
		}
	}
}
