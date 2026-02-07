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

type CursorMoveOption struct {
	Index  int
	Models []Elementer
}

func NewCursorMove(opt CursorMoveOption) CursorMove {
	return CursorMove{
		index:  opt.Index,
		models: opt.Models,
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

type IngoreQuitMsg struct{}

func IgnoreQuitCmd() tea.Msg {
	return IngoreQuitMsg{}
}

type ActionMap struct {
	Msg string
	Act func() (tea.Model, tea.Cmd)
}

var (
	TabToNext = &ActionMap{
		Msg: "tab",
		Act: MoveToNext,
	}
	EnterToNext = &ActionMap{
		Msg: "enter",
		Act: MoveToNext,
	}
	ShiftTabToPrev = &ActionMap{
		Msg: "shift+tab",
		Act: MoveToPrev,
	}
)
