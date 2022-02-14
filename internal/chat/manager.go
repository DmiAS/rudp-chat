package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"net"

	"github.com/rs/zerolog/log"

	"chat/internal/model"
	"chat/pkg/rudp"
)

type Engine interface {
	Start(ctx context.Context)
	Receive() <-chan []byte
	Send(data []byte)
}

const (
	chanSize = 5
)

type Manager struct {
	filesChan    chan []byte
	messagesChan chan []byte
	conn         net.PacketConn
	quit         chan struct{}
	engine       Engine
}

func NewManager(conn net.PacketConn) *Manager {
	return &Manager{
		filesChan:    make(chan []byte, chanSize),
		messagesChan: make(chan []byte, chanSize),
		conn:         conn,
		quit:         make(chan struct{}),
	}
}

func (m *Manager) StartServer() error {
	srv, err := rudp.NewServer(m.conn)
	if err != nil {
		return fmt.Errorf("failure to create rudp server: %s", err)
	}
	m.run(srv)
	log.Debug().Msg("init engine as client")
	return nil
}

func (m *Manager) Connect(addr string) error {
	cli, err := rudp.NewClient(m.conn, addr)
	if err != nil {
		return fmt.Errorf("failure to create rudp client: %s", err)
	}
	m.run(cli)
	log.Debug().Msg("init engine as client")
	return nil
}

func (m *Manager) run(engine Engine) {
	m.engine = engine
	go m.listen()
	go m.engine.Start(context.Background())
}

func (m *Manager) Disconnect() {
	// send signal to interrupt current running engine
	m.quit <- struct{}{}
}

// listens for incoming packets and send it to appropriate channels
func (m *Manager) listen() {
	for {
		data := <-m.engine.Receive()
		log.Debug().Msgf("read data %d from manager", len(data))

		packet := &model.Packet{}
		if err := json.Unmarshal(data, packet); err != nil {
			log.Error().Err(err).Msgf("failure to unmarshal data")
			continue
		}

		// log.Debug().Msgf("received data from engine: %+v", packet)
		switch {
		case packet.IsMessage:
			m.messagesChan <- packet.Data
		case packet.IsFile:
			// sends full data to process files
			m.filesChan <- data
		}
	}
}

func (m *Manager) SendData(data []byte) {
	log.Debug().Msgf("sending data %d bytes", len(data))
	m.engine.Send(data)
}

func (m *Manager) ReceiveFile() []byte {
	return <-m.filesChan
}

func (m *Manager) ReceiveMessage() []byte {
	return <-m.messagesChan
}
