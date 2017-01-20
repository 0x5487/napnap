package napnap

import (
	gcontext "context"
)

const (
	contextKey = "napnap_contextKey"
)

func FromContext(ctx gcontext.Context) (*Context, bool) {
	val, ok := ctx.Value(contextKey).(*Context)
	if ok {
		return val, ok
	}
	return nil, false
}

func newGContext(ctx gcontext.Context, c *Context) gcontext.Context {
	return gcontext.WithValue(ctx, contextKey, c)
}
