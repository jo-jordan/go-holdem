package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type LayoutOption struct {
	Elements []Element
	Pos      lipgloss.Position
	IsRoot   bool
}

type Layout struct {
	elements   []Element
	pos        lipgloss.Position
	focusIndex int
	isRoot     bool
}

func newLayout(opt LayoutOption) *Layout {
	l := &Layout{
		elements:   opt.Elements,
		pos:        opt.Pos,
		focusIndex: -1,
		isRoot:     opt.IsRoot,
	}

	l.initCursor()
	return l
}

func (l *Layout) initCursor() {
	for i, ele := range l.elements {
		if ele.Focused() {
			l.focusIndex = i
		}
	}
}

func (l *Layout) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	// Handle move to next/prev into
	switch msg.(type) {
	case moveToNextMsg:
		l.focusIndex = 0
	case moveToPrevMsg:
		l.focusIndex = len(l.elements) - 1
	}

	if !l.focused() {
		return cmd
	}

	// Update focused element
	current := l.elements[l.focusIndex]
	_, cmd = current.Update(msg)

	if current.Focused() {
		return cmd
	}

	cmds := []tea.Cmd{cmd}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		cmds = append(cmds, l.moveCursor(msg))
	}
	return tea.Batch(cmds...)
}

func (l *Layout) moveCursor(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	switch msg.Type {
	case tea.KeyTab:
		cmd = l.moveToNext()
	case tea.KeyShiftTab:
		cmd = l.moveToPrev()
	}
	return cmd
}

func (l *Layout) moveToPrev() tea.Cmd {
	var cmd tea.Cmd
	l.focusIndex--
	if l.focusIndex == -1 && !l.isRoot {
		return cmd
	} else if l.focusIndex == -1 && l.isRoot {
		l.focusIndex = len(l.elements) - 1
	}
	_, cmd = l.elements[l.focusIndex].Update(moveToPrevMsg{})
	return cmd
}

func (l *Layout) moveToNext() tea.Cmd {
	var cmd tea.Cmd
	l.focusIndex++
	if l.focusIndex == len(l.elements) && !l.isRoot {
		l.focusIndex = -1
		return cmd
	} else if l.focusIndex == len(l.elements) && l.isRoot {
		l.focusIndex = 0
	}
	_, cmd = l.elements[l.focusIndex].Update(moveToNextMsg{})
	return cmd
}

func (l Layout) view() []string {
	strs := make([]string, len(l.elements))
	for i, ele := range l.elements {
		strs[i] = ele.View()
	}
	return strs
}

func (l *Layout) focused() bool {
	return l.focusIndex != -1
}

type Row struct {
	layout *Layout
}

func NewRow(opt LayoutOption) *Row {
	return &Row{
		layout: newLayout(opt),
	}
}

func (r *Row) Init() tea.Cmd {
	return nil
}

func (r *Row) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return r, r.layout.Update(msg)
}

func (r *Row) View() string {
	return lipgloss.JoinHorizontal(r.layout.pos, r.layout.view()...)
}

func (r *Row) Focused() bool {
	return r.layout.focused()
}

type Column struct {
	layout *Layout
}

func NewColumn(opt LayoutOption) *Column {
	return &Column{
		layout: newLayout(opt),
	}
}

func (c *Column) Init() tea.Cmd {
	return nil
}

func (c *Column) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return c, c.layout.Update(msg)
}

func (c *Column) View() string {
	return lipgloss.JoinVertical(c.layout.pos, c.layout.view()...)
}

func (c *Column) Focused() bool {
	return c.layout.focused()
}
