package napnap

import (
	gcontext "context"
)

var (
	ctxKey = &struct {
		name string
	}{
		name: "napnap",
	}
)

// FromContext return a napnap context from the standard context
func FromContext(ctx gcontext.Context) (*Context, bool) {
	val, ok := ctx.Value(ctxKey).(*Context)
	if ok {
		return val, ok
	}
	return nil, false
}

func newGContext(ctx gcontext.Context, c *Context) gcontext.Context {
	return gcontext.WithValue(ctx, ctxKey, c)
}
