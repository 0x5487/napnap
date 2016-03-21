package main

import (
	"github.com/jasonsoft/napnap"
)

func main() {
	router := napnap.NewRouter()

	router.Get("/users/:name", func(c *napnap.Context) {
		name := c.Param("name")
		c.String(200, "Hello, "+name)
	})

	nap := napnap.New()
	nap.Use(router)
	nap.Run(":8080") //run on port 8080
}
