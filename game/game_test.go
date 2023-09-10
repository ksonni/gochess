package game

import "testing"

func TestAttackedSquaresBySide(t *testing.T) {
	title := "Attacked squares"
	g := playGame(title, []testMove{{"e2", "e4"}}, false, t)
	attacked := g.ComputeSquaresAttackedBySide(PieceColor_White, g.Board())
	assertAttackedSquaresEqual(title, attacked, []string{
		"a3", "a6", "b3", "b5", "c3", "c4", "d3", "d5",
		"e2", "e3", "f3", "f5", "g3", "g4", "h3", "h5",
	}, t)
}

func TestIsSideInCheck(t *testing.T) {
	g := NewGame()
	board := g.Board().Clone()
	board.JumpPiece(sq("e8"), sq("e3"))
	if !g.IsSideInCheck(PieceColor_Black, &board) {
		t.Errorf("Black must be in check")
	}
	board.JumpPiece(sq("e3"), sq("e4"))
	if g.IsSideInCheck(PieceColor_Black, &board) {
		t.Errorf("Black must not be in check")
	}
}

// Helpers

type testMove struct {
	from string
	to   string
}

func playGame(title string, moves []testMove, jump bool, t *testing.T) *Game {
	g := NewGame()
	for _, move := range moves {
		fromSq := sq(move.from)
		toSq := sq(move.to)
		if !jump {
			if err := g.Move(fromSq, toSq); err != nil {
				t.Errorf("%s: move %s to %s failed: %v", title, fromSq, toSq, err)
				return g
			}
		} else {
			g.Board().JumpPiece(fromSq, toSq)
		}
	}
	return g
}

func clearSquares(g *Game, squares ...string) {
	for _, s := range squares {
		g.Board().ClearSquare(sq(s))
	}
}
