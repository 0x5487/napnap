package napnap

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type (
	Context struct {
		Request *http.Request
		Writer  http.ResponseWriter
		query   url.Values
	}
)

func NewContext(req *http.Request, writer http.ResponseWriter) *Context {
	return &Context{
		Request: req,
		Writer:  writer,
	}
}

// respnse a string
func (c *Context) String(code int, s string) (err error) {
	c.Writer.WriteHeader(code)
	c.Writer.Write([]byte(s))
	return nil
}

// response JSON format
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
	s := c.Request.FormValue(name)
	return s
}
