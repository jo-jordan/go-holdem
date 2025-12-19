package cmd

import (
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

type Command uint8

const (
	Announcement Command = iota
	Tick
	StartGame
	Deal
	Check
	Call
	Raise
	Fold
	EndGame
	Chat
)

type Cmd interface{}

type GameCmd struct {
	Cmd
	Command Command `json:"command"`
}
type AnnouncementCmd struct {
	GameCmd
	Peer PeerDTO `json:"peer"`
}

type PeerDTO struct {
	ID    string   `json:"id"`
	Addrs []string `json:"addrs"`
}

func PeerDTOFromAddrInfo(info peer.AddrInfo) PeerDTO {
	addrs := make([]string, 0, len(info.Addrs))
	for _, addr := range info.Addrs {
		addrs = append(addrs, addr.String())
	}
	return PeerDTO{
		ID:    info.ID.String(),
		Addrs: addrs,
	}
}

func (dto PeerDTO) AddrInfo() (*peer.AddrInfo, error) {
	id, err := peer.Decode(dto.ID)
	if err != nil {
		return nil, err
	}

	addrs := make([]multiaddr.Multiaddr, 0, len(dto.Addrs))
	for _, s := range dto.Addrs {
		ma, err := multiaddr.NewMultiaddr(s)
		if err != nil {
			return nil, err
		}
		addrs = append(addrs, ma)
	}

	return &peer.AddrInfo{ID: id, Addrs: addrs}, nil
}

type TickCmd struct {
	GameCmd
	Tick string `json:"tick"`
}

type ChatCmd struct {
	GameCmd
	SenderID string `json:"sender_id"`
	Content  string `json:"content"`
}
