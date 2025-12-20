package screens

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jo-jordan/go-holdem/ui"
)

type Game struct {
	screen
	name string
}

type GameOption struct {
	Name  string
	Style lipgloss.Style
}

func NewGame(opt GameOption) *Game {
	game := &Game{
		name: opt.Name,
	}
	game.style = opt.Style

	return game
}

func (g *Game) Init() tea.Cmd {
	return nil
}

func (g *Game) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmd := g.screen.Update(msg)
	if cmd != nil {
		return g, cmd
	}
	return g, cmd
}

func (g *Game) View() string {
	return fmt.Sprintf("%s is creating a new game\n", g.name)
}

type JoinGame struct {
	screen
	ui.CursorMove
	name         string
	target       *ui.InputText
	joinButton   *ui.Button
	cancelButton *ui.Button
}

type JoinGameOption struct {
	Name  string
	Style lipgloss.Style
}

func NewJoinGame(opt JoinGameOption) *JoinGame {
	game := JoinGame{
		name: opt.Name,
		target: ui.NewInputText(ui.InputTextOption{
			Title: "Target: ",
			Focus: true,
			Actions: []*ui.ActionMap{
				tabToNext,
				enterToNext,
				shiftTabToPrev,
			},
		}),
	}

	game.joinButton = ui.NewButton(ui.ButtonOption{
		Value: "Join",
		Actions: []*ui.ActionMap{
			tabToNext,
			shiftTabToPrev,
			{
				Msg: tea.KeyEnter,
				Act: func() (tea.Model, tea.Cmd) {
					target := game.target.Value()
					if target == "" {
						return nil, nil
					}
					return NewGame(GameOption{Name: target, Style: game.style}), nil
				},
			},
		},
	})
	game.cancelButton = ui.NewButton(ui.ButtonOption{
		Value: "Cancel",
		Actions: []*ui.ActionMap{
			tabToNext,
			shiftTabToPrev,
			{
				Msg: tea.KeyEnter,
				Act: func() (tea.Model, tea.Cmd) {
					return NewStartSreen().WithStyle(&game.style), nil
				},
			},
		},
	})
	game.CursorMove = ui.NewCursorMove([]tea.Model{
		game.target,
		game.joinButton,
		game.cancelButton,
	})
	game.style = opt.Style
	return &game
}

func (g *JoinGame) Init() tea.Cmd {
	return textinput.Blink
}

func (g *JoinGame) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var model tea.Model
	cmd := g.screen.Update(msg)
	if cmd != nil {
		return g, cmd
	}

	model, cmd = g.CursorMove.Update(msg)
	if model == nil {
		model = g
	}
	return model, cmd
}

func (g *JoinGame) View() string {
	return g.style.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			fmt.Sprintf("User: %s", g.name),
			g.target.View(),
			lipgloss.JoinHorizontal(
				lipgloss.Center,
				g.joinButton.View(),
				g.cancelButton.View(),
			),
		),
	)

}
