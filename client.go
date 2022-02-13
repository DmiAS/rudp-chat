package rudp

import (
	"crypto/sha1"
	"fmt"
	"net"

	"github.com/xtaci/kcp-go"
	"golang.org/x/crypto/pbkdf2"
)

type Client struct {
	c *kcp.UDPSession
}

const (
	buffSize = 1024
)

func NewClient(conn net.PacketConn, raddr string) (*Client, error) {
	// generate key to secure connection for kcp
	key := pbkdf2.Key([]byte(pass), []byte(salt), iter, keyLen, sha1.New)
	block, err := kcp.NewAESBlockCrypt(key)
	if err != nil {
		return nil, fmt.Errorf("failure to create aes block: %s", err)
	}
	c, err := kcp.NewConn(raddr, block, dataShards, parityShards, conn)
	if err != nil {
		return nil, fmt.Errorf("failure to initate kcp connection: %s", err)
	}
	return &Client{c: c}, nil
}
