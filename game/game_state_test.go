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

func TestStateHashing(t *testing.T) {
	g := NewGame()
	g.Move(Move{From: sq("e2"), To: sq("e4")})
	hash := g.state.repititionHashableString()
	want := "a1=0_2,a2=0_5,a7=1_5,a8=1_2,b1=0_4,b2=0_5,b7=1_5,b8=1_4,c1=0_3,c2=0_5,c7=1_5,c8=1_3,d1=0_1,d2=0_5,d7=1_5,d8=1_1,e1=0_0,e4=0_5,e7=1_5," +
		"e8=1_0,f1=0_3,f2=0_5,f7=1_5,f8=1_3,g1=0_4,g2=0_5,g7=1_5,g8=1_4,h1=0_2,h2=0_5,h7=1_5,h8=1_2,a1=0,a8=0,e1=0,e8=0,h1=0,h8=0,1"
	if hash != want {
		t.Errorf("Got hash: %s want: %s", hash, want)
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
		return &GameState{board, g.numMoves + 1, g.castlingSquares, g.enpassantTarget, 0, 0}
	} else {
		g.board = board
		return g
	}
}

func gameWithPosition(g *GameState, b *Board) *GameState {
	return &GameState{b, g.numMoves, g.castlingSquares, g.enpassantTarget, 0, 0}
}
