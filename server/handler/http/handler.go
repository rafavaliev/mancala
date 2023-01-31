package http

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// NewHandler return a new router with some handy middleware and api routes
func NewHandler(log *zap.SugaredLogger) chi.Router {
	r := chi.NewRouter()

	r.Use(
		traceMiddleware(),
		middleware.SetHeader("Content-Type", "application/json"),
		middleware.Timeout(60*time.Second),
	)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	return r
}
