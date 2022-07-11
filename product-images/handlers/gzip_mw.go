package handlers

import (
	"compress/gzip"
	"net/http"
	"strings"
)

// GzipHandler type
type GzipHandler struct {
}

// GzipMiddleware is a middleware used in GetRouter that comes into play when req header has gzip as Accept-Encoding
func (g *GzipHandler) GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// check if gzip/--compressed
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// create a gziped response
			// rw.Write([]byte("works"))
			wrw := NewWrappedResponseWriter(rw)
			wrw.Header().Set("Content-Encoding", "gzip")
			next.ServeHTTP(wrw, r)
			defer wrw.Flush()
			return
		}
		// handle normal
		next.ServeHTTP(rw, r)
	})
}

// Our Wrapped ResponseWriter which has gzip.Writer embedded to do gzip things
type WrappedResponseWriter struct {
	rw http.ResponseWriter
	gw *gzip.Writer
}

// NewWrappedResponseWriter returns a WrappedResponseWriter instance with given ResponseWriter, and adds gzip instance as well
func NewWrappedResponseWriter(rw http.ResponseWriter) *WrappedResponseWriter {
	gw := gzip.NewWriter(rw)
	return &WrappedResponseWriter{rw: rw, gw: gw}
}

// Header - same as http.ResponseWriter
func (wr *WrappedResponseWriter) Header() http.Header {
	return wr.rw.Header()
}

// Write - overridden to use gzip write
func (wr *WrappedResponseWriter) Write(d []byte) (int, error) {
	return wr.gw.Write(d)
}

// WriteHeader - same as http.ResponseWriter
func (wr *WrappedResponseWriter) WriteHeader(statusCode int) {
	wr.rw.WriteHeader(statusCode)
}

// Flush - cleanup
func (wr *WrappedResponseWriter) Flush() {
	wr.gw.Flush()
	wr.gw.Close()
}
