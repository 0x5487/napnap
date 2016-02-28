package main

import (
	"napnap"
)

type middleware1 struct {
}

func newMiddleware1() middleware1 {
	return middleware1{}
}

func (m middleware1) Execute(c *napnap.Context, next napnap.HandlerFunc) {
	println("=======================New Request")
	c.Set("key", "abcd")
	println("logging1")
	//before...
	next(c)
	//after...
	println("logging3")
}

type middleware2 struct {
}

func newMiddleware2() middleware2 {
	return middleware2{}
}

func (m middleware2) Execute(c *napnap.Context, next napnap.HandlerFunc) {
	println("logging2")
	key := c.Get("key").(string)
	println("your key: " + key)
}
