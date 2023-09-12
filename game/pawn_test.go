package game

import "testing"

var empassantSetupMoves = func() []testMove {
	return []testMove{
		{"a2", "a3"}, {"d7", "d5"},
		{"a3", "a4"}, {"d5", "d4"},
		{"e2", "e4"},
	}
}

var empassantSetupMovesBlack = func() []testMove {
	return []testMove{
		{"e2", "e4"}, {"a7", "a6"},
		{"e4", "e5"}, {"d7", "d5"},
	}
}

func TestPawnAttackedSquares(t *testing.T) {
	runAttackedSquaresTests(map[string]testAttackedSquaresInput{
		"Pawn left flank": {
			[]testMove{},
			"a2",
			[]string{"b3"},
			false,
		},
		"Pawn right flank": {
			[]testMove{},
			"h2",
			[]string{"g3"},
			false,
		},
		"Pawn en-passant": {
			empassantSetupMoves(),
			"d4",
			[]string{"c3", "e3", "e4"},
			false,
		},
		"Pawn en-passant right lost if not used immediately": {
			append(empassantSetupMoves(), testMove{"h7", "h6"}),
			"d4",
			[]string{"c3", "e3"},
			false,
		},
	}, t)
}

func TestEmpassant(t *testing.T) {
	var tests = map[string]struct {
		setupMoves    []testMove
		move          testMove
		captureSquare string
		mustFail      bool
	}{
		"White en-passant": {
			setupMoves:    empassantSetupMoves(),
			move:          testMove{"d4", "e3"},
			captureSquare: "e4",
		},
		"Black en-passant": {
			setupMoves:    empassantSetupMovesBlack(),
			move:          testMove{"e5", "d6"},
			captureSquare: "d4",
		},
		"White en-passant when rights lost": {
			setupMoves: append(
				empassantSetupMoves(),
				testMove{"h7", "h6"},
				testMove{"h2", "h3"},
			),
			move:     testMove{"d4", "e3"},
			mustFail: true,
		},
		"Black en-passant when rights lost": {
			setupMoves: append(
				empassantSetupMovesBlack(),
				testMove{"h2", "h3"},
				testMove{"h7", "h6"},
			),
			move:     testMove{"e5", "d6"},
			mustFail: true,
		},
	}
	for title, test := range tests {
		testEnPassant(title, test.setupMoves, test.move, test.captureSquare, test.mustFail, t)
	}
}

// Helpers

func testEnPassant(title string, setupMoves []testMove, enPassant testMove,
	capturedSq string, mustFail bool, t *testing.T) {
	g := playGame(title, setupMoves, false, t)
	move := enPassant.Move()
	piece, _ := g.Board().GetPiece(move.From)

	if err := g.Move(move); err != nil {
		if !mustFail {
			t.Errorf("%s: valid empassant move from %v to %v failed: %v",
				title, move.From, move.To, err)
		}
		return
	}
	if mustFail {
		t.Errorf("%s: invalid empassant move from %v to %v didn't fail",
			title, move.From, move.To)
		return
	}
	b := g.Board()
	assertSquareEmpty(title, enPassant.from, b, t)
	assertPawn(title, enPassant.to, piece.Color(), b, t)
	assertSquareEmpty(title, capturedSq, b, t)
}

func assertPawn(title string, square string, color PieceColor, board *Board, t *testing.T) {
	p, exists := board.GetPiece(sq(square))
	if !exists {
		t.Errorf("%s: expecting pawn at %s but no piece found", title, square)
		return
	}
	p, isPawn := p.(Pawn)
	if !isPawn && p.Color() != PieceColor_Black {
		t.Errorf("%s: piece at %s is not a %v pawn", title, square, color)
	}
}
