package main

import (
	"./napnap"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func newHelloRouter() *napnap.Router {
	router := napnap.NewRouter()

	router.Get("/hello", func(c *napnap.Context) {
		c.String(200, "Hello World")
	})

	router.Get("/hello/:name", func(c *napnap.Context) {
		name := c.Param("name")
		c.JSON(200, name)
	})

	router.Post("/bind-json", func(c *napnap.Context) {
		var json Person
		err := c.BindJSON(&json)

		if err != nil {
			c.String(400, err.Error())
		}

		c.JSON(200, json)

	})

	router.Post("/hello", func(c *napnap.Context) {
		queryName := c.Query("get")
		formName := c.Form("post")

		println("form: " + formName)
		c.String(200, queryName+","+formName)
	})

	return router
}
