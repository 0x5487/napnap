package napnap

import (
	"net/http/httptest"
)

func CreateTestContext() (*Context, ResponseWriter, *NapNap) {
	nap := New()
	rw := NewResponseWriter()
	rw.reset(httptest.NewRecorder())
	c := NewContext(nap, nil, rw)
	return c, rw, nap
}
