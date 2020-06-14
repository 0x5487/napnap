package middleware

import (
	"strings"

	"github.com/jasonsoft/napnap"
)

// Health is health middleware struct
type Health struct {
}

// NewHealth returns Health middlware instance
func NewHealth() *Health {
	return &Health{}
}

// Invoke function is a middleware entry
func (h *Health) Invoke(c *napnap.Context, next napnap.HandlerFunc) {
	if strings.EqualFold(c.Request.URL.Path, "/health") {
		_ = c.String(200, "OK")
	} else {
		_ = next(c)
	}
}
