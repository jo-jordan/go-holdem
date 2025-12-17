package scenes

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Demo struct {
	program *tea.Program
}

type model struct {
	log       string
	ticks     int
	altscreen bool
	players   []string
}

type tickMsg time.Time
type logMsg struct {
	content string
}

func (d *Demo) New() {
	d.program = tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := d.program.Run(); err != nil {
		d.SendLog(fmt.Sprintf("Alas, there's been an error: %v", err))
		os.Exit(1)
	}
}

func (d *Demo) SendLog(message string) {
	if d.program != nil {
		d.program.Send(logMsg{content: message})
	}
}

func initialModel() model {
	return model{
		log:       "",
		ticks:     0,
		altscreen: true,
		players:   []string{"Player1", "Player2", "Player3"},
	}
}

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Init() tea.Cmd {
	return tick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tickMsg:
		{
			m.ticks += 1
			return m, tick()
		}
	case logMsg:
		// Append new log message with timestamp
		timestamp := time.Now().Format("15:04:05")
		m.log += fmt.Sprintf("[%s] %s\n", timestamp, msg.content)
		return m, nil
	case tea.KeyMsg:
		{
			key := msg.String()
			if len(key) > 0 {
				m.log += fmt.Sprintf("Key %s pressed.\n", key)
			}
			switch key {
			case "q":
				tea.ClearScreen()
				os.Exit(1)
			}
		}
	}

	return m, nil
}

var style = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#7D56F4")).
	Padding(20).
	Width(144).
	BorderForeground(lipgloss.Color("228")).
	BorderBackground(lipgloss.Color("63")).
	BorderTop(true).
	BorderLeft(true).
	Align(lipgloss.Center)

func (m model) View() string {
	// The header
	s := fmt.Sprintf("Players List?\n\n(%d ticks)\n", m.ticks)

	// Iterate over our choices
	for i, player := range m.players {

		// Render the row
		s += fmt.Sprintf("ID: %d, Name: [%s]\n", i, player)
	}

	s += "\n--- Logs ---\n" + m.log

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return style.Render(s)
}
