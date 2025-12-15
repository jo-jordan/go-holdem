package entities

import "math/rand/v2"

type Pack struct {
	Cards []*Card
}

func NewPack() *Pack {
	cardTypes := []CardType{HEARTS, SPADES, DIAMONDS, CLUBS}
	var items uint8 = 13
	cards := make([]*Card, len(cardTypes)*int(items))
	var c int
	for _, i := range cardTypes {
		for j := range items {
			cards[c] = &Card{Point: j + 1, Type: i}
			c++
		}
	}
	return &Pack{
		Cards: cards,
	}
}

func (p *Pack) Shuffle() {
	rand.Shuffle(len(p.Cards), func(i int, j int) {
		p.Cards[i], p.Cards[j] = p.Cards[j], p.Cards[i]
	})
}
