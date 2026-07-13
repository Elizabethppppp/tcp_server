package main

import (
	"bufio"
	"net"
	"time"
)

func handleFunc(conn net.Conn, handler Handler) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	r, err := parseRequest(reader)
	w := &responseWriter{}
	if err != nil {
		w.SetHeader("Content-Type", "text/plain")
		w.WriteHeader(400)
		w.Write([]byte("Bad request"))
		w.Flush(conn)
		return
	}

	handler.ServeHTTP(w, r)

	w.Flush(conn)
}

type HandlerFunc func(w ResponseWriter, r *Request)

func handleHello(w ResponseWriter, r *Request) {
	w.SetHeader("Content-Type", "text/plain")
	w.WriteHeader(200)
	w.Write([]byte("Hello, " + r.Path + "!"))
}

func handleTime(w ResponseWriter, r *Request) {
	now := time.Now()
	w.SetHeader("Content-Type", "text/plain")
	w.WriteHeader(200)
	w.Write([]byte(now.Format(time.RFC3339)))
}

func handleJSON(w ResponseWriter, r *Request) {
	w.SetHeader("Content-Type", "application/json")
	w.WriteHeader(200)

	w.Write([]byte(`{"ok": true}`))
}

func handleTeapot(w ResponseWriter, r *Request) {
	w.SetHeader("Content-Type", "text/plain")
	w.WriteHeader(418)
	w.Write([]byte("I'm a teapot"))
}

func handleNotFound(w ResponseWriter, r *Request) {
	w.WriteHeader(404)
	w.Write([]byte("404 Not Found"))
}
