package ui

import (
	tea "charm.land/bubbletea/v2"
	Cmd "github.com/jo-jordan/go-holdem/cmd"
)

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
	case Cmd.MoveToNextMsg:
		c.index++
		c.index = c.index % len(c.models)
	case Cmd.MoveToPrevMsg:
		c.index--
		if c.index < 0 {
			c.index = len(c.models) - 1
		}
	default:
		model, cmd = c.models[c.index].Update(msg)
		return model, cmd
	}

	_, cmd = c.models[currentIndex].Update(Cmd.BlurMsg{})
	_, cmd = c.models[c.index].Update(Cmd.FocusMsg{})
	return nil, cmd
}

func MoveToNext() (tea.Model, tea.Cmd) {
	return nil, func() tea.Msg {
		return Cmd.MoveToNextMsg{}
	}
}

func MoveToPrev() (tea.Model, tea.Cmd) {
	return nil, func() tea.Msg {
		return Cmd.MoveToPrevMsg{}
	}
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
