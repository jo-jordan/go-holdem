package main

import (
	"fmt"
	"log"

	"github.com/jo-jordan/go-holdem/network"
	"github.com/jo-jordan/go-holdem/scenes"
	"github.com/libp2p/go-libp2p/core/peer"

	golog "github.com/ipfs/go-log/v2"
)

func main() {
	golog.SetAllLoggers(golog.LevelDebug) // Change to INFO for extra info

	demo := &scenes.Demo{}
	go demo.New()

	p2p, err := network.MakeP2PWithHandlers(
		network.EventHandlers{
			OnStarted: func(addr string) {
				demo.SendLog("listening for connections")
				demo.SendLog(fmt.Sprintf("listener ready on %s", addr))
				demo.SendLog(fmt.Sprintf("Now run \"./go-holdem -l 3001 -d %s\" on a different terminal", addr))
			},
			OnMessageReceived: func(from peer.ID, payload []byte) {
				log.Printf("msg from %s: %s\n", from, payload)
			},
			OnSend: func(to peer.ID, payload []byte) {
				log.Printf("sent to %s: %s\n", to, payload)
			},
		},
	)
	if err != nil {
		log.Fatal("This is impossiable!")
	}

	p2p.StartHosting()
}
