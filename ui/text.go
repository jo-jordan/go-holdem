package ui

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type InputText struct {
	text    textinput.Model
	style   lipgloss.Style
	actions []*ActionMap
}

type InputTextOption struct {
	Title       string
	InitText    string
	TextWidth   int
	PlaceHolder string
	Style       *lipgloss.Style
	Focus       bool
	Actions     []*ActionMap
}

func NewInputText(opt InputTextOption) *InputText {
	text := textinput.New()
	text.Prompt = opt.Title
	text.SetValue(opt.InitText)
	text.Placeholder = opt.PlaceHolder
	if opt.TextWidth != 0 {
		text.SetWidth(opt.TextWidth)
		text.CharLimit = opt.TextWidth
	}

	var style lipgloss.Style
	if opt.Style != nil {
		style = *opt.Style
	} else {
		style = lipgloss.NewStyle()
	}

	if opt.Focus {
		text.Focus()
	}

	return &InputText{
		text:    text,
		style:   style,
		actions: opt.Actions,
	}
}

func (i *InputText) Init() tea.Cmd {
	return textinput.Blink
}

func (i *InputText) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	var cmd tea.Cmd
	var model tea.Model
	switch msg := msg.(type) {
	case focusMsg:
		cmd = i.text.Focus()
		cmds = append(cmds, cmd)
	case blurMsg:
		i.text.Blur()
	case tea.KeyPressMsg:
		for _, m := range i.actions {
			if msg.String() == m.Msg {
				model, cmd = m.Act()
				cmds = append(cmds, cmd)
				break
			}
		}
	}
	i.text, cmd = i.text.Update(msg)
	cmds = append(cmds, cmd)
	return model, tea.Batch(cmds...)
}

func (i InputText) View() string {
	style := i.style
	if i.text.Focused() {
		style = style.BorderForeground(lipgloss.Color("#FF5FAF"))
	}
	return style.Render(i.text.View())
}

func (i InputText) Value() string {
	return i.text.Value()
}
