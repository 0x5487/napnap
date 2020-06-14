package napnap

import (
	"crypto/tls"
	"errors"
	"html/template"
	"net/http"
	"path"
	"strings"
	"sync"

	"golang.org/x/crypto/acme/autocert"
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
type HandlerFunc func(c *Context) error

// ErrorHandler defines a function to handle HTTP errors
type ErrorHandler func(c *Context, err error)

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

func (m middleware) Execute(c *Context) error {
	m.handler.Invoke(c, m.next.Execute)
	return nil
}

// WrapHandler wraps `http.Handler` into `napnap.HandlerFunc`.
func WrapHandler(h http.Handler) HandlerFunc {
	return func(c *Context) error {
		h.ServeHTTP(c.Writer, c.Request)
		return nil
	}
}

type NapNap struct {
	pool             sync.Pool
	handlers         []MiddlewareHandler
	middleware       middleware
	template         *template.Template
	templateRootPath string
	router           *Router

	MaxRequestBodySize int64
	ErrorHandler       ErrorHandler
	NotFoundHandler    HandlerFunc
}

// New returns a new NapNap instance
func New(mHandlers ...MiddlewareHandler) *NapNap {
	nap := &NapNap{
		handlers:           mHandlers,
		middleware:         build(mHandlers),
		MaxRequestBodySize: 10485760, // default 10MB for request body size
	}

	nap.router = NewRouter(nap)
	nap.Use(nap.router)
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
	if len(nap.handlers) == 0 {
		nap.handlers = append(nap.handlers, mHandler)
	} else {
		end := len(nap.handlers) - 1
		nap.handlers = append(nap.handlers[:end], mHandler, nap.router)
	}

	nap.middleware = build(nap.handlers)
}

// All is a shortcut for adding all methods
func (nap *NapNap) All(path string, handler HandlerFunc) {
	nap.router.Add(GET, path, handler)
	nap.router.Add(POST, path, handler)
	nap.router.Add(PUT, path, handler)
	nap.router.Add(DELETE, path, handler)
	nap.router.Add(PATCH, path, handler)
	nap.router.Add(OPTIONS, path, handler)
	nap.router.Add(HEAD, path, handler)
}

// Get is a shortcut for router.Add("GET", path, handle)
func (nap *NapNap) Get(path string, handler HandlerFunc) {
	nap.router.Add(GET, path, handler)
}

// Post is a shortcut for router.Add("POST", path, handle)
func (nap *NapNap) Post(path string, handler HandlerFunc) {
	nap.router.Add(POST, path, handler)
}

// Put is a shortcut for router.Add("PUT", path, handle)
func (nap *NapNap) Put(path string, handler HandlerFunc) {
	nap.router.Add(PUT, path, handler)
}

// Delete is a shortcut for router.Add("DELETE", path, handle)
func (nap *NapNap) Delete(path string, handler HandlerFunc) {
	nap.router.Add(DELETE, path, handler)
}

// Patch is a shortcut for router.Add("PATCH", path, handle)
func (nap *NapNap) Patch(path string, handler HandlerFunc) {
	nap.router.Add(PATCH, path, handler)
}

// Options is a shortcut for router.Add("OPTIONS", path, handle)
func (nap *NapNap) Options(path string, handler HandlerFunc) {
	nap.router.Add(OPTIONS, path, handler)
}

// Head is a shortcut for router.Add("HEAD", path, handle)
func (nap *NapNap) Head(path string, handler HandlerFunc) {
	nap.router.Add(HEAD, path, handler)
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

func (nap *NapNap) RunAutoTLS(engine *Server) error {

	whiteLists := []string{}

	for _, domain := range strings.Split(engine.Config.Domain, ",") {
		whiteLists = append(whiteLists, strings.TrimSpace(domain))
	}

	m := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(whiteLists...),
	}

	if engine.Config.CertCachePath != "" {
		m.Cache = autocert.DirCache(engine.Config.CertCachePath)
	}

	go http.ListenAndServe(":http", m.HTTPHandler(nap))

	// https' settings
	engine.Addr = ":https"
	engine.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}
	engine.Handler = nap
	return engine.ListenAndServeTLS("", "")
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
