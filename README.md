# NapNap micro web framework

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

#### Using GET, POST, PUT, PATCH, DELETE and OPTIONS
```go
package main

import (
	"github.com/jasonsoft/napnap"
)

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
```

#### Parameters in path

```go
package main

import (
	"github.com/jasonsoft/napnap"
)

func main() {
	router := napnap.NewRouter()

	router.Get("/users/:name", func(c *napnap.Context) {
		name := c.Param("name")
		c.String(200, "Hello, "+name)
	})

	nap := napnap.New()
	nap.Use(router)
	nap.Run(":8080") //run on port 8080
}
```

#### Get querystring value
```go
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
```

#### Get post form value (Multipart/Urlencoded Form)
```go
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
```

#### JSON binding

```go
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
```

#### JSON rendering

```go
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
```

## Roadmap
We are planning to add those features in the future.
- logging middleware

We support the following features
- routing features (static, parameterized)
- custom middleware
- http/2 (https only)
- rendering
- json binding