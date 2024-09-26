package main

import (
	"fmt"
	"log"
	"sync"
	"y/p2p"
)

type FileServerOpts struct {
	StorageRoot       string
	PathTransformFunc PathTransformFunc
	Transport         p2p.Transport
	BootStarpNodes []string
}

type FileServer struct {
	FileServerOpts
	peerLock sync.Mutex
	peers map[string]p2p.Peer
	store *Store
	quitch chan struct{}
}

func NewFileServer(opts FileServerOpts) *FileServer{
	storeOpts := StoreOpts{
		Root: opts.StorageRoot,
		PathtransformFunc: opts.PathTransformFunc,
	}
	return &FileServer{
		FileServerOpts: opts,
		store: NewStore(storeOpts),
		quitch: make(chan struct{}),
		peers: make(map[string]p2p.Peer),
	}
}

func(s *FileServer) Stop(){
	close(s.quitch)
}

func (s *FileServer) OnPeer( p p2p.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()
	s.peers[p.RemoteAddr().String()] = p
	log.Printf("connected with peer %s",p.RemoteAddr() )
	return nil
}

func (s *FileServer) loop(){
	defer func (){
		log.Println("file server stopped")
		s.Transport.Close()
	}()
	for {
		select{
		case msg := <- s.Transport.Consume():
			fmt.Println(msg)
		
		case <- s.quitch:
			return
	}
}
}

func (s *FileServer) bootStarpNetwork() error {
	for _, addr := range s.BootStarpNodes {
		if len(addr) == 0 {
			continue
		}

		go func(addr string) {
			fmt.Printf("[%s] attemping to connect with remote \n", addr)
			if err := s.Transport.Dial(addr); err != nil {
				log.Println("dial error: ", err)
			}
		}(addr)
	}

	return nil
}
func(s *FileServer) Start() error{
	err := s.Transport.ListenAndAccept()
	if err != nil {
		return err
	}
	
	s.bootStarpNetwork()
	s.loop()
	return nil
}