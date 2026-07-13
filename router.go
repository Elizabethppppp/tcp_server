package main

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
