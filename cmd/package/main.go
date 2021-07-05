package main

import (
	"log"
	"os"

	"github.com/colecrouter/gameboy-go/pkg/reader/gamepak"
	"github.com/colecrouter/gameboy-go/pkg/system"
)

func main() {
	gb := system.NewGameBoy()

	f, err := os.ReadFile("tetris.gb")
	if err != nil {
		log.Fatalln(err.Error())
	}

	game := gamepak.NewGamePak(f)

	gb.InsertCartridge(game)

	gb.Start()
}
