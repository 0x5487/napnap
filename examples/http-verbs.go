package main

import (
	"github.com/jasonsoft/napnap"
)

func myGetEndpoint(c *napnap.Context) {
	c.String(200, "my Get")
}

func myPostEndpoint(c *napnap.Context) {
	c.String(200, "my post")
}

func myPutEndpoint(c *napnap.Context) {
	c.String(200, "my put")
}

func myDeleteEndpoint(c *napnap.Context) {
	c.String(200, "my delete")
}

func myPatchEndpoint(c *napnap.Context) {
	c.String(200, "my patch")
}

func myOptionsEndpoint(c *napnap.Context) {
	c.String(200, "my options")
}

func myHeadEndpoint(c *napnap.Context) {
	w := c.Writer
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func main() {
	router := napnap.NewRouter()

	router.Get("/my-get", myHeadEndpoint)
	router.Post("/my-post", myPostEndpoint)
	router.Put("/my-put", myPutEndpoint)
	router.Delete("/my-delete", myDeleteEndpoint)
	router.Patch("/my-patch", myPatchEndpoint)
	router.Options("/my-options", myOptionsEndpoint)
	router.Head("/my-head", myHeadEndpoint)

	nap := napnap.New()
	nap.Use(router)
	nap.Run(":8080") //run on port 8080
}
