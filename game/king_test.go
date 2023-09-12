package game

import (
	"testing"
)

func TestCastling(t *testing.T) {
	successCastles := map[string]struct {
		clearSquares []string
		kingMove     testMove
		rookMove     testMove
		otherMoves   []testMove
		mustFail     bool
	}{
		"White queenside castles": {
			clearSquares: []string{"b1", "c1", "d1"},
			kingMove:     testMove{"e1", "c1"},
			rookMove:     testMove{"a1", "d1"},
		},
		"White kingside castles": {
			clearSquares: []string{"f1", "g1"},
			kingMove:     testMove{"e1", "g1"},
			rookMove:     testMove{"h1", "f1"},
		},
		"Black queenside castles": {
			clearSquares: []string{"b8", "c8", "d8"},
			otherMoves:   []testMove{{"e2", "e4"}},
			kingMove:     testMove{"e8", "c8"},
			rookMove:     testMove{"a8", "d8"},
		},
		"Black kingside castles": {
			clearSquares: []string{"f8", "g8"},
			otherMoves:   []testMove{{"e2", "e4"}},
			kingMove:     testMove{"e8", "g8"},
			rookMove:     testMove{"h8", "f8"},
		},
		"Can't castle when path obstructed": {
			clearSquares: []string{"f8"},
			otherMoves:   []testMove{{"e2", "e4"}},
			kingMove:     testMove{"e8", "g8"},
			rookMove:     testMove{"h8", "f8"},
			mustFail:     true,
		},
		"Can't castle while in check": {
			clearSquares: []string{"f8", "g8", "e2", "d7"},
			otherMoves: []testMove{
				{"f1", "b5"}, // Give a check
			},
			kingMove: testMove{"e8", "g8"},
			rookMove: testMove{"h8", "f8"},
			mustFail: true,
		},
		"Can't castle when path attacked": {
			clearSquares: []string{"f8", "g8", "e2", "f7"},
			otherMoves: []testMove{
				{"f1", "c4"}, // Attack castling square
			},
			kingMove: testMove{"e8", "g8"},
			rookMove: testMove{"h8", "f8"},
			mustFail: true,
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
			kingMove: testMove{"e8", "g8"},
			rookMove: testMove{"h8", "f8"},
			mustFail: true,
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
			kingMove: testMove{"e8", "g8"},
			rookMove: testMove{"h8", "f8"},
			mustFail: true,
		},
	}

	for title, config := range successCastles {
		g := NewGame()
		clearSquares(g, config.clearSquares...)
		for i, move := range config.otherMoves {
			if err := g.Move(move.Move()); err != nil {
				t.Errorf("%s: move %d failed - %s to %s: %v", title, i+1, move.from, move.to, err)
			}
		}
		kingMove, rookMove, mustFail := config.kingMove, config.rookMove, config.mustFail
		err := g.Move(kingMove.Move())
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
		board := g.Board().Clone()
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
