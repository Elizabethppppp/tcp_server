package main

import (
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":8090")
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
		go handleFunc(conn)
	}
}
