package middleware

import (
	"net/http/pprof"

	"github.com/jasonsoft/napnap"
)

type PPROF struct {
}

func NewPPROF() *PPROF {
	return &PPROF{}
}

func (p *PPROF) Invoke(c *napnap.Context, next napnap.HandlerFunc) {
	pprof.Index(c.Writer, c.Request)
	pprof.Cmdline(c.Writer, c.Request)
	pprof.Profile(c.Writer, c.Request)
	pprof.Symbol(c.Writer, c.Request)
	pprof.Trace(c.Writer, c.Request)
	next(c)
}
