package game

import (
	"fmt"
	"gochess/lib/game"
	"math/rand"

	"github.com/google/uuid"
)

type GameSession struct {
	game  *game.Game
	users map[uuid.UUID]game.PieceColor
	ch    chan sessionCommand
}

type GameSessionSnapshot struct {
	Game  game.GameSnapshot             `json:"game"`
	Users map[uuid.UUID]game.PieceColor `json:"users"`
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

type moveCommand struct {
	userId uuid.UUID
	move   game.Move
	ch     chan<- error
}

type snapshotResult struct {
	snapshot *GameSessionSnapshot
	err      error
}

func NewGameSession(ctrl game.TimeControl, userId uuid.UUID) *GameSession {
	session := GameSession{
		game: game.NewGame(ctrl),
		users: map[uuid.UUID]game.PieceColor{
			userId: game.PieceColor(rand.Intn(2)),
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
		case moveCommand:
			c.ch <- s.makeMove(c.userId, c.move)
		default:
			panic(fmt.Sprintf("Unknown command send to game service: %v", c))
		}
	}
}

func (s *GameSession) joinGame(userId uuid.UUID) error {
	if len(s.users) != 1 {
		return fmt.Errorf("game is already full")
	}

	if _, exists := s.users[userId]; exists {
		return fmt.Errorf("user is already in the game")
	}

	var side game.PieceColor
	for _, v := range s.users {
		side = v
		break
	}

	s.users[userId] = side.Opponent()
	s.game.Start()

	return nil
}

func (s *GameSession) gameSnapshot(userId uuid.UUID) snapshotResult {
	if _, exists := s.users[userId]; !exists {
		return snapshotResult{nil, fmt.Errorf("no permission to access game")}
	}
	snap := &GameSessionSnapshot{
		Game:  s.game.Snapshot(),
		Users: s.users,
	}
	return snapshotResult{snapshot: snap}
}

func (s *GameSession) makeMove(userId uuid.UUID, move game.Move) error {
	if side, exists := s.users[userId]; !exists || side != s.game.MovingSide() {
		return fmt.Errorf("user is not allowed to make this move")
	}
	return s.game.Move(move)
}
