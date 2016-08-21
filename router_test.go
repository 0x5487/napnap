package napnap

import (
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
