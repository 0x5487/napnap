package napnap

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContextRemoteIpAddress(t *testing.T) {
	c, _, _ := CreateTestContext()
	c.NapNap.ForwardRemoteIpAddress = true

	c.Request, _ = http.NewRequest("POST", "/", nil)

	c.Request.Header.Set("X-Real-IP", " 10.10.10.10  ")
	c.Request.Header.Set("X-Forwarded-For", "  20.20.20.20, 30.30.30.30")
	c.Request.RemoteAddr = "  40.40.40.40:42123 "

	assert.Equal(t, "10.10.10.10", c.RemoteIpAddress())

	c.Request.Header.Del("X-Real-IP")
	assert.Equal(t, "20.20.20.20", c.RemoteIpAddress())

	c.Request.Header.Set("X-Forwarded-For", "30.30.30.30  ")
	assert.Equal(t, "30.30.30.30", c.RemoteIpAddress())

	c.Request.Header.Del("X-Forwarded-For")
	assert.Equal(t, "40.40.40.40", c.RemoteIpAddress())
}

func TestContextContentType(t *testing.T) {
	c, _, _ := CreateTestContext()

	c.Request, _ = http.NewRequest("POST", "/", nil)
	c.Request.Header.Set("Content-Type", "application/json; charset=utf-8")

	assert.Equal(t, "application/json", c.ContentType())
}

func TestContextSetCookie(t *testing.T) {
	c, _, _ := CreateTestContext()

	c.SetCookie("user", "jason", 1, "/", "localhost", true, true)
	assert.Equal(t, "user=jason; Path=/; Domain=localhost; Max-Age=1; HttpOnly; Secure", c.Writer.Header().Get("Set-Cookie"))
}

func TestContextGetCookie(t *testing.T) {
	c, _, _ := CreateTestContext()

	c.Request, _ = http.NewRequest("GET", "/get", nil)
	c.Request.Header.Set("Cookie", "user=jason")
	cookie, _ := c.Cookie("user")
	assert.Equal(t, "jason", cookie)
}

func TestContextSetRespHeader(t *testing.T) {
	c, _, _ := CreateTestContext()
	c.RespHeader("Content-Type", "text/plain")
	c.RespHeader("X-Custom", "value")

	assert.Equal(t, c.Writer.Header().Get("Content-Type"), "text/plain")
	assert.Equal(t, c.Writer.Header().Get("X-Custom"), "value")

	c.RespHeader("Content-Type", "text/html")
	c.RespHeader("X-Custom", "")

	assert.Equal(t, c.Writer.Header().Get("Content-Type"), "text/html")
	_, exist := c.Writer.Header()["X-Custom"]
	assert.False(t, exist)
}
