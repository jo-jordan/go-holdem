package network

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/jo-jordan/go-holdem/cmd"
	"github.com/jo-jordan/go-holdem/events"
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

type P2p struct {
	out    chan<- events.NetworkEvent
	peers  sync.Map // key: peer.ID, value: *peerConn
	host   host.Host
	IsHost bool
}

func NewP2p(out chan events.NetworkEvent) *P2p {
	return &P2p{
		IsHost: false,
		out:    out,
	}
}

const (
	protocolID  = "/go-holdem/1.0.0"
	defaultPort = 3000
)

func writeFrame(w *bufio.Writer, payload []byte) error {
	if err := binary.Write(w, binary.BigEndian, uint32(len(payload))); err != nil {
		return err
	}
	if _, err := w.Write(payload); err != nil {
		return err
	}
	return w.Flush()
}

func (p *P2p) SendData(to peer.ID, payload []byte) {
	value, ok := p.peers.Load(to)
	if !ok {
		p.emitErr(fmt.Errorf("no connection to %s", to))
		return
	}

	conn := value.(*peerConn)
	if err := writeFrame(conn.rw.Writer, payload); err != nil {
		p.emitErr(fmt.Errorf("send to %s failed: %v", to, err))
		return
	}
}

func (p *P2p) emitErr(err error) {
	evt := events.NetworkError{
		Err: err,
	}
	p.emit(evt)
}

func (p *P2p) Broadcast(cmd cmd.Cmd) {
	payload, _ := json.Marshal(cmd)
	p.peers.Range(func(_, value any) bool {
		conn := value.(*peerConn)
		if err := writeFrame(conn.rw.Writer, payload); err != nil {
			p.emitErr(fmt.Errorf("broadcast to %s failed: %v", conn.id, err))
		}
		return true
	})
}

func (p *P2p) BroadcastExclude(cmd cmd.Cmd, excludes []peer.ID) {
	payload, _ := json.Marshal(cmd)
	excludeMap := make(map[peer.ID]struct{})
	for _, id := range excludes {
		excludeMap[id] = struct{}{}
	}

	p.peers.Range(func(_, value any) bool {
		conn := value.(*peerConn)
		if _, excluded := excludeMap[conn.id]; excluded {
			return true // skip excluded peer
		}
		if err := writeFrame(conn.rw.Writer, payload); err != nil {
			p.emitErr(fmt.Errorf("broadcast to %s failed: %v", conn.id, err))
		}
		return true
	})
}

func (p *P2p) trackPeer(conn *peerConn) {
	p.peers.Store(conn.id, conn)

	evt := events.PeersUpdated{
		Peers: p.snapshotPeerIDs(),
	}
	p.emit(evt)
}

func (p *P2p) dropPeer(id peer.ID) {
	p.peers.Delete(id)

	evt := events.PeersUpdated{
		Peers: p.snapshotPeerIDs(),
	}
	p.emit(evt)
}

func (p *P2p) snapshotPeerIDs() []peer.ID {
	ids := make([]peer.ID, 0)
	p.peers.Range(func(key, _ any) bool {
		if id, ok := key.(peer.ID); ok {
			ids = append(ids, id)
		}
		return true
	})
	return ids
}

func (p *P2p) handleStream(s network.Stream) {
	remoteID := s.Conn().RemotePeer()
	remoteMA := s.Conn().RemoteMultiaddr()
	if remoteMA != nil {
		p.host.Peerstore().AddAddr(remoteID, remoteMA, peerstore.TempAddrTTL)
	}

	conn := &peerConn{
		id: remoteID,
		rw: bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s)),
	}

	info := peer.AddrInfo{
		ID:    conn.id,
		Addrs: p.host.Peerstore().Addrs(remoteID),
	}
	p.trackPeer(conn)
	go p.announceNewPeer(info)
	go p.readData(conn)
}

func (p *P2p) announceNewPeer(newPeer peer.AddrInfo) {
	payload := cmd.AnnouncementCmd{
		GameCmd: cmd.GameCmd{
			Command: cmd.Announcement,
		},
		Peer: cmd.PeerDTOFromAddrInfo(newPeer),
	}

	p.BroadcastExclude(payload, []peer.ID{newPeer.ID})
}

func (p *P2p) ensurePeerStream(ctx context.Context, id peer.ID) error {
	if _, ok := p.peers.Load(id); ok {
		return nil // already have an active stream
	}

	s, err := p.host.NewStream(ctx, id, protocolID)
	if err != nil {
		return err
	}

	conn := &peerConn{
		id: id,
		rw: bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s)),
	}

	p.trackPeer(conn)
	go p.readData(conn)
	return nil
}

// handleMessage processes an incoming message from a peer.
func (p *P2p) handleMessage(from peer.ID, data []byte) {
	var header cmd.GameCmd
	if err := json.Unmarshal(data, &header); err != nil {
		p.emitErr(fmt.Errorf("invalid payload from %s: %v", from, err))
		return
	}

	switch header.Command {
	case cmd.Announcement:
		var message cmd.AnnouncementCmd
		if err := json.Unmarshal(data, &message); err != nil {
			p.emitErr(fmt.Errorf("bad announcement from %s: %v", from, err))
			return
		}

		info, err := message.Peer.AddrInfo()
		if err != nil {
			p.emitErr(fmt.Errorf("bad peer info in announcement from %s: %v", from, err))
			return
		}

		if p.host == nil {
			p.emitErr(fmt.Errorf("cannot connect to announced peer: host not initialized"))
			return
		}

		p.host.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.TempAddrTTL)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := p.host.Connect(ctx, *info); err != nil {
			p.emitErr(fmt.Errorf("connect to %s failed: %v", message.Peer.ID, err))
		}

		if err := p.ensurePeerStream(context.Background(), info.ID); err != nil {
			p.emitErr(fmt.Errorf("open stream to %s failed: %v", message.Peer.ID, err))
		}
	default:
		evt := events.CmdReceived{
			From:    from,
			Header:  header.Command,
			Payload: data,
		}
		p.emit(evt)
	}
}

func (p *P2p) readData(conn *peerConn) {
	for {
		payload, err := readFrame(conn.rw.Reader)
		if err != nil {
			if err != io.EOF {
				p.emitErr(fmt.Errorf("read from %s failed: %v", conn.id, err))
			}
			p.dropPeer(conn.id)
			return
		}

		if len(payload) == 0 {
			continue // skip empty frames
		}

		p.handleMessage(conn.id, payload)
	}
}

func readFrame(r *bufio.Reader) ([]byte, error) {
	var n uint32
	if err := binary.Read(r, binary.BigEndian, &n); err != nil {
		return nil, err
	}
	buf := make([]byte, n)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, err
	}
	return buf, nil
}

func (p *P2p) StartHosting(dest string, port int) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	h, err := p.makeHost(port)
	if err != nil {

		p.emitErr(fmt.Errorf("failed to create libp2p host: %v", err))

		return
	}

	p.host = h
	p.startPeer(ctx, h, p.handleStream)

	if dest == "" || len(dest) == 0 {
		p.IsHost = true
	} else {
		p.IsHost = false

		conn, err := p.startPeerAndConnect(ctx, h, dest)
		if err != nil {
			p.emitErr(fmt.Errorf("failed to connect to destination: %v", err))
			return
		}
		p.trackPeer(conn)
		go p.readData(conn)
	}

	evt := events.NetworkReady{
		IsHost: p.IsHost,
		Host:   p.host,
	}
	p.emit(evt)
}

func (p *P2p) emit(event events.NetworkEvent) {
	if p.out == nil {
		return
	}

	select {
	case p.out <- event:
		// sent
	default:
		p.emitErr(fmt.Errorf("network lifecycle channel is full; dropping Event"))
	}
}

func (p *P2p) makeHost(port int) (host.Host, error) {
	r := rand.Reader

	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		return nil, err
	}

	if port == 0 {
		port = defaultPort
	}
	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port))

	opts := []libp2p.Option{
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(priv),
		libp2p.DisableRelay(),
	}

	return libp2p.New(opts...)
}

func (p *P2p) startPeer(_ context.Context, h host.Host, streamHandler network.StreamHandler) {
	h.SetStreamHandler(protocolID, streamHandler)

	var port string
	for _, la := range h.Network().ListenAddresses() {
		if p, err := la.ValueForProtocol(multiaddr.P_TCP); err == nil {
			port = p
			break
		}
	}

	if port == "" {
		p.emitErr(fmt.Errorf("failed to determine listening port"))
		return
	}
}

func (p *P2p) startPeerAndConnect(_ context.Context, h host.Host, destination string) (*peerConn, error) {
	for _, la := range h.Addrs() {
		p.emitErr(fmt.Errorf(" - %v", la))
	}

	maddr, err := ma.NewMultiaddr(destination)
	if err != nil {
		p.emitErr(fmt.Errorf("invalid destination multiaddress: %v", err))
		return nil, err
	}

	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		p.emitErr(fmt.Errorf("failed to parse peer address info: %v", err))
		return nil, err
	}

	h.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)

	s, err := h.NewStream(context.Background(), info.ID, protocolID)
	if err != nil {
		p.emitErr(fmt.Errorf("failed to create stream to peer: %v", err))
		return nil, err
	}

	conn := &peerConn{
		id: info.ID,
		rw: bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s)),
	}

	return conn, nil
}
