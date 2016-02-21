package napnap

import "testing"

func TestBasicStaticRoute(t *testing.T) {
	router := NewRouter()

	router.Add(GET, "/hello/get", func() error {
		println("GET Method")
		return nil
	})

	router.Add(POST, "/hello/post", func() error {
		println("POST Method")
		return nil
	})

	router.Add(PUT, "/hello/put", func() error {
		println("PUT Method")
		return nil
	})

	router.Add(DELETE, "/hello/delete", func() error {
		println("DELETE Method")
		return nil
	})

	println("===== result =====")
	h := router.Find(GET, "/hello/get")
	if h == nil {
		t.Error("handler can't be nil")
	}
	h()

	h = router.Find(POST, "/hello/post")
	if h == nil {
		t.Error("handler can't be nil")
	}
	h()

	h = router.Find(PUT, "/hello/put")
	if h == nil {
		t.Error("handler can't be nil")
	}
	h()

	h = router.Find(DELETE, "/hello/delete")
	if h == nil {
		t.Error("handler can't be nil")
	}
	h()

	h = router.Find(GET, "/hello/404")
	if h == nil {
		t.Error("handler can't be nil")
	}
	h()
}
