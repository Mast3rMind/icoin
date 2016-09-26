package main

import (
	"errors"
	"github.com/daqing/icoin/wire"
	"log"
	"net"
	"time"
)

const (
	negotiateTimeout = 30 * time.Second
)

type peer struct {
	server   *server
	conn     net.Conn
	incoming bool
}

func newInboundPeer(s *server, conn net.Conn) *peer {
	return newPeer(s, conn, true)
}

func newOutboundPeer(s *server, conn net.Conn) *peer {
	return newPeer(s, conn, false)
}

func newPeer(s *server, conn net.Conn, incoming bool) *peer {
	return &peer{
		server:   s,
		conn:     conn,
		incoming: incoming,
	}
}

func (p *peer) ID() string {
	return p.conn.RemoteAddr().String()
}

func (p *peer) Start() error {
	ch := make(chan error)

	go func() {
		if p.incoming {
			ch <- p.NegotiateInboundProtocol()
		} else {
			ch <- p.NegotiateOutboundProtocol()
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return err
		}
	case <-time.After(negotiateTimeout):
		return errors.New("Protocol negotiation timeout")
	}

	go p.InHandler()

	return nil
}

func (p *peer) InHandler() {
loop:
	for {
		msg, err := p.readMessage()
		if err != nil {
			log.Println(err)
			break loop
		}

		log.Printf("Received Msg from remote peer: %v\n", msg)
	}
}

func (p *peer) Disconnect() {
	p.conn.Close()
}

func (p *peer) NegotiateInboundProtocol() error {
	return nil
	/*
		if err := p.readRemoteVersionMsg(); err != nil {
			return err
		}

		return p.writeLocalVersionMsg()
	*/
}

func (p *peer) NegotiateOutboundProtocol() error {
	return nil
}

func (p *peer) readRemoteVersionMsg() error {
	/*
		msg, err := p.readMessage()
		if err != nil {
			return err
		}

		versionMsg, ok := msg.(*wire.MsgVersion)
	*/

	return nil
}

func (p *peer) readMessage() (wire.Message, error) {
	msg, err := wire.ReadMessage(p.conn, p.server.netID())
	if err != nil {
		return nil, err
	}

	return msg, err
}
