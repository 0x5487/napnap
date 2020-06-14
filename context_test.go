package napnap

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContextRemoteIpAddress(t *testing.T) {
	c, _, _ := createTestContext()

	c.Request, _ = http.NewRequest("POST", "/", nil)

	c.Request.Header.Set("X-Real-IP", " 10.10.10.10  ")
	assert.Equal(t, "10.10.10.10", c.ClientIP())
	c.Request.Header.Del("X-Real-IP")

	c.Request.Header.Set("X-Forwarded-For", "  20.20.20.20, 30.30.30.30")
	assert.Equal(t, "20.20.20.20", c.ClientIP())
	c.Request.Header.Set("X-Forwarded-For", "30.30.30.30  ")
	assert.Equal(t, "30.30.30.30", c.ClientIP())

	c.Request.Header.Del("X-Forwarded-For")

	c.Request.RemoteAddr = "  40.40.40.40:42123 "
	assert.Equal(t, "40.40.40.40", c.ClientIP())
}

func TestContextContentType(t *testing.T) {
	c, _, _ := createTestContext()

	c.Request, _ = http.NewRequest("POST", "/", nil)
	c.Request.Header.Set("Content-Type", "application/json; charset=utf-8")

	assert.Equal(t, "application/json", c.ContentType())
}

func TestContextSetCookie(t *testing.T) {
	c, _, _ := createTestContext()

	c.SetCookie("user", "jason", 1, "/", "localhost", true, true)
	assert.Equal(t, "user=jason; Path=/; Domain=localhost; Max-Age=1; HttpOnly; Secure", c.Writer.Header().Get("Set-Cookie"))
}

func TestContextGetCookie(t *testing.T) {
	c, _, _ := createTestContext()

	c.Request, _ = http.NewRequest("GET", "/get", nil)
	c.Request.Header.Set("Cookie", "user=jason")
	cookie, _ := c.Cookie("user")
	assert.Equal(t, "jason", cookie)
}

func TestContextSetRespHeader(t *testing.T) {
	c, _, _ := createTestContext()
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

func TestContextRedirectWithAbsolutePath(t *testing.T) {
	c, w, _ := createTestContext()
	c.Request, _ = http.NewRequest("POST", "http://example.com", nil)
	c.Redirect(302, "http://google.com")

	assert.Equal(t, 302, w.Code)
	assert.Equal(t, w.Header().Get("Location"), "http://google.com")
}

func TestContextRedirectWithRelativePath(t *testing.T) {
	c, w, _ := createTestContext()
	c.Request, _ = http.NewRequest("POST", "http://example.com", nil)

	c.Redirect(301, "/path")
	assert.Equal(t, 301, w.Code)
	assert.Equal(t, w.Header().Get("Location"), "/path")
}

func TestContextSetGet(t *testing.T) {
	c, _, _ := createTestContext()
	c.Set("foo", "bar")

	value, err := c.Get("foo")
	assert.Equal(t, value, "bar")
	assert.True(t, err)

	value, err = c.Get("foo2")
	assert.Nil(t, value)
	assert.False(t, err)

	assert.Equal(t, c.MustGet("foo"), "bar")
	assert.Panics(t, func() { c.MustGet("no_exist") })
}

func TestContextSetGetValues(t *testing.T) {
	c, _, _ := createTestContext()
	c.Set("string", "this is a string")
	c.Set("int32", int32(-42))
	c.Set("int64", int64(42424242424242))
	c.Set("uint64", uint64(42))
	c.Set("float32", float32(4.2))
	c.Set("float64", 4.2)
	var a interface{} = 1
	c.Set("intInterface", a)

	assert.Exactly(t, c.MustGet("string").(string), "this is a string")
	assert.Exactly(t, c.MustGet("int32").(int32), int32(-42))
	assert.Exactly(t, c.MustGet("int64").(int64), int64(42424242424242))
	assert.Exactly(t, c.MustGet("uint64").(uint64), uint64(42))
	assert.Exactly(t, c.MustGet("float32").(float32), float32(4.2))
	assert.Exactly(t, c.MustGet("float64").(float64), 4.2)
	assert.Exactly(t, c.MustGet("intInterface").(int), 1)

}

// Tests that the response executes the templates
// and responds with Content-Type set to text/html
// func TestContextRenderHTML(t *testing.T) {
// 	c, w, nap := createTestContext()
// 	templ := template.Must(template.New("t").Parse(`Hello {{.name}}`))
// 	nap.SetTemplate(templ)

// 	data := make(map[string]string, 1)
// 	data["name"] = "NapNap"

// 	c.Render(201, "t", data)

// 	assert.Equal(t, 201, w.Code)
// 	assert.Equal(t, "Hello NapNap", w.Body.String())
// 	assert.Equal(t, "text/html; charset=utf-8", w.HeaderMap.Get("Content-Type"))
// }
