package game

import (
	"encoding/json"
	"gochess/auth"
	"gochess/lib/game"
	"log"
	"net/http"
	"time"
)

var service = NewGameService()

func startGameHandler(w http.ResponseWriter, r *http.Request) {
	var req StartGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "failed to decode JSON", http.StatusBadRequest)
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

	id, err := service.NewGame(control, user.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Created game: %s\n", id)

	res := StartGameResponse{Id: id}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
