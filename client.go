package rudp

import (
	"context"
	"fmt"
	"net"

	"github.com/xtaci/kcp-go"
)

type Client struct {
	c *kcp.UDPSession
}

const (
	buffSize = 1024
)

func NewClient(conn net.PacketConn, raddr string) (*Client, error) {
	block, err := createBlock()
	if err != nil {
		return nil, err
	}
	c, err := kcp.NewConn(raddr, block, dataShards, parityShards, conn)
	if err != nil {
		return nil, fmt.Errorf("failure to initate kcp connection: %s", err)
	}
	return &Client{c: c}, nil
}

func (c *Client) Start(ctx context.Context) {
	m := NewManager(c.c)
	go m.Read(ctx)
	go m.Write(ctx)
	// wait for interrupt
	<-ctx.Done()
}

func (c *Client) Close() error {
	return c.c.Close()
}
