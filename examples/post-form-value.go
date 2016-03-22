package main

import (
	"github.com/jasonsoft/napnap"
)

func main() {
	router := napnap.NewRouter()

	router.Post("/post-form-value", func(c *napnap.Context) {
		userId := c.Form("user_id") //get post form value
		c.String(200, userId)
	})

	nap := napnap.New()
	nap.Use(router)
	nap.Run(":8080") //run on port 8080
}
