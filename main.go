package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
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

func handleFunc(conn net.Conn) {
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

	router(w, r)

	w.Flush(conn)
}

func router(w ResponseWriter, r *Request) {
	switch r.Path {
	case "/time":
		handleTime(w, r)
	case "/json":
		handleJSON(w, r)
	case "/teapot":
		handleTeapot(w, r)
	case "/hello":
		handleHello(w, r)
	default:
		handleNotFound(w, r)
	}
}

type Request struct {
	Method  string
	Path    string
	Proto   string
	Headers map[string]string
	Body    []byte
}

func parseRequestLine(line string) (method, requestURI, proto string, ok bool) {
	method, after, found1 := strings.Cut(line, " ")
	requestURI, proto, found2 := strings.Cut(after, " ")
	if !found1 || !found2 {
		return "", "", "", false
	}
	return method, requestURI, proto, true
}

func parseRequest(r *bufio.Reader) (*Request, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return nil, err
	}

	line = strings.TrimSpace(line)

	method, path, proto, ok := parseRequestLine(line)
	if !ok {
		return nil, fmt.Errorf("invalid request line")
	}

	m := make(map[string]string)

	for {
		s, err := r.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("invalid reading")
		}
		s = strings.TrimSpace(s)
		if s == "" {
			break
		}

		slice := strings.Split(s, ": ")
		key := strings.ToLower(slice[0])
		val := slice[1]
		m[key] = val
	}

	var body []byte
	value, ok := m["content-length"]
	if ok {
		l, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("invalid format")
		}

		if l > 0 {
			b := make([]byte, l)
			_, err := io.ReadFull(r, b)
			if err != nil {
				return nil, fmt.Errorf("invalid length")
			}
			body = b
		}
	}

	request := &Request{
		Method:  method,
		Path:    path,
		Proto:   proto,
		Headers: m,
		Body:    body,
	}

	return request, nil
}

type ResponseWriter interface {
	WriteHeader(status int)
	SetHeader(key, value string)
	Write(body []byte) (int, error)
}

type responseWriter struct {
	status       int
	headers      map[string]string
	body         []byte
	wroteHeaders bool
	wroteBody    bool
}

func (w *responseWriter) WriteHeader(status int) {
	if w.wroteHeaders {
		return
	}
	w.status = status
	w.wroteHeaders = true
}

func (w *responseWriter) SetHeader(key, value string) {
	if w.wroteBody || w.wroteHeaders {
		return
	}
	if w.headers == nil {
		w.headers = make(map[string]string)
	}
	w.headers[key] = value
}

func (w *responseWriter) Write(body []byte) (int, error) {
	if !w.wroteHeaders {
		w.WriteHeader(200)
	}
	w.body = append(w.body, body...)
	w.wroteBody = true
	return len(body), nil
}

type HandlerFunc func(w ResponseWriter, r *Request)

func (w *responseWriter) Flush(conn net.Conn) error {
	if w.status == 0 {
		w.status = 200
	}

	var statusText string
	if w.status == 400 {
		statusText = "Bad request"
	} else if w.status == 418 {
		statusText = "I'm a teapot"
	} else {
		statusText = "OK"
	}

	var response strings.Builder
	response.WriteString(fmt.Sprintf("HTTP/1.1 %d %s\r\n", w.status, statusText))
	for key, val := range w.headers {
		response.WriteString(fmt.Sprintf("%s: %s\r\n", key, val))
	}
	if _, exists := w.headers["Content-Length"]; !exists {
		response.WriteString(fmt.Sprintf("Content-Length: %d\r\n", len(w.body)))
	}
	response.WriteString("Connection: close\r\n")
	response.WriteString("\r\n")
	response.WriteString(string(w.body))

	_, err := conn.Write([]byte(response.String()))
	return err
}

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
