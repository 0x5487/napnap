package main

import (
	"github.com/jasonsoft/napnap"
)

func main() {
	router := napnap.NewRouter()

	router.Get("/querystring-value", func(c *napnap.Context) {
		page := c.Query("page") //get query string value
		c.String(200, page)
	})

	nap := napnap.New()
	nap.Use(router)
	nap.Run(":8080") //run on port 8080
}
