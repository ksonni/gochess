package game

import (
	"testing"
)

type castleData struct {
	clearSquares      []string
	kingMove          testMove
	rookMove          testMove
	otherMoves        []testMove
	possibleKingMoves []string
	castleFails       bool
}

var castles = map[string]castleData{
	"White queen-side castles": {
		clearSquares:      []string{"b1", "c1", "d1"},
		kingMove:          testMove{"e1", "c1"},
		rookMove:          testMove{"a1", "d1"},
		possibleKingMoves: []string{"c1", "d1"},
	},
	"White king-side castles": {
		clearSquares:      []string{"f1", "g1"},
		kingMove:          testMove{"e1", "g1"},
		rookMove:          testMove{"h1", "f1"},
		possibleKingMoves: []string{"f1", "g1"},
	},
	"Black queen-side castles": {
		clearSquares:      []string{"b8", "c8", "d8"},
		otherMoves:        []testMove{{"e2", "e4"}},
		kingMove:          testMove{"e8", "c8"},
		rookMove:          testMove{"a8", "d8"},
		possibleKingMoves: []string{"d8", "c8"},
	},
	"Black king-side castles": {
		clearSquares:      []string{"f8", "g8"},
		otherMoves:        []testMove{{"e2", "e4"}},
		kingMove:          testMove{"e8", "g8"},
		rookMove:          testMove{"h8", "f8"},
		possibleKingMoves: []string{"f8", "g8"},
	},
	"Can't castle when path obstructed": {
		clearSquares:      []string{"f8"},
		otherMoves:        []testMove{{"e2", "e4"}},
		kingMove:          testMove{"e8", "g8"},
		rookMove:          testMove{"h8", "f8"},
		possibleKingMoves: []string{"f8"},
		castleFails:       true,
	},
	"Can't castle while in check": {
		clearSquares: []string{"f8", "g8", "e2", "d7"},
		otherMoves: []testMove{
			{"f1", "b5"}, // Give a check
		},
		kingMove:          testMove{"e8", "g8"},
		rookMove:          testMove{"h8", "f8"},
		possibleKingMoves: []string{"f8"},
		castleFails:       true,
	},
	"Can't castle when path attacked": {
		clearSquares: []string{"f8", "g8", "e2", "f7"},
		otherMoves: []testMove{
			{"f1", "c4"}, // Attack castling square
		},
		kingMove:          testMove{"e8", "g8"},
		rookMove:          testMove{"h8", "f8"},
		possibleKingMoves: []string{"f8"},
		castleFails:       true,
	},
	"Can't castle if king moved previously": {
		clearSquares: []string{"f8", "g8", "d7"},
		otherMoves: []testMove{
			{"a2", "a3"},
			{"e8", "d7"}, // King moved
			{"a3", "a4"},
			{"d7", "e8"}, // King moves back
			{"a4", "a5"},
		},
		kingMove:          testMove{"e8", "g8"},
		rookMove:          testMove{"h8", "f8"},
		possibleKingMoves: []string{"d7", "f8"},
		castleFails:       true,
	},
	"Can't castle if rook moved previously": {
		clearSquares: []string{"f8", "g8", "e2", "h7"},
		otherMoves: []testMove{
			{"a2", "a3"},
			{"h8", "h7"}, // rook moved
			{"a3", "a4"},
			{"a7", "a6"},
			{"a4", "a5"},
			{"h7", "h8"}, // rook moves back
			{"b2", "b3"},
		},
		kingMove:          testMove{"e8", "g8"},
		rookMove:          testMove{"h8", "f8"},
		possibleKingMoves: []string{"f8"},
		castleFails:       true,
	},
}

func TestCastling(t *testing.T) {
	for title, config := range castles {
		g := NewGameState()
		clearSquares(g, config.clearSquares...)
		g = continueGame(title, config.otherMoves, false, g, t)
		kingMove, rookMove, mustFail := config.kingMove, config.rookMove, config.castleFails
		g2, err := g.WithMove(kingMove.Move())
		if mustFail {
			if err == nil {
				t.Errorf("%s: castling did not fail when expected", title)
			}
			continue
		}
		if err != nil {
			t.Errorf("%s: castling failed: %v", title, err)
			continue
		}
		board := g2.Board().Clone()
		if piece, exists := board.GetPiece(sq(kingMove.to)); !exists {
			if _, isKing := piece.(King); !isKing {
				t.Errorf("%s: king not found on %s", title, kingMove.to)
			}
		}
		if _, exists := board.GetPiece(sq(kingMove.from)); exists {
			t.Errorf("%s: king square %s must be cleared after castle", title, kingMove.from)
		}
		if piece, exists := board.GetPiece(sq(rookMove.to)); !exists {
			if _, isRook := piece.(Rook); !isRook {
				t.Errorf("%s: rook not found on %s", title, rookMove.to)
			}
		}
		if _, exists := board.GetPiece(sq(rookMove.from)); exists {
			t.Errorf("%s: rook square %s must be cleared after castle", title, rookMove.from)
		}
	}
}

func TestPossibleCastleMoves(t *testing.T) {
	for title, config := range castles {
		g := NewGameState()
		clearSquares(g, config.clearSquares...)
		g = continueGame(title, config.otherMoves, false, g, t)
		moves := g.PlanPossibleMoves(sq(config.kingMove.from))
		squares := make(map[Square]bool)
		for _, move := range moves {
			squares[move.To] = true
		}
		assertSquareMapEquals(title, squares, config.possibleKingMoves, t)
	}
}
