package core

import (
	"encoding/json"
	"fmt"

	"github.com/jo-jordan/go-holdem/cmd"
	"github.com/jo-jordan/go-holdem/entities"
	"github.com/jo-jordan/go-holdem/events"
	"github.com/jo-jordan/go-holdem/internal"
	"github.com/jo-jordan/go-holdem/network"
	"github.com/jo-jordan/go-holdem/scenes"
	"github.com/libp2p/go-libp2p/core/peer"
)

type MainController struct {
	netToScene chan events.NetworkEvent
	sceneToNet chan events.SceneEvent
	mainScene  *scenes.MainScene
	p2p        *network.P2p
}

func NewController() *MainController {
	netToScene := make(chan events.NetworkEvent, 64)
	sceneToNet := make(chan events.SceneEvent, 64)
	return &MainController{
		netToScene: netToScene,
		sceneToNet: sceneToNet,
		mainScene:  scenes.NewMainScene(sceneToNet),
		p2p:        network.NewP2p(netToScene),
	}
}

func (ctrl *MainController) Run(dest string, port int) {
	mainScene := ctrl.mainScene
	p2p := ctrl.p2p
	mainScene.RunWithEventHandlers(
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
	mainScene.WaitUntilReady()

	p2p.StartHosting(dest, port)

	for {
		select {
		case c := <-ctrl.netToScene:
			ctrl.handleNetworkEvents(c)
		case c := <-ctrl.sceneToNet:
			ctrl.handleScenesEvents(c)
		}
	}
}

func (ctrl *MainController) handleNetworkEvents(evt events.NetworkEvent) {
	switch e := evt.(type) {
	case events.NetworkReady:
		addr := network.GetHostAddress(e.Host)
		ctrl.mainScene.AppendLog("listening for connections")
		ctrl.mainScene.AppendLog(fmt.Sprintf("listener ready on %s", addr))
		ctrl.mainScene.AppendLog(fmt.Sprintf("Now run \"./go-holdem -l 3001 -d %s\" on a different terminal", addr))
		ctrl.mainScene.PlayerID = e.Host.ID().String()
	case events.PeersUpdated:
		ctrl.mainScene.UpdatePlayers(internal.Map(e.Peers, func(id peer.ID) string {
			return id.String()
		}))
	case events.NetworkError:
		ctrl.mainScene.AppendErrorLog(e.Err)
	case events.CmdReceived:
		switch e.Header {
		case cmd.Tick:
			var message cmd.TickCmd
			if err := json.Unmarshal(e.Payload, &message); err != nil {
				ctrl.mainScene.AppendErrorLog(fmt.Errorf("bad tick cmd from %s: %v", e.From, err))
				return
			}
		case cmd.Chat:
			var message cmd.ChatCmd
			if err := json.Unmarshal(e.Payload, &message); err != nil {
				ctrl.mainScene.AppendErrorLog(fmt.Errorf("bad chat cmd from %s: %v", e.From, err))
				return
			}
			ctrl.mainScene.AppendChatMessage(entities.ChatMessage{
				SenderID: message.SenderID,
				Content:  message.Content,
			})
		default:
			ctrl.mainScene.AppendErrorLog(fmt.Errorf("unknown command %d from %s", e.Header, e.From))
		}
	default:
		ctrl.mainScene.AppendErrorLog(fmt.Errorf("unknown network event"))
	}

}

func (ctrl *MainController) handleScenesEvents(cmd cmd.Cmd) {}
