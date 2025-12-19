package screens

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jo-jordan/go-holdem/ui"
)

type StartScreen struct {
	screen
	column ui.Column
}

func NewStartSreen() StartScreen {
	column := ui.NewColumn(ui.ContainerOption{
		Elements: []ui.Element{
			ui.NewInputText(ui.InputTextOption{
				Title: "Enter your name: ",
				Focus: true,
			}),
			ui.NewRow(ui.ContainerOption{
				Elements: []ui.Element{
					ui.NewButton(ui.ButtonOption{
						Value:  "New Game",
						Action: NewGame,
					}),
					ui.NewButton(ui.ButtonOption{
						Value: "Join In",
					}),
				},
			}),
		},
	})
	style := lipgloss.NewStyle().
		Align(lipgloss.Center)
	return StartScreen{
		column: *column,
		screen: screen{
			style: style,
			err:   nil,
		},
	}
}

func (s StartScreen) Init() tea.Cmd {
	return textinput.Blink
}

func (s StartScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var model tea.Model
	cmd := s.screen.Update(msg)
	if cmd != nil {
		return s, cmd
	}
	model, cmd = s.column.Update(msg)
	if model == nil {
		model = s
	}
	return model, cmd
}

func (s StartScreen) View() string {
	return s.style.Render(s.column.View())
}
