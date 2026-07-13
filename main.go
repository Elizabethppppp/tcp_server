package main

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

func main() {
	mux := NewMux()
	mux.Handle("/time", handleTime)
	mux.Handle("/json", handleJSON)
	mux.Handle("/teapot", handleTeapot)
	mux.Handle("/hello", handleHello)

	if err := Listen(":8090", mux); err != nil {
		log.Fatal(err)
	}
}
