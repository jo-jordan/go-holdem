package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Button struct {
	value   string
	focus   bool
	style   lipgloss.Style
	actions []*ActionMap
}

type ButtonOption struct {
	Value   string
	Focus   bool
	Style   *lipgloss.Style
	Actions []*ActionMap
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
	b.actions = opt.Actions
	return b
}

func (b *Button) Init() tea.Cmd {
	return nil
}

func (b *Button) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var model tea.Model
	switch msg := msg.(type) {
	case focusMsg:
		b.focus = true
	case blurMsg:
		b.focus = false
	case tea.KeyMsg:
		for _, m := range b.actions {
			if msg.Type == m.Msg {
				model, cmd = m.Act()
				break
			}
		}
	}
	return model, cmd
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
