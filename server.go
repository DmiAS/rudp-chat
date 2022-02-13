package rudp

import (
	"context"
	"fmt"
	"net"

	"github.com/rs/zerolog/log"

	"github.com/xtaci/kcp-go"
)

type Server struct {
	recvChan chan []byte
	sendChan chan []byte
	listener *kcp.Listener
}

func NewServer(conn net.PacketConn) (*Server, error) {
	block, err := createBlock()
	if err != nil {
		return nil, err
	}
	l, err := kcp.ServeConn(block, dataShards, parityShards, conn)
	if err != nil {
		return nil, fmt.Errorf("failure to initate kcp connection: %s", err)
	}
	return &Server{
		listener: l, recvChan: make(chan []byte, recvChanSize), sendChan: make(chan []byte, sendChanSize),
	}, nil
}

func (s *Server) Close() error {
	return s.listener.Close()
}

func (s *Server) Start(ctx context.Context) {
	// we wait only for one connection thus we should not do it in loop
	conn, err := s.listener.AcceptKCP()
	if err != nil {
		log.Error().Err(err).Msg("failure to accept connection")
	}
	m := NewManager(conn)
	go m.Read(ctx)
	go m.Write(ctx)
	<-ctx.Done()
}
