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
	segments []segment
	handler  HandlerFunc
}

type Mux struct {
	routes []route
}

func NewMux() *Mux {
	return &Mux{
		routes: make([]route, 0),
	}
}

func (m *Mux) Handle(path string, h HandlerFunc) {
	s := strings.Split(strings.Trim(path, "/"), "/")

	seg := make([]segment, len(s))

	for i, value := range s {
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
		path:     path,
		segments: seg,
		handler:  h,
	}

	m.routes = append(m.routes, route)
}

func (m *Mux) ServeHTTP(w ResponseWriter, r *Request) {
	p := strings.Split(strings.Trim(r.Path, "/"), "/")

	for _, value := range m.routes {
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
