package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type LayoutOption struct {
	Elements []Element
	Pos      lipgloss.Position
}

type Layout struct {
	elements   []Element
	pos        lipgloss.Position
	focusIndex int
}

func newLayout(opt LayoutOption) *Layout {
	l := &Layout{
		elements:   opt.Elements,
		pos:        opt.Pos,
		focusIndex: -1,
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
