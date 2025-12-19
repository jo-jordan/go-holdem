package screens

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type screen struct {
	style lipgloss.Style
	err   error
}

func (s *screen) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return tea.Quit
		}
	case tea.WindowSizeMsg:
		s.style = s.style.
			Width(msg.Width).
			Height(msg.Height).
			AlignVertical(lipgloss.Center).
			AlignHorizontal(lipgloss.Center)
	case error:
		s.err = msg
	}
	return nil
}
