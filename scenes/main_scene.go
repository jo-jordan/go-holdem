package scenes

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jo-jordan/go-holdem/cmd"
	"github.com/jo-jordan/go-holdem/entities"
	"github.com/jo-jordan/go-holdem/events"
	"github.com/muesli/reflow/wordwrap"
)

type MainScene struct {
	out             chan events.SceneEvent
	program         *tea.Program
	logs            []string
	ticks           int
	altscreen       bool
	players         []string
	logViewport     viewport.Model
	playersViewport viewport.Model
	textarea        textarea.Model
	chatMsgViewport viewport.Model
	chatMessages    []string
	PlayerID        string
}

type tickMsg time.Time
type logMsg struct {
	isError bool
	content string
}
type playersMsg struct {
	players []string
}

const sceneWidth = 120
const playerViewportHeight = 20
const logViewportHeight = 8
const chatViewportHeight = 8
const textareaHeight = 1
const gap = "\n"

func NewMainScene(out chan events.SceneEvent) *MainScene {
	logViewport := viewport.New(sceneWidth, logViewportHeight)
	logViewport.Style = logViewportStyle
	logViewport.MouseWheelEnabled = true

	playerViewport := viewport.New(sceneWidth, playerViewportHeight)
	playerViewport.Style = playersViewportStyle
	playerViewport.MouseWheelEnabled = true

	chatMsgViewport := viewport.New(sceneWidth, chatViewportHeight)
	chatMsgViewport.Style = logViewportStyle
	chatMsgViewport.MouseWheelEnabled = false

	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()
	ta.Prompt = " > "
	ta.CharLimit = sceneWidth
	ta.SetWidth(sceneWidth)
	ta.SetHeight(textareaHeight)
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false

	d := &MainScene{
		out:             out,
		logs:            []string{"Start Logging"},
		ticks:           0,
		altscreen:       true,
		players:         []string{},
		logViewport:     logViewport,
		playersViewport: playerViewport,
		chatMsgViewport: chatMsgViewport,
		textarea:        ta,
	}

	return d
}

func (d *MainScene) Run() {
	go func() {
		d.program = tea.NewProgram(d, tea.WithAltScreen())
		d.out <- events.SceneReady{}
		if _, err := d.program.Run(); err != nil {
			d.AppendLog(fmt.Sprintf("Alas, there's been an error: %v", err))
			os.Exit(1)
		}
	}()
}

func (d *MainScene) AppendLog(message string) {
	if d.program == nil {
		return
	}
	d.program.Send(logMsg{content: message, isError: false})
}

func (d *MainScene) AppendErrorLog(err error) {
	if d.program == nil {
		return
	}
	d.program.Send(logMsg{content: err.Error(), isError: true})
}

func (d *MainScene) UpdatePlayers(players []string) {
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
				Padding(0).
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("228")).
				BorderBackground(lipgloss.Color("63")).
				PaddingRight(2)

	logViewportStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("62"))

	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render

	senderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Render
)

func tick() tea.Cmd {
	return tea.Batch(
		tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
			return tickMsg(t)
		}),
		textarea.Blink,
	)
}

func (d *MainScene) Init() tea.Cmd {
	return tick()
}

func (d *MainScene) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tickMsg:
		{
			d.ticks += 1
			//gameCmd := cmd.TickCmd{
			//	GameCmd: cmd.GameCmd{
			//		Command: cmd.Tick,
			//	},
			//	Tick: fmt.Sprintf("This is tick cmd: %d.", d.ticks),
			//}
			//cmdData, err := json.Marshal(gameCmd)
			//if err == nil && d.handlers.OnTickCmd != nil {
			//	d.handlers.OnTickCmd(cmdData)
			//}
			return d, tick()
		}
	case playersMsg:
		var c tea.Cmd
		d.players = msg.players
		d.playersViewport, c = d.playersViewport.Update(d.players)
		return d, c
	case logMsg:
		var c tea.Cmd
		d.logs = d.appendLogs(msg.content, msg.isError)

		d.logViewport, c = d.logViewport.Update(d.logs)
		return d, c
	case tea.KeyMsg:
		var (
			tiCmd tea.Cmd
			vpCmd tea.Cmd
		)
		d.textarea, tiCmd = d.textarea.Update(msg)
		d.chatMsgViewport, vpCmd = d.chatMsgViewport.Update(msg)
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			os.Exit(1)
		case tea.KeyEnter:
			cm, _ := strings.CutSuffix(d.textarea.Value(), "\n")
			trimmed := strings.Trim(cm, "\n")
			if trimmed == "" {
				d.textarea.Reset()
				return d, nil
			}
			d.chatMessages = append(d.chatMessages, senderStyle("You: ")+cm)
			d.chatMsgViewport.SetContent(lipgloss.NewStyle().Width(d.chatMsgViewport.Width).Render(strings.Join(d.chatMessages, "\n")))
			d.textarea.Reset()
			d.chatMsgViewport.GotoBottom()
			evt := events.SceneChatMessage{
				ChatCmd: cmd.ChatCmd{
					GameCmd: cmd.GameCmd{
						Command: cmd.Chat,
					},
					SenderID: d.PlayerID,
					Content:  cm,
				},
			}
			d.emit(evt)
		}
		return d, tea.Batch(tiCmd, vpCmd)
	}

	return d, nil
}

func (d *MainScene) emit(evt events.SceneEvent) {
	d.out <- evt
}

func (d *MainScene) appendLogs(content string, isErr bool) []string {
	timestamp := time.Now().Format("15:04:05")
	log := fmt.Sprintf("[%s] %s", timestamp, content)
	if isErr {
		// light red for errors
		log = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(log)
	}
	return append(d.logs, log)
}

func (d *MainScene) View() string {
	return fmt.Sprintf(
		"%s%s%s%s%s%s%s%s",
		d.playersInfoView(),
		gap,
		d.chatMsgViewport.View(),
		gap,
		d.logAreaView(),
		d.helperView(),
		gap,
		d.textarea.View(),
	)
}

func (d *MainScene) playersInfoView() string {
	var s strings.Builder
	fmt.Fprintf(&s, "Players List\n\n(%d ticks)\n", d.ticks)
	for i, player := range d.players {
		fmt.Fprintf(&s, "No.: %d, PeerID: [%s]\n", i, player)
	}

	d.playersViewport.SetContent(s.String())
	d.playersViewport.GotoBottom()
	return d.playersViewport.View()
}

func (d *MainScene) logAreaView() string {
	var s strings.Builder
	for _, log := range d.logs {
		fmt.Fprintf(&s, "%s", wordwrap.String(fmt.Sprintf("%s\n", log), sceneWidth))
	}
	d.logViewport.SetContent(s.String())
	d.logViewport.GotoBottom()
	return d.logViewport.View()
}

func (d *MainScene) helperView() string {
	return helpStyle("\n esc: Quit\n")
}

func (d *MainScene) AppendChatMessage(message entities.ChatMessage) {
	d.chatMessages = append(d.chatMessages, senderStyle(fmt.Sprintf("%s: ", message.SenderID))+message.Content)
	d.chatMsgViewport.SetContent(lipgloss.NewStyle().Width(d.chatMsgViewport.Width).Render(strings.Join(d.chatMessages, "\n")))
	d.chatMsgViewport.GotoBottom()
}
