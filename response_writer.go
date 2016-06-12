package napnap

import "net/http"

const (
	noWritten     = -1
	defaultStatus = 200
)

type beforeFunc func(ResponseWriter)

// ResponseWriter wraps the original http.ResponseWriter
type ResponseWriter interface {
	http.ResponseWriter
	ContentLength() int
	Status() int
	Written() bool
	// Before allows for a function to be called before the ResponseWriter has been written to. This is
	// useful for setting headers or any other operations that must happen before a response has been written.
	Before(func(ResponseWriter))
	reset(writer http.ResponseWriter) ResponseWriter
}

type responseWriter struct {
	http.ResponseWriter
	status        int
	contentLength int
	beforeFuncs   []beforeFunc
}

// NewResponseWriter returns a ResponseWriter which wraps the writer
func NewResponseWriter() ResponseWriter {
	return &responseWriter{
		status:        defaultStatus,
		contentLength: noWritten,
	}
}

// ContentLength returns size of content length
func (rw *responseWriter) ContentLength() int {
	return rw.contentLength
}

// Status returns http status code
func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.Written() {
		// The status will be StatusOK if WriteHeader has not been called yet
		rw.WriteHeader(http.StatusOK)
	}
	n, err := rw.ResponseWriter.Write(b)
	rw.contentLength += n
	return n, err
}

func (rw *responseWriter) Written() bool {
	return rw.status != 0
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	// Store the status code
	rw.status = statusCode
	rw.callBefore()
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriter) Before(before func(ResponseWriter)) {
	rw.beforeFuncs = append(rw.beforeFuncs, before)
}

func (rw *responseWriter) callBefore() {
	for i := len(rw.beforeFuncs) - 1; i >= 0; i-- {
		rw.beforeFuncs[i](rw)
	}
}

func (rw *responseWriter) reset(writer http.ResponseWriter) ResponseWriter {
	rw.ResponseWriter = writer
	rw.contentLength = noWritten
	rw.status = defaultStatus
	return rw
}
