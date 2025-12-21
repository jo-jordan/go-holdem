package ui

import tea "charm.land/bubbletea/v2"

type Elementer interface {
	Update(tea.Msg) (tea.Model, tea.Cmd)
	View() string
}
