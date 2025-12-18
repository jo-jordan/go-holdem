package entities

type Game struct {
	Round   uint
	Dealer  *Dealer
	Players []*Player
}
