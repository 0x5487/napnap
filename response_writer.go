package napnap

import (
	"bufio"
	"net"
	"net/http"
)

const (
	noWritten     = -1
	defaultStatus = 200
)

// ResponseWriter wraps the original http.ResponseWriter
type ResponseWriter interface {
	http.ResponseWriter
	ContentLength() int
	Status() int
	reset(writer http.ResponseWriter) ResponseWriter
}

type responseWriter struct {
	http.ResponseWriter
	committed     bool
	status        int
	contentLength int
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
	if !rw.committed {
		// The status will be StatusOK if WriteHeader has not been called yet
		rw.WriteHeader(http.StatusOK)
	}
	n, err := rw.ResponseWriter.Write(b)
	rw.contentLength += n
	return n, err
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	if rw.committed {
		_logger.debug("Headers were already written.")
		return
	}

	// Store the status code
	rw.status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
	rw.committed = true
}

// Hijack implements the http.Hijacker interface to allow a HTTP handler to
// take over the connection.
// See https://golang.org/pkg/net/http/#Hijacker
func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return rw.ResponseWriter.(http.Hijacker).Hijack()
}

func (rw *responseWriter) reset(writer http.ResponseWriter) ResponseWriter {
	rw.ResponseWriter = writer
	rw.contentLength = noWritten
	rw.status = defaultStatus
	rw.committed = false
	return rw
}
