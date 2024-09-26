package main

import (
	"time"
	"y/p2p"
)


func makeServer (listenAddr string , nodes ...string) *FileServer{
	tcptransportOpts := p2p.TCPTransportOps{
		ListenAddress:    listenAddr,
		HandShakeFunc: p2p.NOPHandShakeFunc,
		Decoder:       p2p.DefaultDecoder{},
	}
	tcpTransport := p2p.NewTCPTransport(tcptransportOpts)

	fileServerOpts := FileServerOpts{
		StorageRoot:       listenAddr + "_network",
		PathTransformFunc: CASPathtranformFunc,
		Transport:         tcpTransport,
		BootStarpNodes:    nodes,
	}

	s := NewFileServer(fileServerOpts)

	tcpTransport.OnPeer = s.OnPeer

	return s
}

func main() {
	s1 := makeServer(":3000","")
	s2 := makeServer(":4000",":3000") 


	go s1.Start()
	time.Sleep(time.Second*2)

	go s2.Start()

}