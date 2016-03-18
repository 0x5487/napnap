package main

import (
	"napnap"
)

func main() {    
	nap := napnap.New()

	helloRouter := newHelloRouter()

	nap.UseFunc(middleware1)
	nap.Use(helloRouter)
	nap.UseFunc(middleware2)
    
	nap.Run(":8080")
}
