package main

import (
	"./napnap"
)

func main() {
	nap := napnap.New()

	nap.Get("/hello-world", func(c *napnap.Context) error {
		c.String(200, "Hello World")
		return nil
	})

	nap.Get("/hello-meican", func(c *napnap.Context) error {
		c.JSON(200, "Hello Meican")
		return nil
	})

	nap.Post("/hello", func(c *napnap.Context) error {
		name := c.Query("name")
		nickName := c.Form("nick_name")

		c.String(200, name+","+nickName)
		return nil
	})

	nap.Run(":8080")
}
