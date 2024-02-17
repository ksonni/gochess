package game

import (
	"testing"
)

func TestAttackedSquaresBySide(t *testing.T) {
	title := "Attacked squares"
	g := playGame(title, []testMove{{"e2", "e4"}}, false, t)
	attacked := g.ComputeSquaresAttackedBySide(PieceColor_White)
	assertSquareMapEquals(title, attacked, []string{
		"a3", "a6", "b3", "b5", "c3", "c4", "d3", "d5",
		"e2", "e3", "f3", "f5", "g3", "g4", "h3", "h5",
	}, t)
}

func TestIsSideInCheck(t *testing.T) {
	g := NewGameState()
	board := g.Board().Clone()

	board.jumpPiece(sq("e8"), sq("e3"))
	g = gameWithPosition(g, board)

	if !g.IsSideInCheck(PieceColor_Black) {
		t.Errorf("Black must be in check")
	}

	board.jumpPiece(sq("e3"), sq("e4"))
	g = gameWithPosition(g, board)

	if g.IsSideInCheck(PieceColor_Black) {
		t.Errorf("Black must not be in check")
	}
}

// Helpers

type testMove struct {
	from string
	to   string
}

func (m testMove) Move() Move {
	return Move{From: sq(m.from), To: sq(m.to)}
}

func playGame(title string, moves []testMove, jump bool, t *testing.T) *GameState {
	return continueGame(title, moves, jump, NewGameState(), t)
}

func continueGame(title string, moves []testMove, jump bool, g *GameState, t *testing.T) *GameState {
	currentGame := g
	for _, move := range moves {
		if !jump {
			nextGame, err := currentGame.WithMove(move.Move())
			if err != nil {
				t.Errorf("%s: move %s to %s failed: %v", title, move.from, move.to, err)
				return currentGame
			}
			currentGame = nextGame
		} else {
			currentGame.Board().jumpPiece(sq(move.from), sq(move.to))
		}
	}
	return currentGame
}

func clearSquares(g *GameState, squares ...string) {
	for _, s := range squares {
		g.Board().clearSquare(sq(s))
	}
}

func emulatePosition(strPieces map[string]Piece, append bool) *GameState {
	g := NewGameState()
	pieces := make(map[Square]Piece)
	for sqStr, piece := range strPieces {
		pieces[sq(sqStr)] = piece
		// Checking if has moved from original square
		if g.Board().pieces[sq(sqStr)] != pieces[sq(sqStr)] {
			if _, tracked := g.castlingSquares[sq(sqStr)]; tracked {
				g.castlingSquares[sq(sqStr)] = SquareMovementStatus_Moved
			}
		}
	}
	board := &Board{pieces: pieces}
	if append {
		return &GameState{board, g.numMoves + 1, g.castlingSquares, g.enpassantTarget}
	} else {
		g.board = board
		return g
	}
}

func gameWithPosition(g *GameState, b *Board) *GameState {
	return &GameState{b, g.numMoves, g.castlingSquares, g.enpassantTarget}
}
