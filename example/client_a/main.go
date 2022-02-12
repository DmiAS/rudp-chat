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
	conn, err := protocol.DialRUDP(&net.UDPAddr{Port: 11000}, &net.UDPAddr{Port: 9000})
	if err != nil {
		log.Fatalf("rudp.DialRUDP error: %v", err)
	}
	fmt.Println(conn.LocalAddr())
	// go func() {
	// <-time.Tick(time.Second * 5)
	// conn.Close()
	// }()
	data := []string{"******", "hello,", " this", " is", " the", " rudp", " client", "******\n"}
	for cnt := 0; ; cnt++ {
		n, err := conn.Write([]byte(data[cnt%len(data)]))
		if err != nil {
			log.Fatalf("conn.Write error: %v", err)
		}
		fmt.Println(n)
		time.Sleep(time.Millisecond * 500)
	}
}
