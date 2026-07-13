package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

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
