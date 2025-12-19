package main

import (
	"flag"

	"github.com/jo-jordan/go-holdem/core"
)

func main() {
	destF := flag.String("d", "", "target peer to dial")
	portF := flag.Int("l", 3000, "wait for incoming connections")

	flag.Parse()

	ctrl := core.NewController()
	ctrl.Run(*destF, *portF)
}
