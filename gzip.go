package napnap

import (
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

// These compression constants are copied from the compress/gzip package.
const (
	encodingGzip = "gzip"

	headerAcceptEncoding  = "Accept-Encoding"
	headerContentEncoding = "Content-Encoding"
	headerContentLength   = "Content-Length"
	headerContentType     = "Content-Type"
	headerVary            = "Vary"
	headerSecWebSocketKey = "Sec-WebSocket-Key"

	BestCompression    = gzip.BestCompression
	BestSpeed          = gzip.BestSpeed
	DefaultCompression = gzip.DefaultCompression
	NoCompression      = gzip.NoCompression
)

// gzipResponseWriter is the ResponseWriter that http.ResponseWriter is
// wrapped in.
type gzipResponseWriter struct {
	gz        *gzip.Writer
	napWriter ResponseWriter
	ResponseWriter
}

// Write writes bytes to the gzip.Writer. It will also set the Content-Type
// header using the net/http library content type detection if the Content-Type
// header was not set yet.
func (grw gzipResponseWriter) Write(b []byte) (int, error) {
	if len(grw.Header().Get(headerContentType)) == 0 {
		grw.Header().Set(headerContentType, http.DetectContentType(b))
	}
	if len(grw.Header().Get(headerContentEncoding)) > 0 {
		// compress the content
		return grw.gz.Write(b)
	}
	// no compress
	grw.gz.Reset(ioutil.Discard)
	return grw.napWriter.Write(b)
}

// handler struct contains the ServeHTTP method
type gzipMiddleware struct {
	pool sync.Pool
}

// NewGzip returns a middleware which will handle the Gzip compression in Invoke.
// Valid values for level are identical to those in the compress/gzip package.
func NewGzip(level int) *gzipMiddleware {
	h := &gzipMiddleware{}
	h.pool.New = func() interface{} {
		gz, err := gzip.NewWriterLevel(ioutil.Discard, level)
		if err != nil {
			panic(err)
		}
		return gz
	}
	return h
}

// Invoke wraps the http.ResponseWriter with a gzip.Writer.
func (h *gzipMiddleware) Invoke(c *Context, next HandlerFunc) {
	r := c.Request
	w := c.Writer
	// Skip compression if the client doesn't accept gzip encoding.
	if !strings.Contains(r.Header.Get(headerAcceptEncoding), encodingGzip) {
		next(c)
		return
	}

	// Skip compression if client attempt WebSocket connection
	if len(r.Header.Get(headerSecWebSocketKey)) > 0 {
		next(c)
		return
	}

	// Skip compression if already compressed
	if w.Header().Get(headerContentEncoding) == encodingGzip {
		next(c)
		return
	}

	// Retrieve gzip writer from the pool. Reset it to use the ResponseWriter.
	// This allows us to re-use an already allocated buffer rather than
	// allocating a new buffer for every request.
	// We defer g.pool.Put here so that the gz writer is returned to the
	// pool if any thing after here fails for some reason (functions in
	// next could potentially panic, etc)
	gz := h.pool.Get().(*gzip.Writer)
	defer h.pool.Put(gz)
	gz.Reset(w)

	// Set the appropriate gzip headers.
	headers := w.Header()
	headers.Set(headerContentEncoding, encodingGzip)
	headers.Set(headerVary, headerAcceptEncoding)

	// Wrap the original http.ResponseWriter
	// and create the gzipResponseWriter.
	grw := gzipResponseWriter{
		gz,
		w,
		w,
	}

	// Call the next handler supplying the gzipResponseWriter instead of
	// the original.
	c.Writer = grw
	next(c)

	gz.Close()
}
