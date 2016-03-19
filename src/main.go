package main

import "napnap"

func main() {
	nap := napnap.New()

	nap.SetViews("views/*")

	nap.Use(napnap.NewStatic("public"))
	//nap.UseFunc(renderMiddleware())
	nap.UseFunc(middleware1)
	nap.Use(newHelloRouter())
	nap.UseFunc(middleware2)
	nap.UseFunc(notFound)

	nap.Run(":8080")
}
