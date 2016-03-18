package napnap

import (
	"html/template"
	"net/http"
	"sync"
)

type HandlerFunc func(c *Context)

// MiddlewareHandler is an interface that objects can implement to be registered to serve as middleware
// in the NapNap middleware stack.
type MiddlewareHandler interface {
	Invoke(c *Context, next HandlerFunc)
}

// MiddlewareFunc is an adapter to allow the use of ordinary functions as NapNap handlers.
type MiddlewareFunc func(c *Context, next HandlerFunc)

func (m MiddlewareFunc) Invoke(c *Context, next HandlerFunc) {
	m(c, next)
}

type middleware struct {
	handler MiddlewareHandler
	next    *middleware
}

func (m middleware) Execute(c *Context) {
	m.handler.Invoke(c, m.next.Execute)
}

type NapNap struct {
	pool       sync.Pool
	handlers   []MiddlewareHandler
	middleware middleware
	template   *template.Template
	//httpErrorHandler HTTPErrorHandler
}

// New returns a new NapNap instance
func New(mHandlers ...MiddlewareHandler) *NapNap {
	nap := &NapNap{
		handlers:   mHandlers,
		middleware: build(mHandlers),
	}

	nap.pool.New = func() interface{} {
		return NewContext(nap, nil, nil)
	}

	return nap
}

// UseFunc adds an anonymous function onto middleware stack.
func (nap *NapNap) UseFunc(aFunc func(c *Context, next HandlerFunc)) {
	nap.Use(MiddlewareFunc(aFunc))
}

// Use adds a Handler onto the middleware stack. Handlers are invoked in the order they are added to a NapNap.
func (nap *NapNap) Use(mHandler MiddlewareHandler) {
	nap.handlers = append(nap.handlers, mHandler)
	nap.middleware = build(nap.handlers)
}

func build(handlers []MiddlewareHandler) middleware {
	var next middleware

	if len(handlers) == 0 {
		return voidMiddleware()
	} else if len(handlers) > 1 {
		next = build(handlers[1:])
	} else {
		next = voidMiddleware()
	}

	return middleware{handlers[0], &next}
}

func voidMiddleware() middleware {
	return middleware{
		MiddlewareFunc(func(c *Context, next HandlerFunc) {}),
		&middleware{},
	}
}

func (nap *NapNap) SetViews(path string) {
	tmpl, err := template.ParseGlob(path)

	if err != nil {
		panic(err)
	}

	template := template.Must(tmpl, err)
	if template == nil {
		println("no template")
	}
	nap.template = template
}

// Run http server
func (nap *NapNap) Run(addr string) {
	//fmt.Println(fmt.Sprintf("listening on %s", addr))
	http.ListenAndServe(addr, nap)
}

// Conforms to the http.Handler interface.
func (nap *NapNap) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := nap.pool.Get().(*Context)
	c.reset(req, w)
	nap.middleware.Execute(c)
	nap.pool.Put(c)
}
