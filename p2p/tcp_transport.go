package p2p

import (
	"fmt"
	"net"
)

type TCPPeer struct {
	conn net.Conn
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer{
	return &TCPPeer{
		conn:conn,
		outbound: outbound,
	}
}

func (p *TCPPeer) Close() error {
	return p.Close()
}

type TCPTransportOps struct {
	ListenAddress string
	Decoder DefaultDecoder
	HandShakeFunc HandShakeFunc
	OnPeer func(Peer) error
}

type TCPTransport struct {
	TCPTransportOps
	listener      net.Listener
	rpcch chan RPC
}

func NewTCPTransport(ops TCPTransportOps) *TCPTransport{
	return &TCPTransport{
		TCPTransportOps: ops,		
	}
}

func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcch
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listener , err = net.Listen("tcp",t.ListenAddress)
	if err != nil{
		return err
	}
	t.startAcceptLoop()

	return nil
}

func (t *TCPTransport) startAcceptLoop(){
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP accept error : %s", err)
		}

		go t.handleConn(conn)
	}
}
type Temp struct{}

func (t *TCPTransport) handleConn(conn net.Conn) {

	var err error
	peer := NewTCPPeer(conn, true)

	if err = t.HandShakeFunc(peer); err != nil{
		conn.Close()
		fmt.Printf("TCP handshake error: %s", err)
		return
	}

	if t.OnPeer != nil {
		err = t.OnPeer(peer)
		if err != nil {
			return
		}
	}

	rpc := RPC{}

	for {
		
		err = t.Decoder.Decode(conn, &rpc)

		if err != nil{
			fmt.Printf("TCP error %s", err)
			return
		}
		rpc.From = conn.RemoteAddr()
		t.rpcch <- rpc
	}

}

