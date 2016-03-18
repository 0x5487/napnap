package main

import (
	"napnap"
)

func main() {
	nap := napnap.New()

	nap.SetViews("views/*")

	//nap.UseFunc(renderMiddleware())
	nap.UseFunc(middleware1)

	helloRouter := newHelloRouter()
	nap.Use(helloRouter)

	nap.UseFunc(middleware2)

	nap.Run(":8080")
}
