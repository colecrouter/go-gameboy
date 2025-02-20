package system

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/colecrouter/gameboy-go/private/reader/gamepak"
)

// RunBlarggTestRom executes a test ROM and waits for a pass/fail result.
func RunBlarggTestRom(t *testing.T, romPath string) {
	// Read ROM file
	romData, err := os.ReadFile(romPath)
	if err != nil {
		t.Fatal(err)
	}
	game := gamepak.NewGamePak(romData)

	// Setup test system with serial device.
	gb, testDevice := SetupBlarggTestSystem()
	gb.CartridgeReader.InsertCartridge(game)

	go gb.Start(true)
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
}
