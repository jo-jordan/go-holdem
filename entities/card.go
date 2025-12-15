package entities

import "fmt"

type CardType string

const (
	HEARTS   CardType = "♥"
	SPADES   CardType = "♠"
	DIAMONDS CardType = "♦"
	CLUBS    CardType = "♣"
)

var PointDisMap map[uint8]string = make(map[uint8]string, 4)

func init() {
	PointDisMap[1] = "A"
	PointDisMap[11] = "J"
	PointDisMap[12] = "Q"
	PointDisMap[13] = "K"
}

type Card struct {
	Point uint8
	Type  CardType
}

func (c *Card) String() string {
	d, ok := PointDisMap[c.Point]
	if !ok {
		d = fmt.Sprintf("%d", c.Point)
	}
	return fmt.Sprintf("%s%s", d, c.Type)
}
