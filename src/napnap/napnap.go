package napnap

import "net/http"

const (
	// CONNECT HTTP method
	CONNECT = "CONNECT"
	// DELETE HTTP method
	DELETE = "DELETE"
	// GET HTTP method
	GET = "GET"
	// HEAD HTTP method
	HEAD = "HEAD"
	// OPTIONS HTTP method
	OPTIONS = "OPTIONS"
	// PATCH HTTP method
	PATCH = "PATCH"
	// POST HTTP method
	POST = "POST"
	// PUT HTTP method
	PUT = "PUT"
	// TRACE HTTP method
	TRACE = "TRACE"
)

type (
	NapNap struct {
		Router           *Router
		httpErrorHandler HTTPErrorHandler
	}

	NapNapHandleFunc func(c *Context) error
	// HTTPErrorHandler is a centralized HTTP error handler.
	HTTPErrorHandler func(error)
)

func New() *NapNap {
	return &NapNap{
		Router: NewRouter(),
	}
}

func (nap *NapNap) Get(path string, handler NapNapHandleFunc) {
	nap.Router.Add(GET, path, handler)
}

func (nap *NapNap) Post(path string, handler NapNapHandleFunc) {
	nap.Router.Add(POST, path, handler)
}

func (nap *NapNap) Run(addr string) {
	http.ListenAndServe(addr, nap)
}

// Conforms to the http.Handler interface.
func (nap *NapNap) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	c := NewContext(req, w)

	// Execute chain
	h := nap.Router.Find(req.Method, req.URL.Path)

	if err := h(c); err != nil {
		nap.httpErrorHandler(err)
	}
}
