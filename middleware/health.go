package middleware

import (
	"strings"

	"github.com/jasonsoft/napnap"
)

type Health struct {
}

func NewHealth() *Health {
	return &Health{}
}

func (h *Health) Invoke(c *napnap.Context, next napnap.HandlerFunc) {
	if strings.EqualFold(c.Request.URL.Path, "/health") {
		c.String(200, "OK")
	} else {
		next(c)
	}
}
