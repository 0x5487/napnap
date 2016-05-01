package napnap

import "net/http"

const (
	noWritten     = -1
	defaultStatus = 200
)

// ResponseWriter wraps the original http.ResponseWriter
type ResponseWriter interface {
	http.ResponseWriter
	ContentLength() int
	Status() int
	reset(writer http.ResponseWriter)
}

type responseWriter struct {
	http.ResponseWriter
	status        int
	contentLength int
}

// NewResponseWriter returns a ResponseWriter which wraps the writer
func NewResponseWriter(writer http.ResponseWriter) ResponseWriter {
	return &responseWriter{
		ResponseWriter: writer,
		status:         defaultStatus,
		contentLength:  noWritten,
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
	n, err := rw.ResponseWriter.Write(b)
	rw.contentLength += n
	return n, err
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	// Store the status code
	rw.status = statusCode
	rw.ResponseWriter.WriteHeader(rw.status)
}

func (rw *responseWriter) reset(writer http.ResponseWriter) {
	rw.ResponseWriter = writer
	rw.contentLength = noWritten
	rw.status = defaultStatus
}
