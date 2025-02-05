package terminal

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/colecrouter/gameboy-go/pkg/renderer"
	"github.com/colecrouter/gameboy-go/pkg/system"
)

type Application struct {
	p  *tea.Program
	gb *system.GameBoy
}

// NewApplication creates a new terminal application.
func NewApplication(gb *system.GameBoy) *Application {
	return &Application{gb: gb}
}

// Run starts the fullscreen Bubble Tea interface.
func (app *Application) Run() error {
	// Create and run bubbletea program with our model.
	m := &terminalModel{app: app}
	app.p = tea.NewProgram(m, tea.WithAltScreen())
	return app.p.Start()
}

// terminalModel implements tea.Model.
type terminalModel struct {
	app       *Application
	frame     string
	showDebug bool   // new: whether the debug viewport is shown
	debugInfo string // new: placeholder text for debug info
}

type tickMsg time.Time

func (m *terminalModel) Init() tea.Cmd {
	// Start GameBoy in background.
	go m.app.gb.Start()
	// Kick off periodic ticks.
	return tickCmd()
}

func tickCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m *terminalModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tickMsg:
		// Update frame from GameBoy display.
		m.frame = renderer.RenderANSI(m.app.gb.Display.Image())
		// update debugInfo placeholder (can be replaced later)
		m.debugInfo = "Debug info: [placeholder]"
		return m, tickCmd()
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			m.app.gb.Stop()
			return m, tea.Quit
		case "d": // toggle debug viewport on key "d"
			m.showDebug = !m.showDebug
		}
		// Also handle ctrl+c.
		if msg.Type == tea.KeyCtrlC {
			m.app.gb.Stop()
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *terminalModel) View() string {
	// If debug viewport is hidden, show only the main display.
	if !m.showDebug {
		return m.frame + "\n\nPress 'd' to show debug info.\nPress q or ctrl+c to exit."
	}

	// If debug viewport is enabled, split main and debug views side-by-side.
	mainLines := strings.Split(m.frame, "\n")
	debugStr := m.debugInfo + "\nPress 'd' to hide debug info."
	debugLines := strings.Split(debugStr, "\n")

	// Determine maximum lines and combine side-by-side.
	maxLines := len(mainLines)
	if len(debugLines) > maxLines {
		maxLines = len(debugLines)
	}

	var combined []string
	for i := 0; i < maxLines; i++ {
		var left, right string
		if i < len(mainLines) {
			left = mainLines[i]
		}
		if i < len(debugLines) {
			right = debugLines[i]
		}
		combined = append(combined, fmt.Sprintf("%-80s   %s", left, right))
	}
	return strings.Join(combined, "\n") + "\nPress q or ctrl+c to exit."
}
