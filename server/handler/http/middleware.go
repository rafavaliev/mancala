package http

import (
	"github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"mancala/internal/trace"
	"net/http"
	"time"
)

const traceIDHeader = "X-Trace-id"

func profilingMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			defer func() {
				start := time.Now()
				zap.S().With("status_code", rw.Status()).
					With("http_verb", r.Method).
					With("bytes", rw.BytesWritten()).
					With("latency", time.Since(start).Seconds()).
					With("uri", r.URL.String()).
					With("trace_id", r.Header.Get(traceIDHeader)).
					Info("router: http request")
			}()
			next.ServeHTTP(rw, r)
		})
	}
}

func traceMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			traceID := r.Header.Get(traceIDHeader)
			if traceID == "" {
				traceID = uuid.New().String()
			}
			r.WithContext(trace.WithValue(r.Context(), traceID))
			w.Header().Set(traceIDHeader, traceID)
			next.ServeHTTP(w, r)
		})
	}
}
