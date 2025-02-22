package terminal

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/colecrouter/gameboy-go/pkg/display/debug/reginfo"
	"github.com/colecrouter/gameboy-go/pkg/system"
	"github.com/colecrouter/gameboy-go/private/display"
	"github.com/colecrouter/gameboy-go/private/display/debug/logs"
	"github.com/colecrouter/gameboy-go/private/display/debug/tilemap"
	"github.com/colecrouter/gameboy-go/private/display/debug/tiles"
	"github.com/colecrouter/gameboy-go/private/display/monochrome"
	"github.com/colecrouter/gameboy-go/private/display/monochrome/lcd"
	"github.com/colecrouter/gameboy-go/private/memory/registers"
	"github.com/colecrouter/gameboy-go/private/ui/logger"
	"github.com/colecrouter/gameboy-go/private/ui/terminal/utils"
	"golang.org/x/term"
)

type Application struct {
	gb          *system.GameBoy
	menus       map[rune]display.Display
	mainDisplay display.Display
	openMenu    rune
	refresh     *time.Ticker
	lastOutput  string
}

// NewApplication creates a new terminal application.
func NewApplication(gb *system.GameBoy) *Application {
	app := &Application{gb: gb}
	app.menus = map[rune]display.Display{
		'v': tiles.NewTileDebug(gb.VRAM, &monochrome.Palette),
		'l': logs.NewLogMenu(),
		'm': tilemap.NewTilemapDebug(gb.VRAM, &monochrome.Palette),
		'r': reginfo.NewLogMenu(gb.IO),
	}
	app.mainDisplay = lcd.NewDisplay(gb.PPU)
	app.refresh = time.NewTicker(16 * time.Millisecond)

	return app
}

func (a *Application) Run() {
	// Set terminal to raw mode.
	oldState, _ := term.MakeRaw(int(os.Stdin.Fd()))
	// Defer terminal restoration in the main goroutine.
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	// Create a panic channel.
	panicChan := make(chan interface{}, 1)

	// Start the GameBoy runtime in a goroutine with panic recovery.
	go func() {
		defer func() {
			if r := recover(); r != nil {
				panicChan <- r
			}
		}()
		a.gb.Start(false)
	}()

	// Redirect stdout and stderr to the log menu.
	if logMenu, ok := a.menus['l'].(*logs.LogMenu); ok {
		logger.RedirectOutput(logMenu)
	}

	// Create channels for keyboard input and OS signals.
	// Changed channel type from rune to string.
	inputChan := make(chan string)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	// Launch goroutine to capture STDIN and handle escape sequences.
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			r, _, err := reader.ReadRune()
			if err != nil {
				continue
			}
			// Check for escape sequence (arrow keys).
			if r == '\x1b' {
				r2, _, err2 := reader.ReadRune()
				if err2 != nil || r2 != '[' {
					inputChan <- string(r)
					continue
				}
				r3, _, err3 := reader.ReadRune()
				if err3 != nil {
					inputChan <- string(r)
					continue
				}
				inputChan <- fmt.Sprintf("%c%c%c", r, r2, r3)
			} else {
				inputChan <- string(r)
			}
		}
	}()

	// Main event loop.
Loop:
	for {
		select {
		case <-a.refresh.C:
			a.render()

			// Reset pressed buttons.
			a.gb.Controller().ResetButtons()

		case key := <-inputChan:
			// Process menu bindings remain unchanged.
			if len(key) == 1 {
				if _, ok := a.menus[rune(key[0])]; ok {
					if a.openMenu == rune(key[0]) {
						a.openMenu = 0
						continue
					} else {
						a.openMenu = rune(key[0])
						continue
					}
				}
			}

			// Handle quit key.
			if key == "q" {
				break Loop
			}

			// Process controller bindings.
			joy := a.gb.Controller()
			var butt registers.Button

			switch key {
			case "\x1b[A":
				butt = registers.Button_Up
			case "\x1b[B":
				butt = registers.Button_Down
			case "\x1b[C":
				butt = registers.Button_Right
			case "\x1b[D":
				butt = registers.Button_Left
			case "+":
				butt = registers.Button_A
			case "-":
				butt = registers.Button_B
			case "*":
				butt = registers.Button_Start
			case "/":
				butt = registers.Button_Select
			default:
				continue
			}

			joy.SetButton(butt, true)
		case <-sigChan:
			break Loop

		case <-panicChan:
			// Optionally log p
			break Loop
		}
	}

	// Stop the GameBoy runtime.
	a.gb.Stop()

	// Clear the screen.
	fmt.Print("\033[H\033[2J")
}

func (a *Application) render() {
	a.gb.PPU.DisplayClock()

	for _, menu := range a.menus {
		menu.Clock()
	}

	clearScreen := "\033[H\033[2J"

	var screens [][]string
	screens = append(screens, utils.DrawBox(a.mainDisplay, &utils.BoxOptions{Border: utils.BorderSingle}))
	if a.openMenu != 0 && a.menus[a.openMenu] != nil {
		m := a.menus[a.openMenu]
		screens = append(screens, utils.DrawBox(m, &utils.BoxOptions{Border: utils.BorderDouble}))
	}

	var output string
	for _, screen := range screens {
		output += "\n\r"
		for _, line := range screen {
			output += line + "\n\r"
		}
	}

	if output != a.lastOutput {
		// Write directly to the original stdout.
		logger.OriginalOutput.Write([]byte(clearScreen))
		logger.OriginalOutput.Write([]byte(output))
		a.lastOutput = output
	}
}
