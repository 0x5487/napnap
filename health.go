package napnap

type Health struct {
}

func NewHealth() *Health {
	return &Health{}
}

func (h *Health) Invoke(c *Context, next HandlerFunc) {
	if c.Request.URL.Path == "/health" {
		c.String(200, "ok")
	} else {
		next(c)
	}
}
