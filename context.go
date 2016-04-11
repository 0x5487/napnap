package napnap

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/url"
	"strings"
)

// Param is a single URL parameter, consisting of a key and a value.
type Param struct {
	Key   string
	Value string
}

type store map[string]interface{}

type Context struct {
	NapNap  *NapNap
	Request *http.Request
	Writer  http.ResponseWriter
	query   url.Values
	params  []Param
	store   store
}

// NewContext returns a new context instance
func NewContext(napnap *NapNap, req *http.Request, writer http.ResponseWriter) *Context {
	return &Context{
		NapNap:  napnap,
		Request: req,
		Writer:  writer,
	}
}

// Render returns html format
func (c *Context) Render(code int, viewName string, data interface{}) (err error) {
	c.Writer.WriteHeader(code)
	c.NapNap.template.ExecuteTemplate(c.Writer, viewName, data)
	return
}

// String returns string format
func (c *Context) String(code int, s string) (err error) {
	c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.Writer.WriteHeader(code)
	c.Writer.Write([]byte(s))
	return
}

// JSON returns json format
func (c *Context) JSON(code int, i interface{}) (err error) {
	b, err := json.Marshal(i)
	if err != nil {
		return err
	}
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(code)
	c.Writer.Write(b)
	return
}

// BindJSON binds the request body into provided type `obj`. The default binder does
// it based on Content-Type header.
func (c *Context) BindJSON(obj interface{}) (err error) {
	req := c.Request
	contentType := req.Header.Get("Content-Type")

	if contentType == "application/json" {
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(obj)
		if err != nil {
			return err
		}
	} else {
		return errors.New("content type doesn't match application/json")
	}

	return
}

// Query returns query parameter by key.
func (c *Context) Query(key string) string {
	if c.query == nil {
		c.query = c.Request.URL.Query()
	}
	return c.query.Get(key)
}

// Form returns form parameter by key.
func (c *Context) Form(key string) string {
	req := c.Request
	if s := req.PostFormValue(key); len(s) > 0 {
		return s
	}
	if req.MultipartForm != nil {
		if values := req.MultipartForm.Value[key]; len(values) > 0 {
			return values[0]
		}
	}
	return ""
}

// Get retrieves data from the context.
func (c *Context) Get(key string) interface{} {
	return c.store[key]
}

// Set saves data in the context.
func (c *Context) Set(key string, val interface{}) {
	if c.store == nil {
		c.store = make(store)
	}
	c.store[key] = val
}

// Param returns form values by parameter
func (c *Context) Param(name string) string {
	for _, param := range c.params {
		if param.Key == name {
			return param.Value
		}
	}
	return ""
}

// RemoteIpAddress returns the remote ip address, it parses
// X-Real-IP and X-Forwarded-For in order to work properly with reverse-proxies such us: nginx or haproxy.
func (c *Context) RemoteIpAddress() string {
	if c.NapNap.ForwardRemoteIpAddress {
		remoteIpAddr := strings.TrimSpace(c.Request.Header.Get("X-Real-Ip"))
		if len(remoteIpAddr) > 0 {
			return remoteIpAddr
		}
		remoteIpAddr = c.Request.Header.Get("X-Forwarded-For")
		if index := strings.IndexByte(remoteIpAddr, ','); index >= 0 {
			remoteIpAddr = remoteIpAddr[0:index]
		}
		remoteIpAddr = strings.TrimSpace(remoteIpAddr)
		if len(remoteIpAddr) > 0 {
			return remoteIpAddr
		}
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(c.Request.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}

// ContentType returns the Content-Type header of the request.
func (c *Context) ContentType() string {
	return filterFlags(c.Request.Header.Get("Content-Type"))
}

// SetCookie allows us to create an cookie
func (c *Context) SetCookie(
	name string,
	value string,
	maxAge int,
	path string,
	domain string,
	secure bool,
	httpOnly bool,
) {
	if path == "" {
		path = "/"
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
}

// Cookie returns cookie value
func (c *Context) Cookie(name string) (string, error) {
	cookie, err := c.Request.Cookie(name)
	if err != nil {
		return "", err
	}
	val, _ := url.QueryUnescape(cookie.Value)
	return val, nil
}

// Status is a intelligent shortcut for c.Writer.WriteHeader(code)
func (c *Context) Status(code int) {
	c.Writer.WriteHeader(code)
}

// RespHeader is a intelligent shortcut for c.Writer.Header().Set(key, value)
// It writes a header in the response.
// If value == "", this method removes the header `c.Writer.Header().Del(key)`
func (c *Context) RespHeader(key, value string) {
	if len(value) == 0 {
		c.Writer.Header().Del(key)
	} else {
		c.Writer.Header().Set(key, value)
	}
}

// RequestHeader is a intelligent shortcut for c.Request.Header.Get(key)
func (c *Context) RequestHeader(key string) string {
	return c.Request.Header.Get(key)
}

func (c *Context) reset(req *http.Request, w http.ResponseWriter) {
	c.Request = req
	c.Writer = w
	c.store = nil
	c.query = nil
	c.params = nil
}
