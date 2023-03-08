package server

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/indiependente/pkg/logger"
)

// requestIDHeader is the request ID HTTP header.
const requestIDHeader = "X-Request-Id"

type logEntry struct {
	reqID  string
	ip     string
	ua     string
	host   string
	method string
	uri    string
}

// HTTP middleware setting a value on the request context
func FastLogger(log logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			r, entry := newLogEntry(r) // override incoming request
			ww := wrapResponseWriter(w, time.Now())

			defer func() {
				log.BytesWritten(ww.writtenBytes).
					Duration(ww.duration).
					Host(entry.host).
					Method(entry.method).
					RemoteAddr(entry.ip).
					RequestID(entry.reqID).
					StatusCode(ww.status).
					URI(entry.uri).
					UserAgent(entry.ua).
					Info("request processed")
			}()

			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}

func newLogEntry(r *http.Request) (*http.Request, *logEntry) {
	req, reqID := getRequestID(r)

	return req, &logEntry{
		reqID:  reqID,
		ip:     r.RemoteAddr,
		ua:     r.UserAgent(),
		host:   r.Host,
		uri:    r.RequestURI,
		method: r.Method,
	}
}

func getRequestID(r *http.Request) (*http.Request, string) {
	ctx := r.Context()
	if ctx == nil {
		reqID := r.Header.Get(requestIDHeader)
		reqIDCtx := context.WithValue(ctx, middleware.RequestIDKey, reqID)
		r = r.WithContext(reqIDCtx)
	}

	reqID, ok := ctx.Value(middleware.RequestIDKey).(string)
	if !ok {
		return r, ""
	}

	return r, reqID
}

var _ http.ResponseWriter = &wrappedResponseWriter{} // compile time interface check

type wrappedResponseWriter struct {
	writtenBytes int
	status       int
	start        time.Time
	duration     time.Duration
	http.ResponseWriter
}

func wrapResponseWriter(w http.ResponseWriter, t time.Time) *wrappedResponseWriter {
	return &wrappedResponseWriter{
		ResponseWriter: w,
		start:          t,
	}
}

func (w *wrappedResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *wrappedResponseWriter) Write(data []byte) (int, error) {
	n, err := w.ResponseWriter.Write(data)
	if err != nil {
		return n, err
	}
	w.writtenBytes = n
	w.duration = time.Since(w.start)

	return n, nil
}
