package ui

import tea "github.com/charmbracelet/bubbletea"

type Element interface {
	Update(tea.Msg) (tea.Model, tea.Cmd)
	View() string
	Focused() bool
	SetTarget(t Element)
}

type moveToPrevMsg struct {
	target Element
}

type moveToNextMsg struct {
	target Element
}
