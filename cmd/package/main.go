package main

import (
	"log"
	"os"

	"github.com/colecrouter/gameboy-go/pkg/system"
	"github.com/colecrouter/gameboy-go/private/reader/gamepak"
	"github.com/colecrouter/gameboy-go/private/ui/terminal"
)

func main() {
	gb := system.NewGameBoy()

	// romData, err := os.ReadFile("tetris.gb")
	romData, err := os.ReadFile("./tests/blargg/cpu_instrs/individual/01-special.gb")
	if err != nil {
		log.Fatalln(err)
	}
	game := gamepak.NewGamePak(romData)
	gb.CartridgeReader.InsertCartridge(game)

	app := terminal.NewApplication(gb)

	app.Run()
}
