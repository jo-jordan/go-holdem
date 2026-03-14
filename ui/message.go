package ui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	Cmd "github.com/jo-jordan/go-holdem/cmd"
)

const (
	TABLE_HEIGHT         = 25 // This value should be same as the height of table
	MESSAGE_BORDER_WIDTH = 2
	MAX_CONTENTS_LENGTH  = 1000
)

type Message struct {
	box      *ViewPort
	style    lipgloss.Style
	text     *InputText
	contents []string
}

func NewMessage(style lipgloss.Style) *Message {
	width := style.GetWidth()
	m := new(Message)
	m.style = lipgloss.NewStyle().
		Width(width).
		Height(style.GetHeight() - TABLE_HEIGHT).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(NORMAL_COLOR))
	m.box = NewViewPort(ViewPortOption{
		SoftWrap: true,
		Style: lipgloss.NewStyle().
			Width(width-MESSAGE_BORDER_WIDTH).
			Height(style.GetHeight()-TABLE_HEIGHT-3).
			Border(lipgloss.ThickBorder(), false, false, true, false).
			BorderBottomForeground(lipgloss.Color(NORMAL_COLOR)),
		Actions: []*ActionMap{
			TabToNext,
			ShiftTabToPrev,
		},
	})
	m.text = NewInputText(InputTextOption{
		Title: ">",
		Focus: false,
	})
	m.contents = make([]string, 0)

	return m
}

func (m *Message) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Key().Code {
		case tea.KeyEnter:
			m.send(m.text.Value())
		case tea.KeyUp:
			m.box.vp.ScrollUp(1)
		case tea.KeyDown:
			m.box.vp.ScrollDown(1)
		case tea.KeyPgUp:
			m.box.vp.PageUp()
		case tea.KeyPgDown:
			m.box.vp.PageDown()
		case tea.KeyEsc:
			m.style = m.style.BorderForeground(lipgloss.Color(NORMAL_COLOR))
			cmd = func() tea.Msg {
				return Cmd.TableFocusMsg{}
			}
		default:
			_, cmd = m.text.Update(msg)
		}
	case tea.WindowSizeMsg:
		width := msg.Width
		height := msg.Height
		m.style = m.style.Width(width).Height(height - TABLE_HEIGHT - 3)
		m.box.SetStyle(m.box.vp.Style.Width(width - 2))
		_, cmd = m.box.Update(msg)
	case Cmd.MsgFocusMsg:
		m.style = m.style.BorderForeground(lipgloss.Color(FOCUS_COLOR))
		_, cmd = m.text.Update(Cmd.FocusMsg{})
	default:
		_, cmd = m.text.Update(msg)
	}
	return nil, cmd
}

func (m Message) View() string {
	return m.style.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			m.box.View(),
			m.text.View(),
		),
	)
}

func (m *Message) send(content string) {
	if content == "" {
		return
	}

	m.contents = append(m.contents, content)
	if len(m.contents) > MAX_CONTENTS_LENGTH {
		m.contents = m.contents[1:]
	}
	// FROM here
	m.text.text.SetValue("")
	m.box.vp.SetContentLines(m.contents)
	m.box.vp.GotoBottom()
}

func (m *Message) Height() int {
	return m.style.GetHeight()
}
