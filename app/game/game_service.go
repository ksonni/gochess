package game

import (
	"fmt"
	"gochess/lib/game"
	"math/rand"
	"sync"

	"github.com/google/uuid"
)

type GameService struct {
	games map[uuid.UUID]*GameSession
	mu    sync.RWMutex
}

type GameSession struct {
	game  *game.Game
	users map[game.PieceColor]uuid.UUID
}

func NewGameService() *GameService {
	return &GameService{
		games: make(map[uuid.UUID]*GameSession),
	}
}

func (s *GameService) NewGame(ctrl game.TimeControl, userId uuid.UUID) (uuid.UUID, error) {
	if !ctrl.Validate() {
		return uuid.Nil, fmt.Errorf("game: invalid time control")
	}

	id := uuid.New()
	session := GameSession{
		game: game.NewGame(ctrl),
		users: map[game.PieceColor]uuid.UUID{
			game.PieceColor(rand.Intn(2)): userId,
		},
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.games[id] = &session
	return id, nil
}
