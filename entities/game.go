package entities

import (
	"errors"
	"sync"
)

type Game struct {
	mut sync.Mutex

	Round   uint
	Players []*Player
	Seat    Seat
}

func NewGame() *Game {
	return &Game{
		Round:   0,
		Players: make([]*Player, 0),
		Seat:    newSeat(),
	}
}

func (g *Game) AddPlayer(player *Player) error {
	g.mut.Lock()
	defer g.mut.Unlock()

	if len(g.Players) >= SEAT_COUNT {
		return errors.New("Full")
	}

	g.Players = append(g.Players, player)
	return nil
}

const (
	SEAT_COUNT = 10
)

type Seat map[int]*Player

func newSeat() Seat {
	seat := make(Seat, SEAT_COUNT)
	for i := range 10 {
		seat[i+1] = nil
	}
	return seat
}
