package main

import (
	"encoding/json"
	"fmt"

	"github.com/jo-jordan/go-holdem/cmd"
	"github.com/jo-jordan/go-holdem/entities"
	"github.com/jo-jordan/go-holdem/internal"
	"github.com/jo-jordan/go-holdem/network"
	"github.com/jo-jordan/go-holdem/scenes"
	"github.com/libp2p/go-libp2p/core/peer"

	golog "github.com/ipfs/go-log/v2"
)

func main() {
	golog.SetAllLoggers(golog.LevelDebug) // Change to INFO for extra info

	p2p := network.NewP2p()
	ms := scenes.NewMainScene()
	ms.RunWithEventHandlers(
		scenes.MainSceneEventHandlers{
			OnTickCmd: func(cmd cmd.TickCmd) {
				if p2p.IsHost {
					p2p.Broadcast(cmd)
				}
			},
			OnSendChatMessage: func(chatCmd cmd.ChatCmd) {
				p2p.Broadcast(chatCmd)
			},
		},
	)
	ms.WaitUntilReady()

	p2p.Handlers = network.P2pEventHandlers{
		OnStarted: func(playerID peer.ID, addr string) {
			ms.AppendLog("listening for connections")
			ms.AppendLog(fmt.Sprintf("listener ready on %s", addr))
			ms.AppendLog(fmt.Sprintf("Now run \"./go-holdem -l 3001 -d %s\" on a different terminal", addr))
			ms.PlayerID = playerID.String()
		},
		OnCmdReceived: func(from peer.ID, command cmd.Command, payload []byte) {
			switch command {
			case cmd.Tick:
				var message cmd.TickCmd
				if err := json.Unmarshal(payload, &message); err != nil {
					ms.AppendLog(fmt.Sprintf("bad tick cmd from %s: %v", from, err))
					return
				}
			case cmd.Chat:
				var message cmd.ChatCmd
				if err := json.Unmarshal(payload, &message); err != nil {
					ms.AppendLog(fmt.Sprintf("bad chat cmd from %s: %v", from, err))
					return
				}
				ms.AppendChatMessage(entities.ChatMessage{
					SenderID: message.SenderID,
					Content:  message.Content,
				})
			default:
				ms.AppendLog(fmt.Sprintf("unknown command %d from %s", command, from))
			}

		},
		OnSent: func(to peer.ID, payload []byte) {
			ms.AppendLog(fmt.Sprintf("sent to %s: %s", to, payload))
		},
		OnLog: func(content string) {
			ms.AppendLog(content)
		},
		OnPeersUpdated: func(peers []peer.ID) {
			ms.UpdatePlayers(internal.Map(peers, func(id peer.ID) string {
				return id.String()
			}))
		},
	}

	p2p.StartHosting()
}
