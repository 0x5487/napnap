package main

import "github.com/jasonsoft/napnap"

func main() {
	router := napnap.NewRouter()

	router.Get("/json-rendering", func(c *napnap.Context) {
		var person struct {
			Name string `json: name`
			Age  int    `json: age`
		}

		person.Name = "napnap"
		person.Age = 18

		c.JSON(200, person)
	})

	nap := napnap.New()
	nap.Use(router)
	nap.Run(":8080") //run on port 8080
}
