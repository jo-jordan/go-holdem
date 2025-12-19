package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ContainerOption struct {
	Elements []Element
	Pos      lipgloss.Position
}

type container struct {
	elements   []Element
	pos        lipgloss.Position
	focusIndex int
	target     Element
	viewFormat func(lipgloss.Position, ...string) string
}

func newContainer(opt ContainerOption) *container {
	l := &container{
		elements:   opt.Elements,
		pos:        opt.Pos,
		focusIndex: -1,
	}

	l.initCursor()
	return l
}

func (c *container) initCursor() {
	for i, ele := range c.elements {
		if ele.Focused() {
			c.focusIndex = i
		}
		ele.SetTarget(c)
	}
}

func (c *container) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var model tea.Model
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case moveToPrevMsg, moveToNextMsg:
		cmd = c.moveCursor(msg)
	default:
		if c.Focused() {
			model, cmd = c.elements[c.focusIndex].Update(msg)
		}
	}

	return model, cmd
}

func (c *container) View() string {
	if c.viewFormat == nil {
		return ""
	}
	strs := make([]string, len(c.elements))
	for i, ele := range c.elements {
		strs[i] = ele.View()
	}
	return c.viewFormat(c.pos, strs...)
}

func (c *container) Focused() bool {
	return c.focusIndex != -1
}

func (c *container) SetTarget(t Element) {
	c.target = t
}

func (c *container) moveCursor(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case moveToNextMsg:
		cmd = c.moveToNext(msg.target)
	case moveToPrevMsg:
		cmd = c.moveToPrev(msg.target)
	}

	return cmd
}

func (c *container) moveToNext(target Element) tea.Cmd {
	var cmd tea.Cmd

	t := target
	switch target := target.(type) {
	case *Row:
		t = target.container
	case *Column:
		t = target.container
	}

	if t == c {
		c.focusIndex++
		if c.focusIndex == len(c.elements) {
			c.focusIndex = -1
			return func() tea.Msg {
				return moveToNextMsg{target: c.target}
			}
		}

		target = c.elements[c.focusIndex]
	} else {
		if t == nil {
			c.focusIndex = 0
			target = c.elements[c.focusIndex]
		} else {
			target = t
		}
	}
	_, cmd = c.elements[c.focusIndex].Update(moveToNextMsg{target: target})
	return cmd
}

func (c *container) moveToPrev(target Element) tea.Cmd {
	var cmd tea.Cmd

	t := target
	switch target := target.(type) {
	case *Row:
		t = target.container
	case *Column:
		t = target.container
	}

	if t == c {
		c.focusIndex--
		if c.focusIndex == -1 {
			return func() tea.Msg {
				return moveToPrevMsg{target: c.target}
			}
		} else if c.focusIndex < -1 {
			c.focusIndex = len(c.elements) - 1
		}

		target = c.elements[c.focusIndex]
	} else {
		if t == nil {
			c.focusIndex = len(c.elements) - 1
			target = c.elements[c.focusIndex]
		} else {
			target = t
		}
	}
	_, cmd = c.elements[c.focusIndex].Update(moveToPrevMsg{target: target})
	return cmd
}

type Row struct {
	*container
}

func NewRow(opt ContainerOption) *Row {
	r := new(Row)
	r.container = newContainer(opt)
	r.container.viewFormat = lipgloss.JoinHorizontal
	return r
}

func (r *Row) Init() tea.Cmd {
	return nil
}

func (r *Row) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return r.container.Update(msg)
}

type Column struct {
	*container
}

func NewColumn(opt ContainerOption) *Column {
	c := new(Column)
	c.container = newContainer(opt)
	c.container.viewFormat = lipgloss.JoinVertical
	return c
}

func (c *Column) Init() tea.Cmd {
	return nil
}

func (c *Column) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return c.container.Update(msg)
}
