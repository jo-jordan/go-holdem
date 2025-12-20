package screens

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jo-jordan/go-holdem/ui"
)

type screen struct {
	style lipgloss.Style
	err   error
}

func (s *screen) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
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

var (
	tabToNext = &ui.ActionMap{
		Msg: tea.KeyTab,
		Act: ui.MoveToNext,
	}
	enterToNext = &ui.ActionMap{
		Msg: tea.KeyEnter,
		Act: ui.MoveToNext,
	}
	shiftTabToPrev = &ui.ActionMap{
		Msg: tea.KeyShiftTab,
		Act: ui.MoveToPrev,
	}
)
