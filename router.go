package server

import "strings"

type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}

type segment struct {
	value   string
	isParam bool
}

type route struct {
	path     string
	method   string
	segments []segment
	handler  HandlerFunc
}

type Mux struct {
	routes      []route
	middlewares []Middleware
}

func NewMux() *Mux {
	return &Mux{
		routes:      make([]route, 0),
		middlewares: make([]Middleware, 0),
	}
}

func (m *Mux) Use(mw Middleware) {
	m.middlewares = append(m.middlewares, mw)
}

func (m *Mux) Handle(r string, h HandlerFunc) {

	s := strings.SplitN(r, " ", 2)
	if len(s) != 2 {
		panic("invalid length: " + r)
	}
	if s[0] == "" || s[1] == "" {
		panic("invalid route: " + r)
	}

	method := strings.ToUpper(s[0])
	path := s[1]
	segments := strings.Split(strings.Trim(path, "/"), "/")

	seg := make([]segment, len(segments))

	for i, value := range segments {
		if strings.HasPrefix(value, "{") && strings.HasSuffix(value, "}") {
			seg[i] = segment{
				value:   strings.Trim(value, "{}"),
				isParam: true,
			}
		} else {
			seg[i] = segment{
				value:   value,
				isParam: false,
			}
		}
	}

	route := route{
		method:   method,
		path:     path,
		segments: seg,
		handler:  h,
	}

	m.routes = append(m.routes, route)
}

func (m *Mux) ServeHTTP(w ResponseWriter, r *Request) {

	p := strings.Split(strings.Trim(r.Path, "/"), "/")

	for _, value := range m.routes {
		if value.method != r.Method {
			continue
		}
		if len(p) != len(value.segments) {
			continue
		}
		params := make(map[string]string)
		match := true

		for i, segm := range value.segments {
			if segm.isParam {
				params[segm.value] = p[i]
			} else if segm.value != p[i] {
				match = false
				break
			}
		}

		if match {
			r.param = params
			value.handler(w, r)
			return
		}

	}
	w.WriteHeader(StatusNotFound)
	w.Write([]byte("Not Found"))
}
