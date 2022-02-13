package rudp

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/xtaci/kcp-go"
)

type Manager struct {
	recvChan chan []byte
	sendChan chan []byte
	session  *kcp.UDPSession
}

func NewManager(session *kcp.UDPSession) *Manager {
	return &Manager{
		recvChan: make(chan []byte, recvChanSize),
		sendChan: make(chan []byte, sendChanSize),
		session:  session,
	}
}

func (m *Manager) Write(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Debug().Msg("stop reading from session: context expired")
		default:
			data := make([]byte, buffSize)
			n, err := m.session.Read(data)
			if err != nil {
				log.Error().Err(err).Msg("failure to read from session")
			} else {
				log.Debug().Msg("sending data to receive channel")
				m.recvChan <- data[:n]
			}
		}
	}
}

func (m *Manager) Read(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Debug().Msg("stop reading from session: context expired")
		case data := <-m.sendChan:
			n, err := m.session.Write(data)
			if err != nil {
				log.Error().Err(err).Msgf("failure to write data(%d) to session", len(data))
			}
			log.Debug().Msgf("send %d bytes to session", n)
		}
	}
}
