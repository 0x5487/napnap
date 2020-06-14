package middleware

import (
	"net/http/pprof"

	"github.com/jasonsoft/napnap"
)

// PPROF is middleware struct
type PPROF struct {
}

// NewPPROF returns a mddleware instance
func NewPPROF() *PPROF {
	return &PPROF{}
}

// Invoke function is a middleware entry
func (p *PPROF) Invoke(c *napnap.Context, next napnap.HandlerFunc) {
	pprof.Index(c.Writer, c.Request)
	pprof.Cmdline(c.Writer, c.Request)
	pprof.Profile(c.Writer, c.Request)
	pprof.Symbol(c.Writer, c.Request)
	pprof.Trace(c.Writer, c.Request)
	_ = next(c)
}
