package server

import (
	"fmt"
	"net"
	"strings"
)

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
