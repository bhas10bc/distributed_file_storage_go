package p2p

import (
	"errors"
	"fmt"
	"log"
	"net"
)

type TCPPeer struct {
	net.Conn
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer{
	return &TCPPeer{
		Conn :conn,
		outbound: outbound,
	}
}

func (p *TCPPeer) Send(b []byte) error {
	_, err := p.Conn.Write(b)
	return err
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

func (t *TCPTransport) Close()error {
	return t.listener.Close()
}

func (t *TCPTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp",addr)
	if err != nil {
		return err
	}

	go t.handleConn(conn, true)
	return nil
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listener , err = net.Listen("tcp",t.ListenAddress)
	if err != nil{
		return err
	}
	go t.startAcceptLoop()

	log.Printf("TCP listeneing on %s", t.ListenAddress)
	return nil
}

func (t *TCPTransport) startAcceptLoop(){
	for {
		conn, err := t.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			return 
		}
		if err != nil {
			fmt.Printf("TCP accept error : %s", err)
		}

		fmt.Printf("new incoming connection %v", conn)

		go t.handleConn(conn, false)
	}
}
type Temp struct{}

func (t *TCPTransport) handleConn(conn net.Conn, outbound bool) {

	var err error
	peer := NewTCPPeer(conn, outbound)

	if err = t.HandShakeFunc(peer); err != nil{
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
		rpc.From = conn.RemoteAddr().String()
		t.rpcch <- rpc
	}

}

