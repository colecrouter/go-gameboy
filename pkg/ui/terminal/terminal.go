package terminal

import (
	"fmt"
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
	app   *Application
	frame string
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
		return m, tickCmd()
	case tea.KeyMsg:
		// Check both the key string and key type for ctrl+c.
		if msg.String() == "q" || msg.Type == tea.KeyCtrlC {
			m.app.gb.Stop()
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *terminalModel) View() string {
	return fmt.Sprintf("%s\nPress q or ctrl+c to exit.", m.frame)
}
