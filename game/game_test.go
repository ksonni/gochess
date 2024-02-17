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
	g := NewGame()
	board := g.Board().Clone()

	board.jumpPiece(sq("e8"), sq("e3"))
	g = g.withPosition(board)

	if !g.IsSideInCheck(PieceColor_Black) {
		t.Errorf("Black must be in check")
	}

	board.jumpPiece(sq("e3"), sq("e4"))
	g = g.withPosition(board)

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

func playGame(title string, moves []testMove, jump bool, t *testing.T) *Game {
	return continueGame(title, moves, jump, NewGame(), t)
}

func continueGame(title string, moves []testMove, jump bool, g *Game, t *testing.T) *Game {
	for _, move := range moves {
		if !jump {
			if err := g.Move(move.Move()); err != nil {
				t.Errorf("%s: move %s to %s failed: %v", title, move.from, move.to, err)
				return g
			}
		} else {
			g.Board().jumpPiece(sq(move.from), sq(move.to))
		}
	}
	return g
}

func clearSquares(g *Game, squares ...string) {
	for _, s := range squares {
		g.Board().clearSquare(sq(s))
	}
}

func createPosition(strPieces map[string]Piece, append bool) *Game {
	g := NewGame()
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
		return g.appendingPosition(board)
	} else {
		g.board = board
		return g
	}
}
