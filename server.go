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

func (s *Server) Listen(ctx context.Context) {
	// we wait only for one connection thus we should not do it in loop
	conn, err := s.listener.AcceptKCP()
	if err != nil {
		log.Error().Err(err).Msg("failure to accept connection")
	}
	go s.read(ctx, conn)
	go s.write(ctx, conn)
	<-ctx.Done()
}

func (s *Server) read(ctx context.Context, session *kcp.UDPSession) {
	for {
		select {
		case <-ctx.Done():
			log.Debug().Msg("stop reading from session: context expired")
		default:
			data := make([]byte, buffSize)
			n, err := session.Read(data)
			if err != nil {
				log.Error().Err(err).Msg("failure to read from session in server")
			} else {
				log.Debug().Msg("sending data to receive channel")
				s.recvChan <- data[:n]
			}
		}
	}
}

func (s *Server) write(ctx context.Context, session *kcp.UDPSession) {
	for {
		select {
		case <-ctx.Done():
			log.Debug().Msg("stop reading from session: context expired")
		case data := <-s.sendChan:
			n, err := session.Write(data)
			if err != nil {
				log.Error().Err(err).Msgf("failure to write data(%d) to session in server", len(data))
			}
			log.Debug().Msgf("send %d bytes to session in server", n)
		}
	}
}
