package ui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const (
	MESSAGE_HEIGHT       = 7
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
		Height(MESSAGE_HEIGHT).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#FF5FAF"))
	m.box = NewViewPort(ViewPortOption{
		SoftWrap: true,
		Style: lipgloss.NewStyle().
			Width(width-MESSAGE_BORDER_WIDTH).
			Height(MESSAGE_HEIGHT-MESSAGE_BORDER_WIDTH-1).
			Border(lipgloss.ThickBorder(), false, false, true, false).
			BorderBottomForeground(lipgloss.Color("#C0C0C0")),
		Actions: []*ActionMap{
			TabToNext,
			ShiftTabToPrev,
		},
	})
	m.text = NewInputText(InputTextOption{
		Title: ">",
		Focus: true,
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
		default:
			_, cmd = m.text.Update(msg)
		}
	case tea.WindowSizeMsg:
		width := msg.Width
		m.style = m.style.Width(width)
		m.box.SetStyle(m.box.vp.Style.Width(width - 2))
		_, cmd = m.box.Update(msg)
	default:
		_, cmd = m.text.Update(msg)
	}
	return nil, cmd
}

func (m Message) View() string {
	return m.style.Render(lipgloss.JoinVertical(
		lipgloss.Left,
		m.box.View(),
		m.text.View(),
	))
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
