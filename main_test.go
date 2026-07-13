package main

import (
	"bytes"
	"testing"
)

func TestParseRequestLine(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		line       string
		wantMethod string
		wantURI    string
		wantProto  string
		wantOK     bool
	}{
		{
			name:       "валидный GET запрос",
			line:       "GET /hello HTTP/1.1",
			wantMethod: "GET",
			wantURI:    "/hello",
			wantProto:  "HTTP/1.1",
			wantOK:     true,
		},
		{
			name:       "валидный POST запрос",
			line:       "POST /api/user HTTP/1.1",
			wantMethod: "POST",
			wantURI:    "/api/user",
			wantProto:  "HTTP/1.1",
			wantOK:     true,
		},
		{
			name: "пустая строка",
			line: "",
		},
		{
			name:       "один элемент (нет версии)",
			line:       "GET /hello",
			wantMethod: "",
			wantURI:    "",
			wantProto:  "",
			wantOK:     false,
		},
		{
			name:       "один элемент вместо трёх",
			line:       "GET",
			wantMethod: "",
			wantURI:    "",
			wantProto:  "",
			wantOK:     false,
		},
		{
			name:       "четыре элемента",
			line:       "GET / HTTP/1.1 EXTRA",
			wantMethod: "GET",
			wantURI:    "/",
			wantProto:  "HTTP/1.1 EXTRA",
			wantOK:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMethod, gotURI, gotProto, gotOK := parseRequestLine(tt.line)

			if gotMethod != tt.wantMethod {
				t.Errorf("Method: got %q, want %q", gotMethod, tt.wantMethod)
			}
			if gotURI != tt.wantURI {
				t.Errorf("URI: got %q, want %q", gotURI, tt.wantURI)
			}
			if gotProto != tt.wantProto {
				t.Errorf("Proto: got %q, want %q", gotProto, tt.wantProto)
			}
			if gotOK != tt.wantOK {
				t.Errorf("OK: got %v, want %v", gotOK, tt.wantOK)
			}
		})
	}
}

// test for handler
type trainResponseWriter struct {
	status  int
	headers map[string]string
	body    *bytes.Buffer
}

func newTrainResponseWriter() *trainResponseWriter {
	return &trainResponseWriter{
		headers: make(map[string]string),
		body:    &bytes.Buffer{},
	}
}

func (tr *trainResponseWriter) WriteHeader(status int) {
	tr.status = status
}

func (tr *trainResponseWriter) SetHeader(key, value string) {
	tr.headers[key] = value
}

func (tr *trainResponseWriter) Write(body []byte) (int, error) {
	// tr.body = append(tr.body, body...)
	// return len(body), nil
	return tr.body.Write(body)
}

func TestHandleHello(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		path     string
		wantBody string
	}{
		{
			name:     "путь есть",
			path:     "/world",
			wantBody: "Hello, /world!",
		},
		{
			name:     "пустой путь",
			path:     "",
			wantBody: "Hello, !",
		},
		{
			name:     "путь с пробелом",
			path:     "/hello world",
			wantBody: "Hello, /hello world!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := newTrainResponseWriter()
			r := &Request{Path: tt.path}

			handleHello(w, r)

			if w.status != 200 {
				t.Errorf("status: got %d, want 200", w.status)
			}
			if w.body.String() != tt.wantBody {
				t.Errorf("body: got %q, want %q", w.body.String(), tt.wantBody)
			}
			if w.headers["Content-Type"] != "text/plain" {
				t.Errorf("Content-Type: got %q, want %q", w.headers["Content-Type"], "text/plain")
			}
		})
	}
}
