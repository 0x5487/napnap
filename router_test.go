package napnap

import "testing"

func TestBasicStaticRoute(t *testing.T) {
	router := NewRouter()

	router.Add(GET, "/hello/get", func(c *Context) {
		println("GET Method")
	})

	router.Add(POST, "/hello/post", func(c *Context) {
		println("POST Method")
	})

	router.Add(PUT, "/hello/put", func(c *Context) {
		println("PUT Method")
	})

	router.Add(DELETE, "/hello/delete", func(c *Context) {
		println("DELETE Method")
	})

	println("===== result =====")
	nap := New()
	c := NewContext(nap, nil, nil)
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

	router.Add(GET, "/users/:user/name", func(c *Context) {
		name := c.Param("user")
		println("user: " + name)
	})

	router.Add(GET, "/users/:first/angela", func(c *Context) {
		first := c.Param("first")
		println("first: " + first)
	})

	router.Add(GET, "/users/:user/phone/:num", func(c *Context) {
		user := c.Param("user")
		name := c.Param("num")
		println("user: " + user)
		println("name: " + name)
	})
	nap := New()
	c := NewContext(nap, nil, nil)
	h := router.Find(GET, "/users/jason/angela", c)

	if h == nil {
		t.Error("handler can't be nil")
	}
	h(c)
}
