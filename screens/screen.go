package screens

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type screen struct {
	style lipgloss.Style
	err   error
}

func (s *screen) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "esc":
			cmd = tea.Quit
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
	return cmd
}
