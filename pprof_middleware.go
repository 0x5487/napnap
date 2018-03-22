package napnap

import (
	"net/http/pprof"
)

type PPROF struct {
}

func NewPPROF() *PPROF {
	return &PPROF{}
}

func (p *PPROF) Invoke(c *Context, next HandlerFunc) {
	pprof.Index(c.Writer, c.Request)
	pprof.Cmdline(c.Writer, c.Request)
	pprof.Profile(c.Writer, c.Request)
	pprof.Symbol(c.Writer, c.Request)
	pprof.Trace(c.Writer, c.Request)
	next(c)
}
