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
		query_name := c.Query("first_name")
		form_name := c.Form("last_name")
		println("form: " + form_name)
		c.String(200, query_name+","+form_name)
	})

	return router
}
