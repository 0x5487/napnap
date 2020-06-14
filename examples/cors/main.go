package main

import (
	"net/http"

	"github.com/jasonsoft/napnap"
	"github.com/jasonsoft/napnap/middleware"
)

func main() {
	nap := napnap.New()

	options := middleware.Options{}
	options.AllowedMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
	options.AllowedHeaders = []string{"*", "Authorization", "Content-Type", "Origin", "Content-Length"}
	nap.Use(middleware.NewCors(options))
	nap.Use(middleware.NewHealth())

	nap.Get("/", func(c *napnap.Context) error {
		return c.String(200, "Hello World")
	})

	http.ListenAndServe("127.0.0.1:10080", nap)
}
