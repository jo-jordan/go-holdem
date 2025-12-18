package ui

import tea "github.com/charmbracelet/bubbletea"

type Element interface {
	tea.Model
	Focused() bool
}

type moveToPrevMsg struct{}

type moveToNextMsg struct{}
