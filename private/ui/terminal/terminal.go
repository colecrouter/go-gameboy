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
	"github.com/colecrouter/gameboy-go/private/ui/logger"
	"github.com/colecrouter/gameboy-go/private/ui/terminal/utils"
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
	// Start the GameBoy runtime.
	go a.gb.Start(true)

	// Set terminal to raw mode.
	// oldState, _ := term.MakeRaw(int(os.Stdin.Fd()))
	// defer func() {
	// 	recover()
	// 	if oldState != nil {
	// 		term.Restore(int(os.Stdin.Fd()), oldState)
	// 	}
	// }()

	// --- NEW: Redirect stdout and stderr to the log menu ---
	// Cast the log menu from the menus map.
	if logMenu, ok := a.menus['l'].(*logs.LogMenu); ok {
		logger.RedirectOutput(logMenu)
	}

	// Create channels for keyboard input and OS signals.
	inputChan := make(chan rune)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	// Launch goroutine to capture STDIN.
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			r, _, err := reader.ReadRune()
			if err != nil {
				continue
			}
			inputChan <- r
		}
	}()

	// Main event loop.
Loop:
	for {
		select {
		case <-a.refresh.C:
			a.render()
		case r := <-inputChan:
			_, ok := a.menus[r]
			if ok {
				if a.openMenu == r {
					a.openMenu = 0
					continue
				} else {
					a.openMenu = r
				}
			}
			if r == 'q' {
				break Loop
			}
		case <-sigChan:
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
