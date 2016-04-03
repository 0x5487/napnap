package napnap

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoteIpAddress(t *testing.T) {
	nap := New()
	nap.ForwardRemoteIpAddress = true
	c := NewContext(nap, nil, nil)

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
