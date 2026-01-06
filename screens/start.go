package screens

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/jo-jordan/go-holdem/ui"
)

type StartScreen struct {
	screen
	ui.CursorMove

	name         *ui.InputText
	createButton *ui.Button
	joinButton   *ui.Button
}

func NewStartSreen() *StartScreen {
	start := &StartScreen{
		name: ui.NewInputText(ui.InputTextOption{
			Title: "Input your name: ",
			Focus: true,
			Actions: []*ui.ActionMap{
				ui.TabToNext,
				ui.EnterToNext,
				ui.ShiftTabToPrev,
			},
		}),
	}
	start.createButton = ui.NewButton(ui.ButtonOption{
		Value: "Create Game",
		Actions: []*ui.ActionMap{
			ui.TabToNext,
			ui.ShiftTabToPrev,
			{
				Msg: "enter",
				Act: func() (tea.Model, tea.Cmd) {
					name := start.name.Value()
					if name == "" {
						return nil, nil
					}
					return NewGame(GameOption{Name: name, Style: start.style}), nil
				},
			},
		},
	})
	start.joinButton = ui.NewButton(ui.ButtonOption{
		Value: "Join in",
		Actions: []*ui.ActionMap{
			ui.TabToNext,
			ui.ShiftTabToPrev,
			{
				Msg: "enter",
				Act: func() (tea.Model, tea.Cmd) {
					name := start.name.Value()
					if name == "" {
						return nil, nil
					}
					return NewJoinGame(JoinGameOption{Name: name, Style: start.style}), nil
				},
			},
		},
	})
	start.CursorMove = ui.NewCursorMove(ui.CursorMoveOption{
		Models: []ui.Elementer{
			start.name,
			start.createButton,
			start.joinButton,
		},
	})
	return start
}

func (s *StartScreen) WithStyle(style *lipgloss.Style) *StartScreen {
	if style != nil {
		s.style = *style
	}
	return s
}

func (s *StartScreen) Init() tea.Cmd {
	return textinput.Blink
}

func (s *StartScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var model tea.Model
	cmd := s.screen.Update(msg)
	if cmd != nil {
		return s, cmd
	}

	model, cmd = s.CursorMove.Update(msg)
	if model == nil {
		model = s
	}
	return model, cmd
}

func (s *StartScreen) View() tea.View {
	v := tea.NewView(s.style.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			s.name.View(),
			lipgloss.JoinHorizontal(
				lipgloss.Center,
				s.createButton.View(),
				s.joinButton.View(),
			),
		),
	))
	v.AltScreen = true
	return v
}
