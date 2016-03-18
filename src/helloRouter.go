package main

import (
	"napnap"
)

func newHelloRouter() *napnap.Router {
	router := napnap.NewRouter()

	router.Get("/hello", func(c *napnap.Context) {
		c.String(200, "Hello World")
	})

	router.Get("/hello/:name", func(c *napnap.Context) {
		name := c.Param("name")
		c.JSON(200, name)
	})

	router.Post("/hello", func(c *napnap.Context) {
		queryName := c.Query("first_name")
		formName := c.Form("first_name")

		println("form: " + formName)
		c.String(200, queryName+","+formName)
	})

	return router
}
