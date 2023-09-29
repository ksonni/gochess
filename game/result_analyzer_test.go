package game

import (
	"testing"
)

func TestResults(t *testing.T) {
	stalemate := DrawReason_Stalemate

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
	}

	for title, pos := range tests {
		g := createPosition(pos.pos, true)
		result := g.computeResult(g)
		if result.Result != pos.result.Result {
			t.Errorf("%s: got game result %v, want %v",
				title, result.Result, GameResult_Checkmate)
		}
		if pos.result.DrawReason != nil {
			if *result.DrawReason != *pos.result.DrawReason {
				t.Errorf("%s: got draw reason %v, want %v",
					title, *result.DrawReason, *pos.result.DrawReason)
			}
		} else if result.DrawReason != nil {
			t.Errorf("%s: draw reason got: %v, want nil", title, *result.DrawReason)
		}
	}
}
