package main

import (
	"fmt"

	"github.com/jo-jordan/go-holdem/entities"
)

func main() {
	pack := entities.NewPack()
	pack.Shuffle()
	fmt.Println(pack)
}
