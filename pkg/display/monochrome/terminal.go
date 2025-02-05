package monochrome

import (
	"fmt"

	"github.com/colecrouter/gameboy-go/pkg/display"
)

type TerminalDisplay struct {
	initialized  bool
	buffer       [display.HEIGHT][display.WIDTH]uint8
	stringBuffer string
}

func (t *TerminalDisplay) DrawScanline(row uint8, line []uint8) {
	for x := 0; x < display.WIDTH; x++ {
		color := line[x]
		t.buffer[row][x] = color
	}
}

func (t *TerminalDisplay) Clock() {
	if !t.initialized {
		panic("TerminalDisplay not initialized")
	}

	var newStringBuffer string

	// Clear the screen
	newStringBuffer += "\033[2J"
	newStringBuffer += "\033[H"
	newStringBuffer += "\n"

	// Print the buffer to the terminal
	for y := 0; y < display.HEIGHT; y++ {
		for x := 0; x < display.WIDTH; x++ {
			color := t.buffer[y][x]
			switch color {
			case 0:
				newStringBuffer += "▓▓"
			case 1:
				newStringBuffer += "▒▒"
			case 2:
				newStringBuffer += "░░"
			case 3:
				newStringBuffer += "  "
			}
		}
		newStringBuffer += "\n"
	}

	// Only update the terminal if the buffer has changed
	if t.stringBuffer != newStringBuffer {
		// Print the buffer to the terminal
		fmt.Print(newStringBuffer)
		t.stringBuffer = newStringBuffer
	}
}

func NewTerminalDisplay() *TerminalDisplay {
	t := &TerminalDisplay{initialized: true}

	for y := 0; y < display.HEIGHT; y++ {
		for x := 0; x < display.WIDTH; x++ {
			t.buffer[y][x] = 0
		}
	}

	return t
}
