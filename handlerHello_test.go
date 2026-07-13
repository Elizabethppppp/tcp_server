package main

import (
	"bytes"
	"testing"
)

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
