package napnap

import "testing"

func TestBasicStaticRoute(t *testing.T) {
	router := NewRouter()

	router.Add(GET, "/hello/get", func(c *Context) error {
		println("GET Method")
		return nil
	})

	router.Add(POST, "/hello/post", func(c *Context) error {
		println("POST Method")
		return nil
	})

	router.Add(PUT, "/hello/put", func(c *Context) error {
		println("PUT Method")
		return nil
	})

	router.Add(DELETE, "/hello/delete", func(c *Context) error {
		println("DELETE Method")
		return nil
	})

	println("===== result =====")
	c := NewContext(nil, nil)
	h := router.Find(GET, "/hello/get", c)
	if h == nil {
		t.Error("handler can't be nil")
	}

	h = router.Find(POST, "/hello/post", c)
	if h == nil {
		t.Error("handler can't be nil")
	}

	h = router.Find(PUT, "/hello/put", c)
	if h == nil {
		t.Error("handler can't be nil")
	}

	h = router.Find(DELETE, "/hello/delete", c)
	if h == nil {
		t.Error("handler can't be nil")
	}

	h = router.Find(GET, "/hello/404", c)
	if h == nil {
		t.Error("handler can't be nil")
	}
}

func TestParameterRoute(t *testing.T) {
	router := NewRouter()

	router.Add(GET, "/users/:user/name", func(c *Context) error {
		name := c.Param("user")
		println("user: " + name)
		return nil
	})

	router.Add(GET, "/users/:first/angela", func(c *Context) error {
		first := c.Param("first")
		println("first: " + first)
		return nil
	})

	router.Add(GET, "/users/:user/phone/:num", func(c *Context) error {
		user := c.Param("user")
		name := c.Param("num")
		println("user: " + user)
		println("name: " + name)
		return nil
	})

	c := NewContext(nil, nil)
	h := router.Find(GET, "/users/jason/angela", c)

	if h == nil {
		t.Error("handler can't be nil")
	}
	h(c)
}
