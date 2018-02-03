package napnap

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testRoute(t *testing.T, method string, path string) {
	passed := false
	_, w, nap := CreateTestContext()

	router := NewRouter()
	router.Add(method, path, func(c *Context) {
		passed = true
		c.SetStatus(200)
	})
	nap.Use(router)

	req, _ := http.NewRequest(method, path, nil)
	nap.ServeHTTP(w, req)

	assert.True(t, passed)
	assert.Equal(t, 200, w.Code)
}

func TestRouterStaticRoute(t *testing.T) {
	testRoute(t, "GET", "/")
	testRoute(t, "GET", "/hello")
	testRoute(t, "POST", "/hello")
	testRoute(t, "PUT", "/hello/put")
	testRoute(t, "DELETE", "/hello/Delet")
}

func TestRouterParameterRoute(t *testing.T) {
	var name string
	_, w, nap := CreateTestContext()

	router := NewRouter()
	router.Add(GET, "/users/:name", func(c *Context) {
		name = c.Param("name")
		c.SetStatus(200)
	})
	nap.Use(router)

	req, _ := http.NewRequest("GET", "/users/john", nil)
	nap.ServeHTTP(w, req)

	assert.Equal(t, "john", name)
	assert.Equal(t, 200, w.Code)
}

func TestRouterMatchAnyRoute(t *testing.T) {
	var action string
	_, w, nap := CreateTestContext()

	router := NewRouter()
	router.Add(GET, "/video/:action1", func(c *Context) {
		fmt.Print("action1")
		action = c.Param("action1")
		c.SetStatus(201)
	})

	router.Add(GET, "/images/*action2", func(c *Context) {
		fmt.Print("action2")
		action = c.Param("action2")
		c.SetStatus(200)
	})
	nap.Use(router)

	req, _ := http.NewRequest("GET", "/images/play/ball.jpg", nil)
	nap.ServeHTTP(w, req)

	assert.Equal(t, "play/ball.jpg", action)
	assert.Equal(t, 200, w.Code)
}
