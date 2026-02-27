package ui

import (
	"strconv"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	P "github.com/jo-jordan/go-holdem/entities"
)

var (
	activeStyle = lipgloss.
			NewStyle().
			Border(
			lipgloss.NormalBorder(),
		).
		Height(3).
		Width(20).
		BorderForeground(lipgloss.Color("#FFFFF"))
	inactiveStyle = activeStyle.
			BorderForeground(lipgloss.Color("#CCCCCC"))
)

type Seat struct {
	num    int
	player *P.Player

	style lipgloss.Style
}

type SeatOpt struct {
	Num    int
	Player *P.Player
}

func NewSeat(opt SeatOpt) *Seat {
	return &Seat{
		num:    opt.Num,
		player: opt.Player,
		style:  inactiveStyle,
	}
}

func (s *Seat) Init() tea.Msg {
	return nil
}

func (s *Seat) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil
}

func (s *Seat) View() string {
	var name string
	var chips string
	var card1 string
	var card2 string
	if s.player != nil {
		name = s.player.Name
		chips = strconv.Itoa(s.player.Account)
		if s.player.Card1 != nil {
			card1 = s.player.Card1.String()
		}
		if s.player.Card2 != nil {
			card2 = s.player.Card2.String()
		}
	}

	return s.style.Render(
		lipgloss.JoinVertical(
			lipgloss.Top,
			name,
			chips,
			lipgloss.JoinHorizontal(lipgloss.Left, card1, card2),
		),
	)
}
