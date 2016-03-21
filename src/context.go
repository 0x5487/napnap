package napnap

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

type (
	// Param is a single URL parameter, consisting of a key and a value.
	Param struct {
		Key   string
		Value string
	}

	Context struct {
		NapNap  *NapNap
		Request *http.Request
		Writer  http.ResponseWriter
		query   url.Values
		params  []Param
		store   store
	}
	store map[string]interface{}
)

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

func (c *Context) reset(req *http.Request, w http.ResponseWriter) {
	c.Request = req
	c.Writer = w
	c.store = nil
	c.query = nil
	c.params = nil
}
