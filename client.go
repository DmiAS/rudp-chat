package rudp

import (
	"context"
	"fmt"
	"net"

	"github.com/xtaci/kcp-go"
)

type Client struct {
	recvChan chan []byte
	sendChan chan []byte
	session  *kcp.UDPSession
}

func NewClient(conn net.PacketConn, raddr string) (*Client, error) {
	block, err := createBlock()
	if err != nil {
		return nil, err
	}
	session, err := kcp.NewConn(raddr, block, dataShards, parityShards, conn)
	if err != nil {
		return nil, fmt.Errorf("failure to initate kcp connection: %s", err)
	}
	return &Client{
		recvChan: make(chan []byte, recvChanSize),
		sendChan: make(chan []byte, sendChanSize),
		session:  session,
	}, nil
}

func (c *Client) Start(ctx context.Context) {
	m := NewManager(c.session, c.recvChan, c.sendChan)
	go m.Read(ctx)
	go m.Write(ctx)
	<-ctx.Done()
}

func (c *Client) Receive() <-chan []byte {
	return c.recvChan
}

func (c *Client) Send(data []byte) {
	c.sendChan <- data
}

func (c *Client) Close() error {
	return c.session.Close()
}
