package entities

import (
	"sync"
)

const (
	MAX_SEAT = 10
)

type Game struct {
	mut sync.Mutex

	Round       uint
	SmallBlind  uint
	InitAccount uint
	Players     []*Player
}

func NewGame() *Game {
	players := make([]*Player, MAX_SEAT)
	for i := range MAX_SEAT {
		players[i] = nil
	}
	return &Game{
		Round:   0,
		Players: players,
	}
}

func (g *Game) AddPlayer(player *Player) error {
	g.mut.Lock()
	defer g.mut.Unlock()

	for i := range MAX_SEAT {
		if g.Players[i] == nil {
			g.Players[i] = player
			break
		}
	}
	return nil
}
