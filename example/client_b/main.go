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
	laddr, err := net.ResolveUDPAddr(network, "localhost:9000")
	if err != nil {
		panic(err)
	}
	conn, err := net.ListenUDP(network, laddr)
	if err != nil {
		panic(err)
	}
	// listener, err := protocol.ListenRUDP(
	// 	&net.UDPAddr{
	// 		Port: 9000,
	// 	},
	// )
	listener := protocol.ListenRUDP(conn)

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
