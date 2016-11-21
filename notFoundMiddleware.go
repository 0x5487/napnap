package napnap

type NotfoundMiddleware struct {
}

func NewNotfoundMiddleware() *NotfoundMiddleware {
	return &NotfoundMiddleware{}
}

func (m *NotfoundMiddleware) Invoke(c *Context, next HandlerFunc) {
	c.SetStatus(404)
}
