package events

import (
	"github.com/jo-jordan/go-holdem/cmd"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

type NetworkEvent interface{ isNetworkEvent() }

type NetworkReady struct {
	IsHost bool
	Host   host.Host
}

type PeersUpdated struct {
	Peers []peer.ID
}

type NetworkError struct {
	Err error
}

type NetworkLog struct {
	Msg string
}

type CmdReceived struct {
	From    peer.ID
	Header  cmd.Command
	Payload []byte
}

// markers
func (NetworkReady) isNetworkEvent() {}
func (PeersUpdated) isNetworkEvent() {}
func (NetworkLog) isNetworkEvent()   {}
func (NetworkError) isNetworkEvent() {}
func (CmdReceived) isNetworkEvent()  {}
