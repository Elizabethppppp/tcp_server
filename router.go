package main

type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}

type Mux struct {
	routes map[string]HandlerFunc
}

func NewMux() *Mux {
	return &Mux{
		routes: make(map[string]HandlerFunc),
	}
}

func (m *Mux) Handle(path string, h HandlerFunc) {
	m.routes[path] = h
}

func (m *Mux) ServeHTTP(w ResponseWriter, r *Request) {
	if h, inMap := m.routes[r.Path]; inMap {
		h(w, r)
	} else {
		handleNotFound(w, r)
	}
}
