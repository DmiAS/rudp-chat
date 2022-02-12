package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"rudp/pkg/protocol"
)

const network = "udp"

func main() {
	protocol.Debug()
	listener, err := protocol.ListenRUDP(
		&net.UDPAddr{
			Port: 9000,
		},
	)

	if err != nil {
		log.Fatalf("rudp.ListenRUDP error: %v", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("listener.Accept error: %v", err)
		}
		go func() {
			<-time.Tick(10 * time.Second)
			conn.Close()
		}()
		go func() {
			buf := make([]byte, 10)
			for {
				n, err := conn.Read(buf)
				if err != nil {
					fmt.Printf("conn.Read error: %v\n", err)
					return
				}
				buf = buf[:n]
				fmt.Printf(string(buf))
			}
		}()
	}
}
