package main

import (
	"napnap"
)

func main() {    
	nap := napnap.New()

	m1 := newMiddleware1()
	m2 := newMiddleware2()
	helloRouter := newHelloRouter()

	nap.Use(m1)
	nap.Use(helloRouter)
	nap.Use(m2)
    
	nap.Run(":8080")
}
