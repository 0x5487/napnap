package main

import (
	"napnap"
)



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
