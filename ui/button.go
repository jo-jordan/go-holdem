package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Button struct {
	value string
	focus bool
	style lipgloss.Style
}

type ButtonOption struct {
	Value string
	Focus bool
	Style *lipgloss.Style
}

func NewButton(opt ButtonOption) *Button {
	b := new(Button)
	b.value = opt.Value
	b.focus = opt.Focus
	if opt.Style == nil {
		b.style = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder())
	} else {
		b.style = *opt.Style
	}
	return b
}

func (b *Button) Init() tea.Cmd {
	return nil
}

func (b *Button) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case moveToPrevMsg, moveToNextMsg:
		b.focus = true
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab, tea.KeyShiftTab:
			b.focus = false
		}
	}
	return b, cmd
}

func (b Button) View() string {
	style := b.style
	if b.focus {
		style = b.style.
			Background(lipgloss.Color("#FF5FAF")).
			Foreground(lipgloss.Color("#FFFFFF"))
	}
	return style.Render(b.value)
}

func (b *Button) Focused() bool {
	return b.focus
}
