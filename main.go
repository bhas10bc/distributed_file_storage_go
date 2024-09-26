package main

import (
	"bytes"
	"log"
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


	go func ()  {
		log.Fatal(s1.Start())
	}()
	time.Sleep(1*time.Second)
	go s2.Start()
	time.Sleep(1*time.Second)

	data := bytes.NewReader([]byte("my big data"))

	s2.StoreData("key", data)

	select{}

}