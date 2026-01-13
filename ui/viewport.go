package ui

import (
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type ViewPort struct {
	vp      viewport.Model
	actions []*ActionMap
}

type ViewPortOption struct {
	Style    lipgloss.Style
	Focus    bool
	Actions  []*ActionMap
	SoftWrap bool
}

func NewViewPort(opt ViewPortOption) *ViewPort {
	vp := viewport.New(
		viewport.WithWidth(opt.Style.GetWidth()),
		viewport.WithHeight(opt.Style.GetHeight()),
	)
	vp.Style = opt.Style
	vp.SoftWrap = opt.SoftWrap
	return &ViewPort{
		vp:      vp,
		actions: opt.Actions,
	}
}

func (v *ViewPort) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var model tea.Model
	var cmd tea.Cmd
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		for _, m := range v.actions {
			if msg.String() == m.Msg {
				model, cmd = m.Act()
				cmds = append(cmds, cmd)
				break
			}
		}
	case tea.WindowSizeMsg:
		v.vp.SetWidth(msg.Width)
		v.vp.SetHeight(msg.Height)
	}

	return model, cmd
}

func (v ViewPort) View() string {
	return v.vp.View()
}

func (v *ViewPort) SetStyle(style lipgloss.Style) {
	v.vp.Style = style
}
