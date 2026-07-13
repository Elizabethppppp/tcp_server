package server

import (
	"log"
	"net"
)

func Listen(addr string, handler Handler) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("Error", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting", err)
			continue
		}
		go handleFunc(conn, handler)
	}
}
