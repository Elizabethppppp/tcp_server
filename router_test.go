package server

import "testing"

type fakeResponseWriter struct {
	status int
}

func (f *fakeResponseWriter) WriteHeader(status int) {
	f.status = status
}

func (f *fakeResponseWriter) SetHeader(key, value string) {
}

func (f *fakeResponseWriter) Write(body []byte) (int, error) {
	return len(body), nil
}

func newFakeWriter() *fakeResponseWriter {
	return &fakeResponseWriter{}
}

func TestMux_DispatchesToRegisteredHandler(t *testing.T) {
	m := NewMux()

	var called string
	m.Handle("/a", func(w ResponseWriter, r *Request) { called = "a" })
	m.Handle("/b", func(w ResponseWriter, r *Request) { called = "b" })

	m.ServeHTTP(newFakeWriter(), &Request{Method: "GET", Path: "/b"})

	if called != "b" {
		t.Errorf("вызван хендлер %q, want %q", called, "b")
	}
}

func TestMux_DispatchesToRegisteredHandlerPOST(t *testing.T) {
	m := NewMux()

	var calledP string
	m.Handle("/a", func(w ResponseWriter, r *Request) { calledP = "a" })
	m.Handle("/b", func(w ResponseWriter, r *Request) { calledP = "b" })

	m.ServeHTTP(newFakeWriter(), &Request{Method: "POST", Path: "/b"})

	if calledP != "b" {
		t.Errorf("вызван хендлер %q, want %q", calledP, "b")
	}
}

func TestMux_NotFound(t *testing.T) {
	m := NewMux()

	fw := newFakeWriter()

	m.ServeHTTP(fw, &Request{Method: "GET", Path: "/unknown"})

	if fw.status != 404 {
		t.Errorf("Статус получен %d, want %d", fw.status, 404)
	}
}

func TestMux_PathParam(t *testing.T) {
	m := NewMux()

	var userID string
	m.Handle("/user/{id}", func(w ResponseWriter, r *Request) {
		userID = r.Param("id")
	})

	m.ServeHTTP(newFakeWriter(), &Request{Method: "GET", Path: "/user/123"})

	if userID != "123" {
		t.Errorf("r.Param(\"id\") = %q, want %q", userID, "123")
	}
}
