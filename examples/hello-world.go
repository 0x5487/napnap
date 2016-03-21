package main

import (
	"github.com/jasonsoft/napnap"
)

func main() {
	router := napnap.NewRouter()

	router.Get("/hello-world", func(c *napnap.Context) {
		c.String(200, "Hello, World")
	})

	nap := napnap.New()
	nap.Use(router)
	nap.Run(":8080") //run on port 8080
}
