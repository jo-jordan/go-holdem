package entities

type Dealer struct {
	Pack *Pack
}

func (d *Dealer) StartNewRound() {
	d.Pack = NewPack()
	d.Pack.Shuffle()
}
