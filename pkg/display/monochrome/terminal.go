package monochrome

import (
	"fmt"

	"github.com/colecrouter/gameboy-go/pkg/display"
	"github.com/colecrouter/gameboy-go/pkg/memory"
)

const WIDTH = 160
const HEIGHT = 144

// const WIDTH = 8
// const HEIGHT = 8
const TILES_WIDTH = WIDTH / 8
const TILES_HEIGHT = HEIGHT / 8

type TerminalDisplay struct {
	initialized bool
	vram        *memory.Device
	buffer      [HEIGHT][WIDTH]uint8
}

func (t *TerminalDisplay) populateBuffer() {
	// First, construct an intermediate buffer of tiles
	tiles := [TILES_HEIGHT][TILES_WIDTH]display.Tile{}
	for y := 0; y < TILES_HEIGHT; y++ {
		for x := 0; x < TILES_WIDTH; x++ {
			tileAddr := uint16(y*TILES_WIDTH+x) * 16
			tile := display.Tile{}
			for row := uint8(0); row < 8; row++ {
				tile.Bytes[row*2] = (*t.vram).Read(tileAddr + uint16(row)*2)
				tile.Bytes[row*2+1] = (*t.vram).Read(tileAddr + uint16(row)*2 + 1)
			}
			tiles[y][x] = tile
		}
	}

	// Then, draw the tiles to the screen
	for y := 0; y < HEIGHT; y++ {
		for x := 0; x < WIDTH; x++ {
			tileX := x / 8
			tileY := y / 8
			tile := tiles[tileY][tileX]
			color := tile.ReadPixel(uint8(y%8), uint8(x%8))
			t.buffer[y][x] = color
		}
	}
}

func (t *TerminalDisplay) Clock() {
	if !t.initialized {
		panic("TerminalDisplay not initialized")
	}

	t.populateBuffer()

	// Clear the screen
	fmt.Print("\033c")

	// Print the buffer to the terminal
	for y := 0; y < HEIGHT; y++ {
		for x := 0; x < WIDTH; x++ {
			color := t.buffer[y][x]
			switch color {
			case 0:
				print("  ")
			case 1:
				print("░░")
			case 2:
				print("▒▒")
			case 3:
				print("▓▓")
			}
		}
		print("\n")
	}
}

func NewTerminalDisplay(vram memory.Device) *TerminalDisplay {
	t := &TerminalDisplay{initialized: true}

	t.vram = &vram

	for y := 0; y < HEIGHT; y++ {
		for x := 0; x < WIDTH; x++ {
			t.buffer[y][x] = 0
		}
	}

	return t
}
