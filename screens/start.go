package screens

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jo-jordan/go-holdem/ui"
)

type StartScreen struct {
	column ui.Column
	style  lipgloss.Style
	err    error
}

func NewStartSreen() StartScreen {
	column := ui.NewColumn(ui.LayoutOption{
		IsRoot: true,
		Elements: []ui.Element{
			ui.NewInputText(ui.InputTextOption{
				Title: "Enter your name: ",
				Focus: true,
			}),
			ui.NewRow(ui.LayoutOption{
				Elements: []ui.Element{
					ui.NewButton(ui.ButtonOption{
						Value: "New Game",
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
		style:  style,
		err:    nil,
	}
}

func (s StartScreen) Init() tea.Cmd {
	return textinput.Blink
}

func (s StartScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return s, tea.Quit
		}
	case tea.WindowSizeMsg:
		s.style = s.style.
			Width(msg.Width).
			Height(msg.Height).
			AlignVertical(lipgloss.Center).
			AlignHorizontal(lipgloss.Center)
	case error:
		s.err = msg
		return s, nil
	}
	var cmd tea.Cmd
	_, cmd = s.column.Update(msg)
	return s, cmd
}

func (s StartScreen) View() string {
	return s.style.Render(s.column.View())
}
