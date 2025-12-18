package network

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	mrand "math/rand"
	"sync"
	"time"

	"github.com/jo-jordan/go-holdem/cmd"
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
	peers    sync.Map // key: peer.ID, value: *peerConn
	host     host.Host
	Handlers P2pEventHandlers
	IsHost   bool
}

func NewP2p() *P2p {
	return &P2p{
		IsHost: false,
	}
}

type P2pEventHandlers struct {
	OnStarted         func(listenAddr string)
	OnMessageReceived func(from peer.ID, payload []byte)
	OnSent            func(to peer.ID, payload []byte)
	OnLog             func(content string)
	OnPeersUpdated    func(peers []peer.ID)
}

const protocolID = "/go-holdem/1.0.0"

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
		if p.Handlers.OnLog != nil {
			p.Handlers.OnLog(fmt.Sprintf("no connection to %s", to))
		}
		return
	}

	conn := value.(*peerConn)
	if err := writeFrame(conn.rw.Writer, payload); err != nil && p.Handlers.OnLog != nil {
		p.Handlers.OnLog(fmt.Sprintf("send to %s failed: %v", to, err))
		return
	}

	if p.Handlers.OnSent != nil {
		go p.Handlers.OnSent(to, payload)
	}
}

func (p *P2p) Broadcast(payload []byte) {
	p.peers.Range(func(_, value any) bool {
		conn := value.(*peerConn)
		if err := writeFrame(conn.rw.Writer, payload); err != nil && p.Handlers.OnLog != nil {
			p.Handlers.OnLog(fmt.Sprintf("broadcast to %s failed: %v", conn.id, err))
		} else if p.Handlers.OnSent != nil {
			go p.Handlers.OnSent(conn.id, payload)
		}
		return true
	})
}

func (p *P2p) trackPeer(conn *peerConn) {
	p.peers.Store(conn.id, conn)

	if p.Handlers.OnPeersUpdated != nil {
		go p.Handlers.OnPeersUpdated(p.snapshotPeerIDs())
	}
}

func (p *P2p) dropPeer(id peer.ID) {
	p.peers.Delete(id)

	if p.Handlers.OnPeersUpdated != nil {
		go p.Handlers.OnPeersUpdated(p.snapshotPeerIDs())
	}
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
	payload, _ := json.Marshal(cmd.AnnouncementCmd{
		GameCmd: cmd.GameCmd{
			Command: cmd.Announcement,
		},
		Peer: newPeer,
	})

	p.Broadcast(payload)
}

// handleMessage processes an incoming message from a peer. TODO need to move this to outside
func (p *P2p) handleMessage(from peer.ID, data []byte) {
	var header cmd.GameCmd
	if err := json.Unmarshal(data, &header); err != nil {
		if p.Handlers.OnLog != nil {
			p.Handlers.OnLog(fmt.Sprintf("invalid payload from %s: %v", from, err))
		}
		return
	}

	switch header.Command {
	case cmd.Announcement:
		var message cmd.AnnouncementCmd
		if err := json.Unmarshal(data, &message); err != nil {
			if p.Handlers.OnLog != nil {
				p.Handlers.OnLog(fmt.Sprintf("bad announcement from %s: %v", from, err))
			}
			return
		}

		if p.host == nil {
			if p.Handlers.OnLog != nil {
				p.Handlers.OnLog("cannot connect to announced peer: host not initialized")
			}
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := p.host.Connect(ctx, message.Peer); err != nil {
			p.Handlers.OnLog(fmt.Sprintf("connect to %s failed: %v", message.Peer.ID, err))
		}
	case cmd.Tick:
		var message cmd.TickCmd
		if err := json.Unmarshal(data, &message); err != nil {
			if p.Handlers.OnLog != nil {
				p.Handlers.OnLog(fmt.Sprintf("bad tick cmd from %s: %v", from, err))
			}
			return
		}
	default:
		if p.Handlers.OnLog != nil {
			p.Handlers.OnLog(fmt.Sprintf("unknown command %d from %s", header.Command, from))
		}
	}

	if p.Handlers.OnMessageReceived != nil {
		go p.Handlers.OnMessageReceived(from, data)
	}
}

func (p *P2p) readData(conn *peerConn) {
	for {
		payload, err := readFrame(conn.rw.Reader)
		if err != nil {
			if err != io.EOF && p.Handlers.OnLog != nil {
				p.Handlers.OnLog(fmt.Sprintf("read from %s failed: %v", conn.id, err))
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

func (p *P2p) StartHosting() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	listenF := flag.Int("l", 3000, "wait for incoming connections")
	targetF := flag.String("d", "", "target peer to dial")
	seedF := flag.Int64("seed", 0, "set random seed for id generation")
	flag.Parse()

	if *listenF == 0 {
		p.Handlers.OnLog("Please provide a port to bind on with -l")
	}

	h, err := p.makeHost(*listenF, *seedF)
	if err != nil {
		p.Handlers.OnLog(err.Error())
	}

	p.host = h

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

	<-ctx.Done()
}

func (p *P2p) makeHost(listenPort int, randSeed int64) (host.Host, error) {
	var r io.Reader
	if randSeed == 0 {
		r = rand.Reader
	} else {
		r = mrand.New(mrand.NewSource(randSeed))
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

func (p *P2p) startPeer(_ context.Context, h host.Host, streamHandler network.StreamHandler) {
	h.SetStreamHandler(protocolID, streamHandler)

	var port string
	for _, la := range h.Network().ListenAddresses() {
		if p, err := la.ValueForProtocol(multiaddr.P_TCP); err == nil {
			port = p
			break
		}
	}

	if port == "" && p.Handlers.OnLog != nil {
		p.Handlers.OnLog("was not able to find actual local port")
		return
	}
}

func (p *P2p) startPeerAndConnect(_ context.Context, h host.Host, destination string) (*peerConn, error) {
	for _, la := range h.Addrs() {
		p.Handlers.OnLog(fmt.Sprintf(" - %v", la))
	}

	maddr, err := ma.NewMultiaddr(destination)
	if err != nil {
		p.Handlers.OnLog(err.Error())
		return nil, err
	}

	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		p.Handlers.OnLog(err.Error())
		return nil, err
	}

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

	return conn, nil
}

func (p *P2p) getHostAddress(ha host.Host) string {
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/p2p/%s", ha.ID()))
	addr := ha.Addrs()[0]
	return addr.Encapsulate(hostAddr).String()
}
