package napnap

import (
	"net/http/httptest"
)

func CreateTestContext() (*Context, ResponseWriter, *NapNap) {
	nap := New()
	w := NewResponseWriter(httptest.NewRecorder())
	c := NewContext(nap, nil, w)
	return c, w, nap
}
