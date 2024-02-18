package game

import (
	"testing"
)

func TestResults(t *testing.T) {
	stalemate := DrawReason_Stalemate
	insufficientMat := DrawReason_InusfficientMaterial

	tests := map[string]struct {
		pos    map[string]Piece
		result ResultData
	}{
		// Checkmate
		"Sole king checkmate": {
			map[string]Piece{
				"c8": NewKing(PieceColor_Black),
				"h8": NewRook(PieceColor_White),
				"g7": NewRook(PieceColor_White),
				"c1": NewKing(PieceColor_White),
			},
			ResultData{Result: GameResult_Checkmate},
		},
		"Checkmate with an unhelpful friendly piece present": {
			map[string]Piece{
				"c8": NewKing(PieceColor_Black),
				"a1": NewBishop(PieceColor_Black),
				"h8": NewRook(PieceColor_White),
				"g7": NewRook(PieceColor_White),
				"c1": NewKing(PieceColor_White),
			},
			ResultData{Result: GameResult_Checkmate},
		},
		"Checkmate evadable if friendly piece can capture": {
			map[string]Piece{
				"c8": NewKing(PieceColor_Black),
				"h1": NewRook(PieceColor_Black),
				"h8": NewRook(PieceColor_White),
				"g7": NewRook(PieceColor_White),
				"c1": NewKing(PieceColor_White),
			},
			ResultData{Result: GameResult_Active},
		},
		"Checkmate evadable if friendly piece can block": {
			map[string]Piece{
				"c8": NewKing(PieceColor_Black),
				"a5": NewBishop(PieceColor_Black),
				"h8": NewRook(PieceColor_White),
				"g7": NewRook(PieceColor_White),
				"c1": NewKing(PieceColor_White),
			},
			ResultData{Result: GameResult_Active},
		},
		"Checkmate evadable if sole king can move": {
			map[string]Piece{
				"c8": NewKing(PieceColor_Black),
				"h8": NewRook(PieceColor_White),
				"g6": NewRook(PieceColor_White),
				"c1": NewKing(PieceColor_White),
			},
			ResultData{Result: GameResult_Active},
		},

		// Stalemate
		"Stalemate with a sole king": {
			map[string]Piece{
				"a8": NewKing(PieceColor_Black),
				"a7": NewPawn(PieceColor_White),
				"a6": NewKing(PieceColor_White),
			},
			ResultData{Result: GameResult_Draw, DrawReason: &stalemate},
		},
		"Stalemate with other pieces on the board": {
			map[string]Piece{
				"a8": NewKing(PieceColor_Black),
				"a7": NewPawn(PieceColor_White),
				"a6": NewKing(PieceColor_White),
				"d5": NewPawn(PieceColor_Black),
				"d4": NewPawn(PieceColor_White),
			},
			ResultData{Result: GameResult_Draw, DrawReason: &stalemate},
		},
		"Not stalemate if king blocked but other pieces movable": {
			map[string]Piece{
				"a8": NewKing(PieceColor_Black),
				"a7": NewPawn(PieceColor_White),
				"a6": NewKing(PieceColor_White),
				"d5": NewPawn(PieceColor_Black),
			},
			ResultData{Result: GameResult_Active},
		},

		// Insufficient material
		"King vs King insufficient material": {
			map[string]Piece{
				"a8": NewKing(PieceColor_Black),
				"a6": NewKing(PieceColor_White),
			},
			ResultData{Result: GameResult_Draw, DrawReason: &insufficientMat},
		},
		"King vs King & white knight insufficient material": {
			map[string]Piece{
				"a8": NewKing(PieceColor_Black),
				"a6": NewKing(PieceColor_White),
				"a1": NewKnight(PieceColor_White),
			},
			ResultData{Result: GameResult_Draw, DrawReason: &insufficientMat},
		},
		"King vs King & black knight insufficient material": {
			map[string]Piece{
				"a8": NewKing(PieceColor_Black),
				"a6": NewKing(PieceColor_White),
				"a1": NewKnight(PieceColor_Black),
			},
			ResultData{Result: GameResult_Draw, DrawReason: &insufficientMat},
		},
		"King vs King & white bishop insufficient material": {
			map[string]Piece{
				"a8": NewKing(PieceColor_Black),
				"a6": NewKing(PieceColor_White),
				"a1": NewBishop(PieceColor_White),
			},
			ResultData{Result: GameResult_Draw, DrawReason: &insufficientMat},
		},
		"King vs King & black bishop insufficient material": {
			map[string]Piece{
				"a8": NewKing(PieceColor_Black),
				"a6": NewKing(PieceColor_White),
				"a1": NewBishop(PieceColor_Black),
			},
			ResultData{Result: GameResult_Draw, DrawReason: &insufficientMat},
		},
		"King vs King & opposing square colour bishops insufficient material": {
			map[string]Piece{
				"a8": NewKing(PieceColor_Black),
				"a1": NewBishop(PieceColor_Black),
				"a6": NewKing(PieceColor_White),
				"b1": NewBishop(PieceColor_White),
			},
			ResultData{Result: GameResult_Draw, DrawReason: &insufficientMat},
		},
		"King vs King & same square colour bishops, not insufficient material": {
			map[string]Piece{
				"a8": NewKing(PieceColor_Black),
				"a1": NewBishop(PieceColor_Black),
				"a6": NewKing(PieceColor_White),
				"a3": NewBishop(PieceColor_White),
			},
			ResultData{Result: GameResult_Active},
		},
	}

	for title, pos := range tests {
		g := &Game{
			state: emulatePosition(pos.pos, true),
		}
		result := g.ComputeResult()
		if result.Result != pos.result.Result {
			t.Errorf("%s: got game result %v, want %v",
				title, result.Result, GameResult_Checkmate)
		}
		if pos.result.DrawReason != nil {
			if result.DrawReason == nil {
				t.Errorf("%s: got draw reason nil, want %v",
					title, *pos.result.DrawReason)
			} else if *result.DrawReason != *pos.result.DrawReason {
				t.Errorf("%s: got draw reason %v, want %v",
					title, *result.DrawReason, *pos.result.DrawReason)
			}
		} else if result.DrawReason != nil {
			t.Errorf("%s: draw reason got: %v, want nil", title, *result.DrawReason)
		}
	}
}

func Test3FoldRepetition(t *testing.T) {
	type ThreeFoldParams struct {
		moves []testMove
		draw  bool
	}

	tests := map[string]ThreeFoldParams{
		"Repeated knight moves": {
			moves: []testMove{
				{from: "b1", to: "c3"},
				{from: "g8", to: "f6"},
				{from: "c3", to: "b1"},
				{from: "f6", to: "g8"},

				{from: "b1", to: "c3"},
				{from: "g8", to: "f6"},
				{from: "c3", to: "b1"},
				{from: "f6", to: "g8"},

				{from: "b1", to: "c3"},
			},
			draw: true,
		},

		"Must be move of the same player when repetition occurs": {
			moves: []testMove{
				{from: "d2", to: "d4"},
				{from: "d7", to: "d5"}, // Repetition 1 = white's turn next
				{from: "c1", to: "h6"},
				{from: "c8", to: "h3"},

				{from: "h6", to: "c1"},
				{from: "h3", to: "c8"}, // Repetition 2 = white's turn next
				{from: "c1", to: "h6"},
				{from: "c8", to: "h3"},

				{from: "h6", to: "d2"},
				{from: "h3", to: "d7"},
				{from: "d2", to: "e3"},
				{from: "d7", to: "c8"},

				{from: "e3", to: "c1"}, // Position appears again, but it's black's turn
				{from: "c8", to: "d7"},
				{from: "c1", to: "f4"},
				{from: "d7", to: "e6"},

				{from: "f4", to: "e3"},
				{from: "e6", to: "d7"},
				{from: "e3", to: "c1"},
				{from: "d7", to: "c8"}, // Repetition 3 = white's trun next = draw
			},
			draw: true,
		},

		"Must have same en-passant rights when repetition occurs": {
			moves: []testMove{
				{from: "e2", to: "e4"},
				{from: "g8", to: "f6"},
				{from: "e4", to: "e5"},
				{from: "b8", to: "c6"},

				{from: "b1", to: "c3"},
				{from: "d7", to: "d5"}, // Apparant repetition 1 - white's turn and has en-passant rights
				{from: "d1", to: "g4"}, // Actual repetition 1 and black's turn
				{from: "c8", to: "d7"},

				{from: "g4", to: "d1"},
				{from: "d7", to: "c8"}, // Apparant repetition 2 - white's turn but no en-passant rights
				{from: "d1", to: "g4"}, // Actual repetition 2 and black's turn
				{from: "c8", to: "d7"},

				{from: "g4", to: "d1"},
				{from: "d7", to: "c8"}, // Apparant repetition 3 - white's turn but no en-passant rights
				{from: "d1", to: "g4"}, // Actual repetition 3 and black's turn = draw
			},
			draw: true,
		},

		"Must have same castling rights status when repetition occurs - rook variation": {
			moves: []testMove{
				{from: "a2", to: "a3"},
				{from: "a7", to: "a6"}, // Apparant repetition 1 - white's turn and has king-side castling rights

				{from: "a1", to: "a2"}, // Actual repetition 1 and black's turn
				{from: "g8", to: "f6"},
				{from: "a2", to: "a1"},
				{from: "f6", to: "g8"}, // Apparant repetition 2 - white's turn but no castling rights on this king-side

				{from: "a1", to: "a2"}, // Actual repetition 2 and black's turn
				{from: "g8", to: "f6"},
				{from: "a2", to: "a1"},
				{from: "f6", to: "g8"}, // Apparant repetition 3 - white's turn but no castling rights on this king-side

				{from: "a1", to: "a2"}, // Actual repetition 2 and black's turn = draw
			},
			draw: true,
		},

		"Must have same castling rights status when repetition occurs - king variation": {
			moves: []testMove{
				{from: "e2", to: "e3"},
				{from: "e7", to: "e6"}, // Apparant repetition 1 - white's turn and has castling rights

				{from: "e1", to: "e2"}, // Actual repetition 1 and black's turn
				{from: "d8", to: "e7"},
				{from: "e2", to: "e1"},
				{from: "e7", to: "d8"}, // Apparant repetition 2 - white's turn but no castling rights on this king-side

				{from: "e1", to: "e2"}, // Actual repetition 2 and black's turn
				{from: "d8", to: "e7"},
				{from: "e2", to: "e1"},
				{from: "e7", to: "d8"}, // Apparant repetition 3 - white's turn but no castling rights on this king-side

				{from: "e1", to: "e2"}, // Actual repetition 3 and black's turn = draw
			},
			draw: true,
		},
	}

	assertDrawAtEnd := func(title string, params ThreeFoldParams) {
		g := NewGame()
		moves, mustDraw := params.moves, params.draw
		for i, move := range moves {
			err := g.Move(move.Move())
			if err != nil {
				t.Errorf("%s: move %d: failed before testing 3-fold repetition, %v", title, i, err)
				break
			}
			result := g.ComputeResult()
			if i == len(moves)-1 && mustDraw {
				if result.Result != GameResult_Draw {
					t.Errorf("%s: move %d: got game result %d but want %d", title, i, result.Result, GameResult_Draw)
				}
				if result.DrawReason == nil {
					t.Errorf("%s: move %d: got no draw reason, but want %d", title, i, GameResult_Draw)
				} else if *result.DrawReason != DrawReason_3FoldRepetition {
					t.Errorf("%s: move %d: got draw reason %d but want %d", title, i, result.DrawReason, DrawReason_3FoldRepetition)
				}
			} else if result.Result != GameResult_Active {
				t.Errorf("%s: move %d: got game result %d when in progress but want %d", title, i, result.Result, GameResult_Active)
			}
		}
	}

	for title, params := range tests {
		assertDrawAtEnd(title, params)
	}
}
