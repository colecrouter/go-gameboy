package reginfo

import (
	"fmt"

	"github.com/colecrouter/gameboy-go/private/display"
	"github.com/colecrouter/gameboy-go/private/memory/io"
	"github.com/colecrouter/gameboy-go/private/ui/terminal/utils"
)

type LogMenu struct {
	config display.Config
	reg    *io.Registers
}

func NewLogMenu(r *io.Registers) *LogMenu {
	return &LogMenu{reg: r, config: display.Config{Width: 48 * utils.CHAR_WIDTH, Title: "Registers", Height: 12 * utils.CHAR_HEIGHT}}
}

func (l *LogMenu) Clock() {
}

// Text returns the logged messages.
func (l *LogMenu) Text() []string {
	var text []string

	text = append(text, fmt.Sprintf("Scroll X/Y: %d/%d", l.reg.ScrollX, l.reg.ScrollY))
	text = append(text, fmt.Sprintf("Window X/Y: %d/%d", l.reg.WindowX, l.reg.WindowY))
	text = append(text, fmt.Sprintf("Enabled BG/Window: %t", l.reg.LCDControl.EnableBackgroundAndWindow))
	text = append(text, fmt.Sprintf("Enabled Sprites: %t", l.reg.LCDControl.EnableSprites))
	text = append(text, fmt.Sprintf("2nd Sprite Map BG/Window: %t/%t", l.reg.LCDControl.BackgroundUseSecondTileMap, l.reg.LCDControl.WindowUseSecondTileMap))
	return text
}

func (l *LogMenu) Config() *display.Config {
	return &l.config
}
