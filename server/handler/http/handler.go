package http

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
	"mancala/internal/socket"
	"mancala/lobby"
	"mancala/mancala"
	"net/http"
	"time"
)

// NewHandler return a new router with some handy middleware and api routes
func NewHandler(log *zap.SugaredLogger, lobbySvc *lobby.Repository, service *mancala.Repository) chi.Router {
	r := chi.NewRouter()

	hub := socket.NewHub(lobbySvc, service)
	// Wait for socket messages
	go hub.Run()

	r.Use(
		traceMiddleware(),
		profilingMiddleware(),
		middleware.SetHeader("Content-Type", "application/json"),
		middleware.Timeout(60*time.Second),
	)

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.With(cors.AllowAll().Handler).Route("/v1/lobby", func(r chi.Router) {
		r.Post("/", createLobby(lobbySvc))
		r.Get("/{slug}", getLobby(lobbySvc))
		r.Put("/{slug}/join", joinLobby(lobbySvc))
		r.Delete("/{slug}", deleteLobby(lobbySvc))
	})

	r.With(cors.AllowAll().Handler).Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		socket.ServeWs(hub, w, r)
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	return r
}
