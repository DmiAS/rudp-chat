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
	laddr, err := net.ResolveUDPAddr(network, "localhost:11000")
	if err != nil {
		panic(err)
	}
	raddr, err := net.ResolveUDPAddr(network, "localhost:9000")
	if err != nil {
		panic(err)
	}
	conn, err := net.ListenUDP(network, laddr)
	if err != nil {
		panic(err)
	}

	rconn, err := protocol.DialRUDP(conn, raddr)
	if err != nil {
		log.Fatalf("rudp.DialRUDP error: %v", err)
	}
	fmt.Println(rconn.LocalAddr())
	data := []string{"******", "hello,", " this", " is", " the", " rudp", " client", "******\n"}
	for cnt := 0; ; cnt++ {
		n, err := rconn.Write([]byte(data[cnt%len(data)]))
		if err != nil {
			log.Fatalf("conn.Write error: %v", err)
		}
		fmt.Println(n)
		time.Sleep(time.Millisecond * 500)
	}
}
