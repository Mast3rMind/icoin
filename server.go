package main

import (
	"github.com/daqing/icoin/console"
	"github.com/daqing/icoin/wire"
	"log"
	"net"
)

type server struct {
	conf  *config
	peers map[string]*peer
}

func newServer(conf *config) *server {
	s := &server{
		conf:  conf,
		peers: make(map[string]*peer),
	}

	return s
}

func (s *server) Start() {
	log.Printf("Server started on port: %s.", s.conf.params.port)

	go s.Listener()
	go s.ConnectToPeer()

	waitChan := make(chan bool)

	go console.WaitForKill(waitChan)

	<-waitChan
}

func (s *server) Listener() {
	ln, err := net.Listen("tcp", ":"+s.conf.params.port)

	if err != nil {
		log.Fatalln(err)
	}

	for {
		conn, err := ln.Accept()

		if err != nil {
			continue
		}

		peer := newInboundPeer(s, conn)

		s.AddPeer(peer)

	}
}

func (s *server) ConnectToPeer() {
	log.Println("Peer To Connect:", s.conf.connect)
	conn, err := net.Dial("tcp", s.conf.connect+":"+s.conf.params.port)
	if err != nil {
		log.Println(err)
		return
	}

	peer := newOutboundPeer(s, conn)
	s.AddPeer(peer)
}

func (s *server) AddPeer(peer *peer) {
	log.Printf("Added New Peer %s (%s) To Server\n", peer.ID(), peer.Direction())
	s.peers[peer.ID()] = peer

	go func() {
		if err := peer.Start(); err != nil {
			log.Println(err)
			peer.Disconnect()
			s.RemovePeer(peer)
		}
	}()
}

func (s *server) RemovePeer(peer *peer) {
	delete(s.peers, peer.ID())
}

func (s *server) netID() wire.NetID {
	return s.conf.params.netID
}
