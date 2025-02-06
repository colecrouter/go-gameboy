package terminal

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/colecrouter/gameboy-go/pkg/display"
	"github.com/colecrouter/gameboy-go/pkg/display/debug"
	"github.com/colecrouter/gameboy-go/pkg/display/monochrome"
	"github.com/colecrouter/gameboy-go/pkg/system"
	"github.com/colecrouter/gameboy-go/pkg/ui/terminal/utils"
	"golang.org/x/term"
)

type Application struct {
	gb         *system.GameBoy
	menus      map[rune]display.Display
	openMenu   rune
	refresh    *time.Ticker
	lastOutput string
}

// NewApplication creates a new terminal application.
func NewApplication(gb *system.GameBoy) *Application {
	app := &Application{gb: gb}
	app.menus = map[rune]display.Display{
		'v': debug.NewTileDebug(gb.VRAM, &monochrome.Palette),
	}
	app.refresh = time.NewTicker(16 * time.Millisecond)

	return app
}

func (a *Application) Run() {
	// Start the GameBoy runtime.
	go a.gb.Start()

	// Set terminal to raw mode.
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

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
	for _, menu := range a.menus {
		menu.Clock()
	}

	clearScreen := "\033[H\033[2J"

	var screens [][]string
	if a.gb.Display != nil {
		screens = append(screens, utils.DrawBox(a.gb.Display.Image(), &utils.BoxOptions{Border: utils.BorderSingle, Title: "GameBoy Display"}))
	}
	if a.openMenu != 0 && a.menus[a.openMenu] != nil {
		screens = append(screens, utils.DrawBox(a.menus[a.openMenu].Image(), &utils.BoxOptions{Border: utils.BorderDouble, Title: "Tile Map"}))
	}

	var output string
	for _, screen := range screens {
		output += "\n\r"
		for _, line := range screen {
			output += line + "\n\r"
		}
	}

	if output != a.lastOutput {
		fmt.Print(clearScreen)
		fmt.Print(output)
		a.lastOutput = output
	}
}
