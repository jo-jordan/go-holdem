package cmd

type Command uint8

const (
	Announcement Command = iota + 0
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

type TickCmd struct {
	GameCmd
	Tick string `json:"tick"`
}

type ChatCmd struct {
	GameCmd
	SenderID string `json:"sender_id"`
	Content  string `json:"content"`
}
