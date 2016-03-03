package napnap

import (
	"encoding/json"
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
		Request *http.Request
		Writer  http.ResponseWriter
		query   url.Values
		params  []Param
		store   store
	}
	store map[string]interface{}
)

// NewContext returns a new context instance
func NewContext(req *http.Request, writer http.ResponseWriter) *Context {
	return &Context{
		Request: req,
		Writer:  writer,
	}
}


// String returns string format
func (c *Context) String(code int, s string) (err error) {
	c.Writer.WriteHeader(code)
	c.Writer.Write([]byte(s))
	return nil
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
	return nil
}

// Query returns query parameter by name.
func (c *Context) Query(name string) string {
	if c.query == nil {
		c.query = c.Request.URL.Query()
	}
	return c.query.Get(name)
}

// Form returns form parameter by name.
func (c *Context) Form(name string) string {
	s := c.Request.FormValue(name) // bug, it will also get value from querystring as well
	return s
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