package network

import (
	"bufio"
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	mrand "math/rand"
	"sync"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"

	"github.com/multiformats/go-multiaddr"
	ma "github.com/multiformats/go-multiaddr"
)

type peerConn struct {
	rw *bufio.ReadWriter
	id peer.ID
}

type p2p struct {
	mu       sync.RWMutex
	peers    map[peer.ID]*peerConn
	Handlers EventHandlers
	IsHost   bool
}

func NewP2p() *p2p {
	return &p2p{
		peers:  make(map[peer.ID]*peerConn),
		IsHost: false,
	}
}

type EventHandlers struct {
	OnStarted         func(listenAddr string)
	OnMessageReceived func(from peer.ID, payload []byte)
	OnSent            func(to peer.ID, payload []byte)
	OnLog             func(content string)
	OnPeersUpdated    func(peers []peer.ID)
}

const protocolID = "/go-holdem/1.0.0"

func (p *p2p) SendData(to peer.ID, payload []byte) {
	p.mu.RLock()
	conn, ok := p.peers[to]
	p.mu.RUnlock()
	if !ok {
		if p.Handlers.OnLog != nil {
			p.Handlers.OnLog(fmt.Sprintf("no connection to %s", to))
		}
		return
	}

	fmt.Fprintf(conn.rw, "%s\n", payload)
	conn.rw.Flush()

	if p.Handlers.OnSent != nil {
		go p.Handlers.OnSent(to, payload)
	}
}

func (p *p2p) Broadcast(payload []byte) {
	p.mu.RLock()
	peers := make([]*peerConn, 0, len(p.peers))
	for _, conn := range p.peers {
		peers = append(peers, conn)
	}
	p.mu.RUnlock()

	for _, conn := range peers {
		fmt.Fprintf(conn.rw, "%s\n", payload)
		conn.rw.Flush()
		if p.Handlers.OnSent != nil {
			go p.Handlers.OnSent(conn.id, payload)
		}
	}
}

func (p *p2p) trackPeer(conn *peerConn) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peers[conn.id] = conn

	if p.Handlers.OnPeersUpdated != nil {
		go p.Handlers.OnPeersUpdated(keys(p.peers))
	}
}

func keys[K comparable, V any](m map[K]V) []K {
	out := make([]K, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

func (p *p2p) dropPeer(id peer.ID) {
	p.mu.Lock()
	delete(p.peers, id)
	p.mu.Unlock()
}

func (p *p2p) handleStream(s network.Stream) {
	conn := &peerConn{
		id: s.Conn().RemotePeer(),
		rw: bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s)),
	}
	p.trackPeer(conn)
	go p.readData(conn)
}

func (p *p2p) readData(conn *peerConn) {
	for {
		str, err := conn.rw.ReadString('\n')
		if err != nil {
			p.dropPeer(conn.id)
		}
		if str == "" {
			return
		}
		if str != "\n" {
			// Green console colour: 	\x1b[32m
			// Reset console colour: 	\x1b[0m
			p.Handlers.OnMessageReceived(conn.id, fmt.Appendf(nil, "\x1b[32m%s\x1b[0m ", str))
		}
	}
}

func (p *p2p) StartHosting() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Parse options from the command line
	listenF := flag.Int("l", 3000, "wait for incoming connections")
	targetF := flag.String("d", "", "target peer to dial")
	seedF := flag.Int64("seed", 0, "set random seed for id generation")
	flag.Parse()

	if *listenF == 0 {
		p.Handlers.OnLog("Please provide a port to bind on with -l")
	}

	// Make a host that listens on the given multiaddress
	h, err := p.makeHost(*listenF, *seedF)
	if err != nil {
		p.Handlers.OnLog(err.Error())
	}

	if *targetF == "" {
		p.startPeer(ctx, h, p.handleStream)
		p.IsHost = true
		if p.Handlers.OnStarted != nil {
			fullAddr := p.getHostAddress(h)
			go p.Handlers.OnStarted(fullAddr)
		}
	} else {
		conn, err := p.startPeerAndConnect(ctx, h, *targetF)
		if err != nil {
			p.Handlers.OnLog(err.Error())
			return
		}

		p.trackPeer(conn)
		go p.readData(conn)
	}

	// Run until canceled.
	<-ctx.Done()
}

func (p *p2p) makeHost(listenPort int, randseed int64) (host.Host, error) {
	var r io.Reader
	if randseed == 0 {
		r = rand.Reader
	} else {
		r = mrand.New(mrand.NewSource(randseed))
	}

	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		return nil, err
	}

	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", listenPort))

	opts := []libp2p.Option{
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(priv),
		libp2p.DisableRelay(),
	}

	return libp2p.New(opts...)
}

func (p *p2p) startPeer(_ context.Context, h host.Host, streamHandler network.StreamHandler) {
	h.SetStreamHandler(protocolID, streamHandler)

	// Let's get the actual TCP port from our listen multiaddr, in case we're using 0 (default; random available port).
	var port string
	for _, la := range h.Network().ListenAddresses() {
		if p, err := la.ValueForProtocol(multiaddr.P_TCP); err == nil {
			port = p
			break
		}
	}

	if port == "" {
		p.Handlers.OnLog("was not able to find actual local port")
		return
	}
}

func (p *p2p) startPeerAndConnect(_ context.Context, h host.Host, destination string) (*peerConn, error) {
	for _, la := range h.Addrs() {
		p.Handlers.OnLog(fmt.Sprintf(" - %v\n", la))
	}

	maddr, err := ma.NewMultiaddr(destination)
	if err != nil {
		p.Handlers.OnLog(err.Error())
		return nil, err
	}

	// Extract the peer ID from the multiaddr.
	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		p.Handlers.OnLog(err.Error())
		return nil, err
	}

	// We have a peer ID and a targetAddr, so we add it to the peerstore
	// so LibP2P knows how to contact it
	h.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)

	s, err := h.NewStream(context.Background(), info.ID, protocolID)
	if err != nil {
		p.Handlers.OnLog(err.Error())
		return nil, err
	}

	p.Handlers.OnLog("Established connection to destination")

	conn := &peerConn{
		id: info.ID,
		rw: bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s)),
	}

	return conn, err
}

func (p *p2p) getHostAddress(ha host.Host) string {
	// Build host multiaddress
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/p2p/%s", ha.ID()))

	// Now we can build a full multiaddress to reach this host
	// by encapsulating both addresses:
	addr := ha.Addrs()[0]
	return addr.Encapsulate(hostAddr).String()
}
