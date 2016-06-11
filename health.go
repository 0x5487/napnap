package napnap

import (
	"strings"
)

type Health struct {
}

func NewHealth() *Health {
	return &Health{}
}

func (h *Health) Invoke(c *Context, next HandlerFunc) {
	if strings.EqualFold(c.Request.URL.Path, "/health") {
		c.String(200, "OK")
	} else {
		next(c)
	}
}
