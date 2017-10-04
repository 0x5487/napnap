package napnap

import (
	"errors"
	"html/template"
	"net/http"
	"path"
	"sync"
)

var (
	_logger *logger
)

func init() {
	_logger = &logger{
		mode: off,
	}
}

// HandlerFunc defines a function to server HTTP requests
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

// WrapHandler wraps `http.Handler` into `napnap.HandlerFunc`.
func WrapHandler(h http.Handler) HandlerFunc {
	return func(c *Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

type NapNap struct {
	pool             sync.Pool
	handlers         []MiddlewareHandler
	middleware       middleware
	template         *template.Template
	templateRootPath string

	ForwardRemoteIpAddress bool
	MaxRequestBodySize     int64
}

// New returns a new NapNap instance
func New(mHandlers ...MiddlewareHandler) *NapNap {
	nap := &NapNap{
		handlers:           mHandlers,
		middleware:         build(mHandlers),
		MaxRequestBodySize: 10485760, // default 10MB for request body size
	}

	nap.pool.New = func() interface{} {
		rw := NewResponseWriter()
		return NewContext(nap, nil, rw)
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

// SetTemplate function allows user to set their own template instance.
func (nap *NapNap) SetTemplate(t *template.Template) {
	nap.template = t
}

// SetRender function allows user to set template location.
func (nap *NapNap) SetRender(templateRootPath string) {
	sharedTemplatePath := path.Join(templateRootPath, "shares/*")
	tmpl, err := template.ParseGlob(sharedTemplatePath)
	template := template.Must(tmpl, err)
	if template == nil {
		_logger.debug("no template")
		template = template.New("")
	}
	nap.template = template
	nap.templateRootPath = templateRootPath
}

// Run will run http server
func (nap *NapNap) Run(engine *Server) error {
	engine.Handler = nap
	return engine.ListenAndServe()
}

// RunTLS will run http/2 server
func (nap *NapNap) RunTLS(engine *Server) error {
	engine.Handler = nap
	return engine.ListenAndServeTLS(engine.Config.TLSCertFile, engine.Config.TLSKeyFile)
}

// RunAll will listen on multiple port
func (nap *NapNap) RunAll(addrs []string) error {
	if len(addrs) == 0 {
		return errors.New("addrs can't be empty")
	}

	wg := &sync.WaitGroup{}

	for _, addr := range addrs {
		wg.Add(1)
		go func(newAddr string) {
			err := http.ListenAndServe(newAddr, nap)
			if err != nil {
				panic(err)
			}
			wg.Done()
		}(addr)
	}

	wg.Wait()
	return nil
}

// Conforms to the http.Handler interface.
func (nap *NapNap) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	req.Body = http.MaxBytesReader(w, req.Body, nap.MaxRequestBodySize)
	c := nap.pool.Get().(*Context)
	c.reset(w, req)
	nap.middleware.Execute(c)
	nap.pool.Put(c)
}
