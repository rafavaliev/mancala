package http

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"go.uber.org/zap"
	"mancala/lobby"
	"net/http"
)

func deleteLobby(svc *lobby.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")

		if err := svc.Delete(r.Context(), slug); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func getLobby(svc *lobby.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")

		l, err := svc.Get(r.Context(), slug)
		switch {
		case err == nil:
			break
		case err == lobby.ErrNotFound:
			w.WriteHeader(http.StatusNotFound)
			return
		default:
			zap.S().Error("could not get lobby", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, l)
	}
}

func createLobby(svc *lobby.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		newLobby, err := svc.Create(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, newLobby)
	}
}

func joinLobby(svc *lobby.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")

		l, err := svc.Join(r.Context(), slug)
		switch {
		case err == nil:
			break
		case err == lobby.ErrNotFound:
			w.WriteHeader(http.StatusNotFound)
			return
		default:
			zap.S().Error("could not get lobby", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, l)
	}
}
