package main

import (
	"log"
	"os"
	"runtime/pprof"
	"testing"
	"time"

	"github.com/colecrouter/gameboy-go/pkg/system"
	"github.com/colecrouter/gameboy-go/private/reader/gamepak"
	"github.com/colecrouter/gameboy-go/private/ui/terminal"
)

func BenchmarkGameBoy_BootROMPerformance(b *testing.B) {
	stdout := os.Stdout
	null, _ := os.Open(os.DevNull)

	// Kill stdout from the terminal application.
	os.Stdout = null
	defer func() {
		recover()

		os.Stdout = stdout
	}()

	f, err := os.ReadFile("../../tests/blargg/cpu_instrs/cpu_instrs.gb")
	if err != nil {
		log.Fatalln(err.Error())
	}

	for b.Loop() {

		game := gamepak.NewGamePak(f)
		gb := system.NewGameBoy()
		gb.FastMode = true
		gb.CartridgeReader.InsertCartridge(game)
		app := terminal.NewApplication(gb)

		var profFile *os.File
		// if b.N == 0 {
		// Start CPU profiling.
		profFile, err := os.Create("cpu.prof")
		if err != nil {
			b.Fatal(err)
		}
		err = pprof.StartCPUProfile(profFile)
		if err != nil {
			b.Fatal(err)
		}
		// }

		b.ResetTimer()

		go app.Run(false)

		// for {
		// 	// Wait until CPU passes the boot ROM (PC > 0xFF).
		// 	if gb.CPU.Registers.PC > 0xFF {
		// 		gb.Stop()
		// 		break
		// 	}
		// }

		// Wait for 5 seconds.
		time.Sleep(5 * time.Second)
		gb.Stop()

		// Temporarily restore stdout.
		os.Stdout = stdout
		b.Logf("Boot ROM completed. Cycles processed: %d", gb.TotalCycles())
		os.Stdout = null

		// if b.N == 0 {
		// Stop CPU profiling.
		pprof.StopCPUProfile()
		profFile.Close()
		// }

	}

	// Restore stdout.
	os.Stdout = stdout

	b.Logf("Benchmark complete; run \"go tool pprof -http :12345 -no_browser ./cmd/package/cpu.prof\" to view the profile.\n")
}
