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
	"github.com/colecrouter/gameboy-go/pkg/renderer"
	"github.com/colecrouter/gameboy-go/pkg/system"
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

		case <-sigChan:
			break Loop
		}
	}

	a.gb.Stop()
}

func (a *Application) render() {
	for _, menu := range a.menus {
		menu.Clock()
	}

	var buffer string

	// Clear the screen
	buffer += "\033[H\033[2J"

	// Render the main screen
	img := a.gb.Display.Image()
	buffer += renderer.RenderSixel(img)

	if a.openMenu != 0 {
		menu := a.menus[a.openMenu]
		img := menu.Image()
		buffer += renderer.RenderSixel(img)
	}

	if buffer != a.lastOutput {
		fmt.Print(buffer)
		a.lastOutput = buffer
	}
}
