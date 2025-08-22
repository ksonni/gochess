package game

import "testing"

var enPassantSetupMoves = func() []testMove {
	return []testMove{
		{"a2", "a3"}, {"d7", "d5"},
		{"a3", "a4"}, {"d5", "d4"},
		{"e2", "e4"},
	}
}

var enPassantSetupMovesBlack = func() []testMove {
	return []testMove{
		{"e2", "e4"}, {"a7", "a6"},
		{"e4", "e5"}, {"d7", "d5"},
	}
}

func TestPawnAttackedSquares(t *testing.T) {
	runAttackedSquaresTests(map[string]targetedSquaresTestInput{
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
			enPassantSetupMoves(),
			"d4",
			[]string{"c3", "e3", "e4"},
			false,
		},
		"Pawn en-passant right lost if not used immediately": {
			append(enPassantSetupMoves(), testMove{"h7", "h6"}),
			"d4",
			[]string{"c3", "e3"},
			false,
		},
	}, t)
}

func TestEnPassant(t *testing.T) {
	var tests = map[string]struct {
		setupMoves     []testMove
		move           testMove
		captureSquare  string
		possibleMoves  []string
		enPassantFails bool
	}{
		"White en-passant": {
			setupMoves:    enPassantSetupMoves(),
			move:          testMove{"d4", "e3"},
			possibleMoves: []string{"d3", "e3"},
			captureSquare: "e4",
		},
		"Black en-passant": {
			setupMoves:    enPassantSetupMovesBlack(),
			move:          testMove{"e5", "d6"},
			possibleMoves: []string{"d6", "e6"},
			captureSquare: "d4",
		},
		"White en-passant when rights lost": {
			setupMoves: append(
				enPassantSetupMoves(),
				testMove{"h7", "h6"},
				testMove{"h2", "h3"},
			),
			move:           testMove{"d4", "e3"},
			possibleMoves:  []string{"d3"},
			enPassantFails: true,
		},
		"Black en-passant when rights lost": {
			setupMoves: append(
				enPassantSetupMovesBlack(),
				testMove{"h2", "h3"},
				testMove{"h7", "h6"},
			),
			move:           testMove{"e5", "d6"},
			possibleMoves:  []string{"e6"},
			enPassantFails: true,
		},
	}
	for title, test := range tests {
		testEnPassant(title, test.setupMoves, test.move,
			test.captureSquare, test.enPassantFails, test.possibleMoves, t)
	}
}

func TestPromotion(t *testing.T) {
	var tests = map[string]struct {
		pawnSquare    string
		move          testMove
		setupMoves    []testMove
		promoted      PieceType
		possibleMoves []string
		cantPromote   bool
	}{
		"Can promote to knight": {
			pawnSquare:    "a2",
			move:          testMove{"a7", "a8"},
			possibleMoves: []string{"a8", "b8"},
			promoted:      PieceType_Knight,
		},
		"Can promote to bishop": {
			pawnSquare:    "a2",
			move:          testMove{"a7", "a8"},
			possibleMoves: []string{"a8", "b8"},
			promoted:      PieceType_Bishop,
		},
		"Can promote to rook": {
			pawnSquare:    "a2",
			move:          testMove{"a7", "a8"},
			possibleMoves: []string{"a8", "b8"},
			promoted:      PieceType_Rook,
		},
		"Can promote to queen": {
			pawnSquare:    "a2",
			move:          testMove{"a7", "a8"},
			possibleMoves: []string{"a8", "b8"},
			promoted:      PieceType_Queen,
		},
		"Can't promote to king": {
			pawnSquare:    "a2",
			move:          testMove{"a7", "a8"},
			possibleMoves: []string{"a8", "b8"},
			promoted:      PieceType_King,
			cantPromote:   true,
		},
		"Can't promote to another pawn": {
			pawnSquare:    "a2",
			move:          testMove{"a7", "a8"},
			possibleMoves: []string{"a8", "b8"},
			promoted:      PieceType_Pawn,
			cantPromote:   true,
		},
		"Promotion from the middle of the board": {
			pawnSquare:    "e2",
			move:          testMove{"e7", "e8"},
			possibleMoves: []string{"e8", "d8", "f8"},
			promoted:      PieceType_Rook,
		},
	}
	for title, test := range tests {
		testPromotion(title, test.pawnSquare, test.move, test.setupMoves,
			test.promoted, test.cantPromote, test.possibleMoves, t)
	}
}

// Helpers

func testEnPassant(title string, setupMoves []testMove, enPassant testMove,
	capturedSq string, mustFail bool, possibleMoves []string, t *testing.T) {
	g := playGame(title, setupMoves, false, t)
	move := enPassant.Move()
	piece, _ := g.Board().GetPiece(move.From)
	moves := g.PlanPossibleMoves(move.From)
	assertMovePlanSquares(title+": possible moves: ", moves, possibleMoves, t)
	g2, err := g.WithMove(move)
	if err != nil {
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
	b := g2.Board()
	assertSquareEmpty(title, enPassant.from, b, t)
	assertPawn(title, enPassant.to, piece.Color(), b, t)
	assertSquareEmpty(title, capturedSq, b, t)
}

func testPromotion(title string, pawnSquare string, promotionMove testMove,
	setupMoves []testMove, promoted PieceType, mustFail bool, possibleMoves []string, t *testing.T) {
	g := playGame(title, setupMoves, false, t)
	clearSquares(g, promotionMove.from, promotionMove.to)
	g.Board().jumpPiece(sq(pawnSquare), sq(promotionMove.from))
	move := Move{
		From:      sq(promotionMove.from),
		To:        sq(promotionMove.to),
		Promotion: &promoted,
	}
	moves := g.PlanPossibleMoves(move.From)
	if lGot, lWant := len(moves), len(possibleMoves)*4; lGot != lWant {
		t.Errorf("%s: promotion got %d possible positions, want %d", title, lGot, lWant)
	}
	assertMovePlanSquares(title, moves, possibleMoves, t)
	g2, err := g.WithMove(move)
	if err != nil {
		if !mustFail {
			t.Errorf("%s: promotion failed: %v", title, err)
		}
		return
	}
	if mustFail {
		t.Errorf("%s: invalid promotion from %s to %s succeeded",
			title, promotionMove.from, promotionMove.to)
		return
	}
	if piece, ok := g2.Board().GetPiece(sq(promotionMove.to)); !ok || piece.Type() != promoted {
		t.Errorf("%s: promoted piece not found on %s", title, promotionMove.to)
	}
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
