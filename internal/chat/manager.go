package chat

import (
	"context"
	"fmt"
	"net"

	"chat/pkg/rudp"
)

type Engine interface {
	Start(ctx context.Context)
}

type Manager struct {
	conn net.PacketConn
	quit chan struct{}
}

func NewManager(conn net.PacketConn) *Manager {
	return &Manager{conn: conn, quit: make(chan struct{})}
}

func (m *Manager) Connect(addr string) error {
	cli, err := rudp.NewClient(m.conn, addr)
	if err != nil {
		return fmt.Errorf("failure to create rudp client: %s", err)
	}
	// run client in separate goroutine
	go m.spin(cli)
	return nil
}

func (m *Manager) spin(engine Engine) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	engine.Start(ctx)
	<-m.quit
}

func (m *Manager) Disconnect() {
	// send signal to interrupt current running engine
	m.quit <- struct{}{}
}
