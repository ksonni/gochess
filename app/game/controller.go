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

var service = NewGameService()

func startGameHandler(w http.ResponseWriter, r *http.Request) {
	var req StartGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
	}

	user, ok := auth.Claims(r.Context())
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	control := game.TimeControl{
		Total:     time.Duration(req.DurationMillis) * time.Millisecond,
		Increment: time.Duration(req.IncrementMillis) * time.Millisecond,
	}

	gameId, err := service.NewGame(control, user.Id)
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

func joinGameHandler(w http.ResponseWriter, r *http.Request) {
	gameId, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid game ID", http.StatusBadRequest)
		return
	}
	user, ok := auth.Claims(r.Context())
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	if err = service.JoinGame(gameId, user.Id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("User %s joined game %s\n", user.Id, gameId)

	w.WriteHeader(http.StatusOK)
}

func gameSnapshotHandler(w http.ResponseWriter, r *http.Request) {
	gameId, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid game ID", http.StatusBadRequest)
		return
	}
	user, ok := auth.Claims(r.Context())
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	snap, err := service.SessionSnapshot(gameId, user.Id)
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
