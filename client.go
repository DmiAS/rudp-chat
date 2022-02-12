package rudp

import (
	"fmt"
	"net"

	"github.com/l0k18/rudp"

	"rudp/pkg/protocol"
)

type Client struct {
	listener *rudp.RudpListener
	rconn    *protocol.RUDPConn
	conn     *net.UDPConn
}

const (
	buffSize = 1024
)

func NewClient(conn *net.UDPConn) *Client {
	return &Client{conn: conn}
}

func (c *Client) DialConn(laddr, raddr *net.UDPAddr) error {
	var err error
	c.rconn, err = protocol.DialRUDP(laddr, raddr)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Send(data []byte) error {
	if _, err := c.rconn.Write(data); err != nil {
		return err
	}
	return nil
}

func (c *Client) Listen(laddr *net.UDPAddr) {
	listener, err := protocol.ListenRUDP(laddr)
	if err != nil {
		panic(err)
	}
	for {
		rconn, err := listener.Accept()
		if err != nil {
			fmt.Printf("accept err %v\n", err)
			break
		}
		c.handleClient(rconn)
	}
}

func (c *Client) handleClient(conn net.Conn) {
	for {
		fmt.Printf(conn.RemoteAddr().String())
		data := make([]byte, buffSize)
		n, err := conn.Read(data)
		if err != nil {
			fmt.Printf("read err %s\n", err)
			break
		}
		fmt.Printf("receive %s", string(data[:n]))
		fmt.Printf(" from <%v>\n", conn.RemoteAddr())
	}
}
