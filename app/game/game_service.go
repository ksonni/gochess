package game

import (
	"fmt"
	"gochess/lib/game"
	"sync"

	"github.com/google/uuid"
)

type GameService struct {
	games map[uuid.UUID]*GameSession
	mu    sync.RWMutex
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
	session := NewGameSession(ctrl, userId)

	s.mu.Lock()
	defer s.mu.Unlock()

	s.games[id] = session
	return id, nil
}

func (s *GameService) JoinGame(gameId uuid.UUID, userId uuid.UUID) error {
	ch := make(chan error)
	cmd := joinGameCommand{userId, ch}
	if err := s.sendCommand(cmd, gameId); err != nil {
		return err
	}
	return <-ch
}

func (s *GameService) MakeMove(gameId uuid.UUID, userId uuid.UUID, move game.Move) error {
	ch := make(chan error)
	cmd := moveCommand{userId, move, ch}
	if err := s.sendCommand(cmd, gameId); err != nil {
		return err
	}
	return <-ch

}

func (s *GameService) Resign(gameId uuid.UUID, userId uuid.UUID) error {
	ch := make(chan error)
	cmd := resignCommand{userId, ch}
	if err := s.sendCommand(cmd, gameId); err != nil {
		return err
	}
	return <-ch

}

func (s *GameService) SessionSnapshot(gameId uuid.UUID, userId uuid.UUID) (*GameSessionSnapshot, error) {
	ch := make(chan snapshotResult)
	cmd := snapshotCommand{userId, ch}

	s.mu.RLock()
	defer s.mu.RUnlock()

	g, ok := s.games[gameId]
	if !ok {
		return nil, fmt.Errorf("game not found")
	}

	g.ch <- cmd
	result := <-ch

	return result.snapshot, result.err
}

func (s *GameService) CloseSession(gameId uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if g, ok := s.games[gameId]; ok {
		g.Close()
		delete(s.games, gameId)
	}
}

func (s *GameService) sendCommand(cmd sessionCommand, gameId uuid.UUID) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	g, ok := s.games[gameId]
	if !ok {
		return fmt.Errorf("game not found")
	}
	g.ch <- cmd
	return nil
}
