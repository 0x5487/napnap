package napnap

import (
	"net/http/httptest"
)

func CreateTestContext() (c *Context, w *httptest.ResponseRecorder, nap *NapNap) {
	nap = New()
	w = httptest.NewRecorder()
	c = NewContext(nap, nil, w)
	return
}
