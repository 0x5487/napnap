# NapNap micro web framework

## Roadmap
We are planning to add those features in the future.
- logging middleware

We support the following features
- routing features (static, parameterized)
- custom middleware
- http/2 (https only)
- rendering
- json binding

## Start using it
1. Download and install it:

    ```sh
    $ go get github.com/jasonsoft/napnap
    ```

2. Import it in your code:

    ```go
    import "github.com/jasonsoft/napnap"
    ```

## Getting Started

#### hello world example
```go
package main

import (
	"github.com/jasonsoft/napnap"
)

func main() {
	router := napnap.NewRouter()

	router.Get("/hello-world", func(c *napnap.Context) {
		c.String(200, "Hello, World")
	})

	nap := napnap.New()
	nap.Use(router)
	nap.Run(":8080") //run on port 8080
}
```