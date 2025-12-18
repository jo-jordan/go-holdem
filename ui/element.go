package ui

import tea "github.com/charmbracelet/bubbletea"

type Element interface {
	tea.Model
	Focused() bool
}

type MoveToNextMsg struct{}

func moveToNextCmd() tea.Msg {
	return MoveToNextMsg{}
}

type MoveToPreMsg struct{}

func moveToPrevCmd() tea.Msg {
	return MoveToPreMsg{}
}
