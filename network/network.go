package network

func MakeP2PWithHandlers(handlers EventHandlers) (*p2p, error) {
	p2p := p2p{}
	p2p.handlers = handlers
	return &p2p, nil
}
