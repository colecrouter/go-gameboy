package monochrome

import (
	"fmt"
)

const WIDTH = 160
const HEIGHT = 144
const TILES_WIDTH = WIDTH / 8
const TILES_HEIGHT = HEIGHT / 8

type TerminalDisplay struct {
	initialized  bool
	buffer       [HEIGHT][WIDTH]uint8
	stringBuffer string
}

func (t *TerminalDisplay) DrawScanline(row uint8, line []uint8) {
	for x := 0; x < WIDTH; x++ {
		color := line[x]
		t.buffer[row][x] = color
	}
}

func (t *TerminalDisplay) Clock() {
	if !t.initialized {
		panic("TerminalDisplay not initialized")
	}

	newStringBuffer := ""

	// Clear the screen
	newStringBuffer += "\033[2J"
	newStringBuffer += "\033[H"
	newStringBuffer += "\n"

	// Print the buffer to the terminal
	for y := 0; y < HEIGHT; y++ {
		for x := 0; x < WIDTH; x++ {
			color := t.buffer[y][x]
			switch color {
			case 0:
				newStringBuffer += "  "
			case 1:
				newStringBuffer += "░░"
			case 2:
				newStringBuffer += "▒▒"
			case 3:
				newStringBuffer += "▓▓"
			}
		}
		newStringBuffer += "\n"
	}

	// Move cursor to top left
	// newStringBuffer += "\033[H"

	if t.stringBuffer != newStringBuffer {
		// Print the buffer to the terminal
		fmt.Print(newStringBuffer)
		t.stringBuffer = newStringBuffer
	}
}

func NewTerminalDisplay() *TerminalDisplay {
	t := &TerminalDisplay{initialized: true}

	for y := 0; y < HEIGHT; y++ {
		for x := 0; x < WIDTH; x++ {
			t.buffer[y][x] = 0
		}
	}

	return t
}
