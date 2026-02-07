package entities

type PlayerStatus uint8

const (
	// JOINING is that the server receives the join request before the user joined
	JOINING PlayerStatus = iota

	// JOINED is that the use is in the game, but is not ready
	JOINED

	// READY is the user is READY
	READY

	// SUSPEND is the user is in a suspend status, waiting for the dealer or other players
	SUSPEND

	// ACTION is the user's turn
	ACTION

	// FOLD is that the user folds in this round
	FOLD
)

type Player struct {
	Name string
	// Addr   netip.AddrPort
	Status PlayerStatus
}

func NewPlayer(name string) Player {
	return Player{
		Name: name,
		// Addr: ,
		Status: JOINING,
	}
}
