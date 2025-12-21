package main

import (
	"log"

	tea "charm.land/bubbletea/v2"
	"github.com/jo-jordan/go-holdem/screens"
)

func main() {
	p := tea.NewProgram(screens.NewStartSreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
