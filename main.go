package main

import (
	"fmt"

	"github.com/jo-jordan/go-holdem/internal"
	"github.com/jo-jordan/go-holdem/network"
	"github.com/jo-jordan/go-holdem/scenes"
	"github.com/libp2p/go-libp2p/core/peer"

	golog "github.com/ipfs/go-log/v2"
)

func main() {
	golog.SetAllLoggers(golog.LevelDebug) // Change to INFO for extra info

	p2p := network.NewP2p()
	demo := scenes.NewDemo()
	demo.RunWithEventHandlers(
		scenes.DemoEventHandlers{
			OnGameCmd: func(cmd []byte) {
				if p2p.IsHost {
					p2p.Broadcast(cmd)
				}

			},
		},
	)
	demo.WaitUntilReady()

	p2p.Handlers = network.P2pEventHandlers{
		OnStarted: func(addr string) {
			demo.AppendLog("listening for connections")
			demo.AppendLog(fmt.Sprintf("listener ready on %s", addr))
			demo.AppendLog(fmt.Sprintf("Now run \"./go-holdem -l 3001 -d %s\" on a different terminal", addr))
		},
		OnMessageReceived: func(from peer.ID, payload []byte) {
			//demo.AppendLog(fmt.Sprintf("msg from %s: %s", from, payload))
		},
		OnSent: func(to peer.ID, payload []byte) {
			demo.AppendLog(fmt.Sprintf("sent to %s: %s", to, payload))
		},
		OnLog: func(content string) {
			demo.AppendLog(content)
		},
		OnPeersUpdated: func(peers []peer.ID) {
			demo.UpdatePlayers(internal.Map(peers, func(id peer.ID) string {
				return id.String()
			}))
		},
	}

	p2p.StartHosting()
}
