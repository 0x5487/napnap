package main

import (
	"napnap"
)

func renderMiddleware() napnap.MiddlewareFunc {

	return func(c *napnap.Context, next napnap.HandlerFunc) {
		println("rendering")

		c.Render(200, "basic", nil)
	}
}

func middleware1(c *napnap.Context, next napnap.HandlerFunc) {
	println("=======================New Request")
	c.Set("key", "abcd")
	println("logging1")
	//before...
	next(c)
	//after...
	println("logging3")
}

func middleware2(c *napnap.Context, next napnap.HandlerFunc) {
	println("logging2")
	key := c.Get("key").(string)
	println("your key: " + key)
}
