package main

import (
	"fmt"
	"log"
	"y/p2p"
)

func OnPeer(peer p2p.Peer) error {
	peer.Close()
	// fmt.Println("doing some logic with peer outside TCP")
	return nil
}

func main() {
	tcpOps := p2p.TCPTransportOps{
		ListenAddress: ":3000",
		HandShakeFunc: p2p.NOPHandShakeFunc,
		Decoder: p2p.DefaultDecoder{},
		OnPeer: OnPeer,
	}
	tr := p2p.NewTCPTransport(tcpOps)

	go func() {
		for {
			msg := <- tr.Consume()
			fmt.Printf("%+v \n", msg)
		}
	}()
	err := tr.ListenAndAccept();
	if  err != nil {
		log.Fatal(err)
	}

	select{}
}