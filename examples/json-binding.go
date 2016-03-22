package main

import "github.com/jasonsoft/napnap"

func main() {
	router := napnap.NewRouter()

	router.Post("/json-binding", func(c *napnap.Context) {
		var person struct {
			Name string `json: name`
			Age  int    `json: age`
		}
		if c.BindJSON(&person) == nil {
			c.String(200, person.Name)
		}
	})

	nap := napnap.New()
	nap.Use(router)
	nap.Run(":8080") //run on port 8080
}
