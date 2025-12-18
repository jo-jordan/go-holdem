package cmd

const (
	Announcement uint8 = iota + 0
	Tick
	StartGame
	Deal
	Check
	Call
	Raise
	Fold
	EndGame
)

type GameCmd struct {
	Command uint8 `json:"command"`
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
