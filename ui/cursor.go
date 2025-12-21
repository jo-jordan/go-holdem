package ui

import (
	tea "charm.land/bubbletea/v2"
)

type MoveToPrevMsg struct{}

type MoveToNextMsg struct{}

type focusMsg struct{}

type blurMsg struct{}

type CursorMove struct {
	index  int
	models []Elementer
}

func NewCursorMove(models []Elementer) CursorMove {
	return CursorMove{
		models: models,
	}
}

func (c *CursorMove) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	currentIndex := c.index
	var model tea.Model
	var cmd tea.Cmd
	switch msg.(type) {
	case MoveToNextMsg:
		c.index++
		c.index = c.index % len(c.models)
	case MoveToPrevMsg:
		c.index--
		if c.index < 0 {
			c.index = len(c.models) - 1
		}
	default:
		model, cmd = c.models[c.index].Update(msg)
		return model, cmd
	}

	_, cmd = c.models[currentIndex].Update(blurMsg{})
	_, cmd = c.models[c.index].Update(focusMsg{})
	return nil, cmd
}

func MoveToNext() (tea.Model, tea.Cmd) {
	return nil, func() tea.Msg {
		return MoveToNextMsg{}
	}
}

func MoveToPrev() (tea.Model, tea.Cmd) {
	return nil, func() tea.Msg {
		return MoveToPrevMsg{}
	}
}

type ActionMap struct {
	Msg string
	Act func() (tea.Model, tea.Cmd)
}
