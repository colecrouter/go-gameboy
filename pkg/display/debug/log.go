package debug

import (
	"strings"

	"github.com/colecrouter/gameboy-go/pkg/display"
	"github.com/colecrouter/gameboy-go/pkg/ui/terminal/utils"
)

type LogMenu struct {
	logs   []string
	config display.Config
}

func NewLogMenu() *LogMenu {
	return &LogMenu{logs: make([]string, 0), config: display.Config{Width: 32 * utils.CHAR_WIDTH, Title: "Logs", Height: 10 * utils.CHAR_HEIGHT}}
}

func (l *LogMenu) Clock() {
}

// Text returns the logged messages.
func (l *LogMenu) Text() []string {
	// Limit line length to 32 characters.
	trimmed := make([]string, len(l.logs))
	for i, line := range l.logs {
		if len(line) > 32 {
			trimmed[i] = line[:32] + "..."
		} else {
			trimmed[i] = line
		}
		// Add reset color code to the end of each line.
		trimmed[i] += "\033[0m"
	}
	return l.logs
}

// Optionally, add a method to append log messages.
func (l *LogMenu) AddLog(log string) {
	lines := strings.Split(log, "\n")

	l.logs = append(l.logs, lines...)
	// Optionally limit the length of logs.
	if len(l.logs) > 10 {
		l.logs = l.logs[len(l.logs)-10:]
	}
}

func (l *LogMenu) Config() *display.Config {
	return &l.config
}
