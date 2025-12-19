package screens

import tea "github.com/charmbracelet/bubbletea"

type Game struct {
	screen
}

func newGame() *Game {
	return &Game{}
}

func NewGame() tea.Model {
	return newGame()
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
	return "Game Start"
}
