package napnap

import (
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Param is a single URL parameter, consisting of a key and a value.
type Param struct {
	Key   string
	Value string
}

type Context struct {
	NapNap  *NapNap
	Request *http.Request
	Writer  ResponseWriter
	query   url.Values
	params  []Param
	store   map[string]interface{}
}

// NewContext returns a new context instance
func NewContext(napnap *NapNap, req *http.Request, writer ResponseWriter) *Context {
	return &Context{
		NapNap:  napnap,
		Request: req,
		Writer:  writer,
	}
}

// Render returns html format
func (c *Context) Render(code int, viewName string, data interface{}) (err error) {
	c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
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

// Redirect returns a HTTP redirect to the specific location.
func (c *Context) Redirect(code int, location string) error {
	if (code < 300 || code > 308) && code != 201 {
		return fmt.Errorf("Cannot redirect with status code %d", code)
	}
	http.Redirect(c.Writer, c.Request, location, code)
	return nil
}

// BindJSON binds the request body into provided type `obj`. The default binder does
// it based on Content-Type header.
func (c *Context) BindJSON(obj interface{}) error {
	req := c.Request
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(obj)
	if err != nil {
		return err
	}
	return nil
}

// Query returns query parameter by key.
func (c *Context) Query(key string) string {
	if c.query == nil {
		c.query = c.Request.URL.Query()
	}
	return c.query.Get(key)
}

// QueryInt returns query parameter by key and cast the value to int.
func (c *Context) QueryInt(key string) (int, error) {
	return strconv.Atoi(c.Query(key))
}

// QueryIntWithDefault returns query parameter by key and cast the value to int.  If the value doesn't exist, the default value will be used.
func (c *Context) QueryIntWithDefault(key string, defaultValue int) (int, error) {
	data := c.Query(key)
	if len(data) > 0 {
		return strconv.Atoi(c.Query(key))
	}
	return defaultValue, nil
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

// FormFile returns file.
func (c *Context) FormFile(key string) (*multipart.FileHeader, error) {
	_, fh, err := c.Request.FormFile(key)
	return fh, err
}

// Get retrieves data from the context.
func (c *Context) Get(key string) (interface{}, bool) {
	var value interface{}
	var exists bool
	if c.store != nil {
		value, exists = c.store[key]
	}
	return value, exists
}

// MustGet returns the value for the given key if it exists, otherwise it panics.
func (c *Context) MustGet(key string) interface{} {
	if value, exists := c.Get(key); exists {
		return value
	}
	panic("Key \"" + key + "\" does not exist")
}

// Set saves data in the context.
// It also lazy initializes  c.Keys if it was not used previously.
func (c *Context) Set(key string, val interface{}) {
	if c.store == nil {
		c.store = make(map[string]interface{})
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

// ParamInt returns parameter by key and cast the value to int.
func (c *Context) ParamInt(key string) (int, error) {
	return strconv.Atoi(c.Param(key))
}

// RemoteIPAddress returns the remote ip address, it parses
// X-Real-IP and X-Forwarded-For in order to work properly with reverse-proxies such us: nginx or haproxy.
func (c *Context) RemoteIPAddress() string {
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

// SetStatus is a intelligent shortcut for c.Writer.WriteHeader(code)
func (c *Context) SetStatus(code int) {
	c.Writer.WriteHeader(code)
}

// Status is a intelligent shortcut for c.Writer.Status()
func (c *Context) Status() int {
	return c.Writer.Status()
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

// StdContext return golang standard context
func (c *Context) StdContext() context.Context {
	ctx := c.Request.Context()
	ctx = newGContext(ctx, c)
	return ctx
}

// SetStdContext allow us to save the golang context to request
func (c *Context) SetStdContext(ctx context.Context) {
	c.Request = c.Request.WithContext(ctx)
}

// DeviceType returns user's device type which includes web, mobile, tab, tv
func (c *Context) DeviceType() string {
	userAgent := c.RequestHeader("User-Agent")
	deviceType := "web"

	if strings.Contains(userAgent, "Android") ||
		strings.Contains(userAgent, "webOS") ||
		strings.Contains(userAgent, "iPhone") ||
		strings.Contains(userAgent, "BlackBerry") ||
		strings.Contains(userAgent, "Windows Phone") {
		deviceType = "mobile"
	} else if strings.Contains(userAgent, "iPad") ||
		strings.Contains(userAgent, "iPod") ||
		(strings.Contains(userAgent, "tablet") ||
			strings.Contains(userAgent, "RX-34") ||
			strings.Contains(userAgent, "FOLIO")) ||
		(strings.Contains(userAgent, "Kindle") ||
			strings.Contains(userAgent, "Mac OS") &&
				strings.Contains(userAgent, "Silk")) ||
		(strings.Contains(userAgent, "AppleWebKit") &&
			strings.Contains(userAgent, "Silk")) {
		deviceType = "tab"
	} else if strings.Contains(userAgent, "TV") ||
		strings.Contains(userAgent, "NetCast") ||
		strings.Contains(userAgent, "boxee") ||
		strings.Contains(userAgent, "Kylo") ||
		strings.Contains(userAgent, "Roku") ||
		strings.Contains(userAgent, "DLNADOC") {
		deviceType = "tv"
	}
	return deviceType
}

func (c *Context) reset(w http.ResponseWriter, req *http.Request) {
	c.Request = req
	c.Writer = c.Writer.reset(w)
	c.store = nil
	c.query = nil
	c.params = nil
}
