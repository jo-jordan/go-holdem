package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jo-jordan/go-holdem/screens"
)

func main() {
	p := tea.NewProgram(screens.NewStartSreen(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
