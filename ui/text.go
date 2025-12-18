package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type InputText struct {
	text  textinput.Model
	style lipgloss.Style
}

type InputTextOption struct {
	Title       string
	InitText    string
	TextWidth   int
	PlaceHolder string
	Style       *lipgloss.Style
	Focus       bool
}

func NewInputText(opt InputTextOption) *InputText {
	text := textinput.New()
	text.Prompt = opt.Title
	text.SetValue(opt.InitText)
	text.Placeholder = opt.PlaceHolder
	if opt.TextWidth != 0 {
		text.Width = opt.TextWidth
		text.CharLimit = opt.TextWidth
	} else {
		text.Width = 24
		text.CharLimit = 24
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
		text:  text,
		style: style,
	}
}

func (i *InputText) Init() tea.Cmd {
	return textinput.Blink
}

func (i *InputText) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyTab:
			i.text.Blur()
			cmd = moveToNextCmd
		case tea.KeyShiftTab:
			i.text.Blur()
			cmd = moveToPrevCmd
		default:
			i.text, cmd = i.text.Update(msg)
		}
	case MoveToPreMsg, MoveToNextMsg:
		cmd = i.text.Focus()
	case cursor.BlinkMsg:
		i.text, cmd = i.text.Update(msg)
	}
	return i, cmd
}

func (i *InputText) View() string {
	return i.style.Render(fmt.Sprintf("%s\n", i.text.View()))
}

func (i *InputText) Focused() bool {
	return i.text.Focused()
}
