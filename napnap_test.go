package napnap

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultHandlers(t *testing.T) {
	_, w, nap := CreateTestContext()

	isError := false
	nap.ErrorHandler = func(c *Context, err error) {
		isError = true
		assert.Equal(t, "oops", err.Error())
	}

	isNotFound := false
	nap.NotFoundHandler = func(c *Context) error {
		isNotFound = true
		return nil
	}

	nap.Get("/error", func(c *Context) error {
		return errors.New("oops")
	})

	req, _ := http.NewRequest("GET", "/error", nil)
	nap.ServeHTTP(w, req)

	req, _ = http.NewRequest("GET", "/not_found", nil)
	nap.ServeHTTP(w, req)

	assert.Equal(t, true, isError)
	assert.Equal(t, true, isNotFound)
}

func TestMidderwareOrder(t *testing.T) {

}
