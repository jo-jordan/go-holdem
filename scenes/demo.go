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
	program         *tea.Program
	handlers        DemoEventHandlers
	logs            []string
	ready           chan struct{}
	ticks           int
	altscreen       bool
	players         []string
	logViewport     viewport.Model
	playersViewport viewport.Model
}

type DemoEventHandlers struct {
	OnGameCmd func(cmd []byte)
}

type tickMsg time.Time
type logMsg struct {
	content string
}
type playersMsg struct {
	players []string
}

const sceneWidth = 144
const playerViewportHeight = 48
const logViewportHeight = 12

func NewDemo() *Demo {
	logViewport := viewport.New(sceneWidth, logViewportHeight)
	logViewport.Style = logViewportStyle

	playerViewport := viewport.New(sceneWidth, playerViewportHeight)
	playerViewport.Style = playersViewportStyle

	logViewport.MouseWheelEnabled = true
	playerViewport.MouseWheelEnabled = true

	d := &Demo{
		logs:            []string{"Start Logging"},
		ready:           make(chan struct{}),
		ticks:           0,
		altscreen:       true,
		players:         []string{"Player1", "Player2", "Player3"},
		logViewport:     logViewport,
		playersViewport: playerViewport,
	}

	return d
}

func (d *Demo) WaitUntilReady() {
	<-d.ready
}

func (d *Demo) RunWithEventHandlers(handlers DemoEventHandlers) {
	d.handlers = handlers
	go func() {
		d.program = tea.NewProgram(d, tea.WithAltScreen())
		close(d.ready)
		if _, err := d.program.Run(); err != nil {
			d.AppendLog(fmt.Sprintf("Alas, there's been an error: %v", err))
			os.Exit(1)
		}
	}()
}

func (d *Demo) AppendLog(message string) {
	if d.program == nil {
		return
	}
	d.program.Send(logMsg{content: message})
}

func (d *Demo) UpdatePlayers(players []string) {
	if d.program == nil {
		return
	}
	d.program.Send(playersMsg{players: players})
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

func tick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (d *Demo) Init() tea.Cmd {
	return tick()
}

func (d *Demo) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tickMsg:
		{
			d.ticks += 1
			d.handlers.OnGameCmd(fmt.Appendf(nil, "This is tick cmd: %d.\n", d.ticks))
			return d, tick()
		}
	case playersMsg:
		var cmd tea.Cmd
		d.players = msg.players
		d.playersViewport, cmd = d.playersViewport.Update(d.players)
		return d, cmd
	case logMsg:
		var cmd tea.Cmd
		d.logs = d.appendLogs(msg.content)

		d.logViewport, cmd = d.logViewport.Update(d.logs)
		return d, cmd
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
				d.logs = d.appendLogs(fmt.Sprintf("Key %s pressed.", key))
				d.logViewport, cmd = d.logViewport.Update(d.logs)
				return d, cmd
			}
		}
	}

	return d, nil
}

func (d *Demo) appendLogs(content string) []string {
	timestamp := time.Now().Format("15:04:05")
	return append(d.logs, fmt.Sprintf("[%s] %s", timestamp, content))
}

func (d *Demo) View() string {
	// Send the UI for rendering
	return d.playersInfoView() + "\n" + d.logAreaView() + d.helperView()
}

func (d *Demo) playersInfoView() string {
	var s strings.Builder
	fmt.Fprintf(&s, "Players List\n\n(%d ticks)\n", d.ticks)
	for i, player := range d.players {
		fmt.Fprintf(&s, "No.: %d, PeerID: [%s]\n", i, player)
	}

	d.playersViewport.SetContent(s.String())
	d.playersViewport.GotoBottom()
	return d.playersViewport.View()
}

func (d *Demo) logAreaView() string {
	var s strings.Builder
	for _, log := range d.logs {
		fmt.Fprintf(&s, "%s\n", log)
	}
	d.logViewport.SetContent(s.String())
	d.logViewport.GotoBottom()
	return d.logViewport.View()
}

func (d *Demo) helperView() string {
	return helpStyle("\n  ↑/↓: Navigate • q: Quit\n")
}
