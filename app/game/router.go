package game

import (
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router) {
	c := NewController()

	r.Post("/game/start", c.startGameHandler)
	r.Post("/game/{id}/join", c.joinGameHandler)
	r.Get("/game/{id}", c.gameSnapshotHandler)
	r.Post("/game/{id}/move", c.gameMoveHandler)
	r.Post("/game/{id}/resign", c.gameResignHandler)
}
