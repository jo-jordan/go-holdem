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
	contents *Contents
}

func NewMessage(style lipgloss.Style) *Message {
	width := style.GetWidth()
	// height := style.GetHeight()
	m := new(Message)
	m.style = lipgloss.NewStyle().
		Width(width).
		Height(MESSAGE_HEIGHT).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#FF5FAF"))
	m.box = NewViewPort(ViewPortOption{
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
	m.contents = NewContents(MAX_CONTENTS_LENGTH)

	return m
}

func (m *Message) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Key().Code {
		case tea.KeyEnter:
			m.send()
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

func (m *Message) send() {
	v := m.text.Value()
	if v == "" {
		return
	}
	m.contents.Add(v)
	m.text.text.SetValue("")
	m.box.vp.SetContentLines(m.contents.Contents())
	m.box.vp.GotoBottom()
}

type Contents struct {
	start    int
	end      int
	contents []string
}

func NewContents(maxLength int) *Contents {
	return &Contents{
		start:    0,
		end:      0,
		contents: make([]string, maxLength),
	}
}

func (c *Contents) Add(item string) {
	c.contents[c.end] = item
	c.end++

	if c.end == len(c.contents) {
		c.end = 0
	}
	if c.end == c.start {
		c.start++
		if c.start == len(c.contents) {
			c.start = 0
		}
	}
}

func (c *Contents) Contents() []string {
	var o []string
	if c.start == c.end {
		return o
	} else if c.start < c.end {
		o = make([]string, c.end-c.start)
		for i := range c.end - c.start {
			o[i] = c.contents[i]
		}
		return o
	} else {
		o = make([]string, len(c.contents))
		// 0 1 2 3 4 5 6 7 8 9
		l := len(c.contents) - c.end
		for i := range l {
			o[i] = c.contents[c.end+i]
		}

		for i := range c.end {
			o[i+l] = c.contents[i]
		}
		return o
	}
}
