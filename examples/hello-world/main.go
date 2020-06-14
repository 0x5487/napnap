package main

import (
	"net/http"

	"github.com/jasonsoft/napnap"
)

func main() {
	nap := napnap.New()

	nap.Get("/", func(c *napnap.Context) error {
		return c.String(200, "Hello World")
	})

	_ = http.ListenAndServe("127.0.0.1:10080", nap)
}
