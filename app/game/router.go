package game

import (
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router) {
	r.Post("/game/start", startGameHandler)
	r.Post("/game/{id}/join", joinGameHandler)
	r.Get("/game/{id}", gameSnapshotHandler)
}
