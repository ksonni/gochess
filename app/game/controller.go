package game

import (
	"encoding/json"
	"gochess/auth"
	"gochess/lib/game"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Controller struct {
	service *GameService
}

func NewController() *Controller {
	return &Controller{
		service: NewGameService(),
	}
}

func (c *Controller) startGameHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.Claims(r.Context())
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	var req StartGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}

	control := game.TimeControl{
		Total:     time.Duration(req.DurationMillis) * time.Millisecond,
		Increment: time.Duration(req.IncrementMillis) * time.Millisecond,
	}

	gameId, err := c.service.NewGame(control, user.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("User %s created game %s\n", user.Id, gameId)

	res := StartGameResponse{Id: gameId}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (c *Controller) joinGameHandler(w http.ResponseWriter, r *http.Request) {
	gameId, user, ok := c.gameParams(w, r)
	if !ok {
		return
	}

	if err := c.service.JoinGame(gameId, user.Id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("User %s joined game %s\n", user.Id, gameId)

	w.WriteHeader(http.StatusOK)
}

func (c *Controller) gameSnapshotHandler(w http.ResponseWriter, r *http.Request) {
	gameId, user, ok := c.gameParams(w, r)
	if !ok {
		return
	}

	snap, err := c.service.SessionSnapshot(gameId, user.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err := json.NewEncoder(w).Encode(snap); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (c *Controller) gameMoveHandler(w http.ResponseWriter, r *http.Request) {
	gameId, user, ok := c.gameParams(w, r)
	if !ok {
		return
	}

	var move game.Move
	if err := json.NewDecoder(r.Body).Decode(&move); err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}

	if err := c.service.MakeMove(gameId, user.Id, move); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("User %s made move %v in game %s\n", user.Id, move, gameId)

	w.WriteHeader(http.StatusOK)
}

func (c *Controller) gameResignHandler(w http.ResponseWriter, r *http.Request) {
	gameId, user, ok := c.gameParams(w, r)
	if !ok {
		return
	}

	if err := c.service.Resign(gameId, user.Id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("User %s resigned from in game %s\n", user.Id, gameId)

	w.WriteHeader(http.StatusOK)
}

func (c *Controller) gameParams(w http.ResponseWriter, r *http.Request) (uuid.UUID, *auth.UserClaims, bool) {
	gameId, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid game ID", http.StatusBadRequest)
		return uuid.Nil, nil, false
	}
	user, ok := auth.Claims(r.Context())
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return uuid.Nil, nil, false
	}
	return gameId, &user, true
}
