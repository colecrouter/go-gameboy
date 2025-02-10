package main

import (
	"log"
	"os"

	"github.com/colecrouter/gameboy-go/pkg/reader/gamepak"
	"github.com/colecrouter/gameboy-go/pkg/system"
	"github.com/colecrouter/gameboy-go/pkg/ui/terminal"
)

func main() {
	gb := system.NewGameBoy()

	romData, err := os.ReadFile("tetris.gb")
	// romData, err := os.ReadFile("01-special.gb")
	if err != nil {
		log.Fatalln(err)
	}
	game := gamepak.NewGamePak(romData)
	gb.CartridgeReader.InsertCartridge(game)

	app := terminal.NewApplication(gb)

	app.Run()
}
