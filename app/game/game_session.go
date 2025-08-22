package game

import (
	"fmt"
	"gochess/lib/game"
	"math/rand"

	"github.com/google/uuid"
)

type GameSession struct {
	game  *game.Game
	users map[game.PieceColor]uuid.UUID
	ch    chan sessionCommand
}

type GameSessionSnapshot struct {
	Game  game.GameSnapshot             `json:"game"`
	Users map[game.PieceColor]uuid.UUID `json:"users"`
}

type sessionCommand interface{}

type joinGameCommand struct {
	userId uuid.UUID
	ch     chan<- error
}

type snapshotCommand struct {
	userId uuid.UUID
	ch     chan<- snapshotResult
}

type snapshotResult struct {
	snapshot *GameSessionSnapshot
	err      error
}

func NewGameSession(ctrl game.TimeControl, userId uuid.UUID) *GameSession {
	session := GameSession{
		game: game.NewGame(ctrl),
		users: map[game.PieceColor]uuid.UUID{
			game.PieceColor(rand.Intn(2)): userId,
		},
		ch: make(chan sessionCommand),
	}
	go startSession(&session, session.ch)
	return &session
}

func (s *GameSession) Close() {
	close(s.ch)
}

func startSession(s *GameSession, ch <-chan sessionCommand) {
	for cmd := range ch {
		switch c := cmd.(type) {
		case joinGameCommand:
			c.ch <- s.joinGame(c.userId)
		case snapshotCommand:
			c.ch <- s.gameSnapshot(c.userId)
		default:
			panic(fmt.Sprintf("Unknown command send to game service: %v", c))
		}
	}
}

func (s *GameSession) joinGame(userId uuid.UUID) error {
	if len(s.users) != 1 {
		return fmt.Errorf("game is already full")
	}

	var side game.PieceColor
	for k := range s.users {
		side = k
		break
	}

	if existing := s.users[side]; existing == userId {
		return fmt.Errorf("user is already in the game")
	}

	opponent := side.Opponent()
	s.users[opponent] = userId
	s.game.Start()

	return nil
}

func (s *GameSession) gameSnapshot(userId uuid.UUID) snapshotResult {
	var ownsGame bool
	for _, id := range s.users {
		ownsGame = ownsGame || id == userId
	}
	if !ownsGame {
		return snapshotResult{nil, fmt.Errorf("no permission to access game")}
	}
	snap := &GameSessionSnapshot{
		Game:  s.game.Snapshot(),
		Users: s.users,
	}
	return snapshotResult{snapshot: snap}
}
