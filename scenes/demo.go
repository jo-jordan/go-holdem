package scenes

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Demo struct {
	program *tea.Program
}

type model struct {
	logs      []string
	ticks     int
	altscreen bool
	players   []string
	logViewport  viewport.Model
	playersViewport viewport.Model
}

type tickMsg time.Time
type logMsg struct {
	content string
}

const sceneWidth = 144
const playerViewportHeight = 48
const logViewportHeight = 12

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

var (
	playersViewportStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#000000ff")).
		Padding(20).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("228")).
		BorderBackground(lipgloss.Color("63")).
		PaddingRight(2)


	logViewportStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingRight(2)

	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render
)

func initialModel() model {
	logViewport := viewport.New(sceneWidth, logViewportHeight)
	logViewport.Style = logViewportStyle

	playerViewport := viewport.New(sceneWidth, playerViewportHeight)
	playerViewport.Style = playersViewportStyle

	logViewport.MouseWheelEnabled = true
	playerViewport.MouseWheelEnabled = true

	return model{
		logs:      []string{"Start Logging"},
		ticks:     0,
		altscreen: true,
		players:   []string{"Player1", "Player2", "Player3"},
		logViewport:  logViewport,
		playersViewport: playerViewport,
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

		var cmd tea.Cmd
		m.logs = m.appendLogs(msg.content)
		
		m.logViewport, cmd = m.logViewport.Update(m.logs)
		return m, cmd
	case tea.KeyMsg:
		{
			key := msg.String()
			switch key {
			case "q":
				tea.ClearScreen()
				os.Exit(1)
			}
			if len(key) > 0 {
				var cmd tea.Cmd
				m.logs = m.appendLogs(fmt.Sprintf("Key %s pressed.", key))
				m.logViewport, cmd = m.logViewport.Update(m.logs)
				return m, cmd
			}
		}
	}

	return m, nil
}

func (m model) appendLogs(content string) []string {
	timestamp := time.Now().Format("15:04:05")
	return append(m.logs, fmt.Sprintf("[%s] %s", timestamp, content))
}

func (m model) View() string {
	// Send the UI for rendering
	return m.playersInfoView() + "\n" + m.logAreaView() + m.helperView()
}

func (m model) playersInfoView() string {
	var s strings.Builder
	fmt.Fprintf(&s, "Players List\n\n(%d ticks)\n", m.ticks)
	for i, player := range m.players {
		fmt.Fprintf(&s, "ID: %d, Name: [%s]\n", i, player)
	}

	m.playersViewport.SetContent(s.String())
	m.playersViewport.GotoBottom()
	return m.playersViewport.View()
}

func (m model) logAreaView() string {
	var s strings.Builder
	for _, log := range m.logs {
		fmt.Fprintf(&s, "%s\n", log)
	}
	m.logViewport.SetContent(s.String())
	m.logViewport.GotoBottom()
	return m.logViewport.View()
}

func (m model) helperView() string {
	return helpStyle("\n  ↑/↓: Navigate • q: Quit\n")
}
