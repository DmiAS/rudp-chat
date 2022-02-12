package main

import (
	"crypto/sha1"
	"fmt"
	"net"
	"time"

	"github.com/xtaci/kcp-go"
	"golang.org/x/crypto/pbkdf2"
)

const network = "udp"

func main() {
	// raddr, err := net.ResolveUDPAddr(network, "localhost:9000")
	// if err != nil {
	// 	panic(err)
	// }
	// raddr := "localhost:9000"
	key := pbkdf2.Key([]byte("demo pass"), []byte("demo salt"), 1024, 32, sha1.New)
	block, err := kcp.NewAESBlockCrypt(key)
	if err != nil {
		panic(err)
	}
	conn, err := net.ListenPacket(network, "localhost:11000")
	if err != nil {
		panic(err)
	}
	c, err := kcp.ServeConn(block, 10, 3, conn)
	if err != nil {
		panic(err)
	}
	defer c.Close()
	go read(c)
	// go write(c)
	<-time.After(time.Second * 20)
}

func read(conn *kcp.Listener) {
	fmt.Println("start reading")
	for {
		s, err := conn.AcceptKCP()
		if err != nil {
			panic(err)
		}
		go handle(s)
	}
}

func handle(s *kcp.UDPSession) {
	for {
		data := make([]byte, 1024)
		n, err := s.Read(data)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(data[:n]))
		s.Write([]byte("pong"))
	}
}

func write(conn *kcp.UDPSession) {
	fmt.Println("start writing")
	for {
		_, err := conn.Write([]byte("hello, world 1"))
		if err != nil {
			fmt.Println("cant write", err)
		}
		<-time.After(time.Second * 3)
	}
}
