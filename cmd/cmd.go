package cmd

import "github.com/libp2p/go-libp2p/core/peer"

const (
	Announcement uint8 = iota + 0
	Tick
)

type GameCmd struct {
	Command uint8 `json:"command"`
}
type AnnouncementCmd struct {
	GameCmd
	Peer peer.AddrInfo `json:"peer"`
}

type TickCmd struct {
	GameCmd
	Tick string `json:"tick"`
}
