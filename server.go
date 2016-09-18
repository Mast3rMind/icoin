package main

import (
	"github.com/daqing/icoin/console"
	"log"
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

	waitChan := make(chan bool)

	go console.WaitForKill(waitChan)

	<-waitChan
}
