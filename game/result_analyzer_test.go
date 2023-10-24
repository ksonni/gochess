package game

import (
	"runtime"
	"strings"
	"testing"
)

type ResultTestCase struct {
	pos    map[string]Piece
	result ResultData
}

func TestResults(t *testing.T) {
	stalemate := DrawReason_Stalemate
	insufficientMat := DrawReason_InusfficientMaterial
	fiftyMoves := DrawReason_50Moves
	// threeFold := DrawReason_3FoldRepetition

	tests := map[string]ResultTestCase{
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

	fiftyMoveTests := map[string]ResultTestCase{
		// 50-move rule
		"50-move rule": {
			map[string]Piece{
				"a8": NewKing(PieceColor_White),
				"h1": NewKing(PieceColor_Black),
				"a7": NewRook(PieceColor_White),
				"h2": NewRook(PieceColor_Black),
			},
			ResultData{Result: GameResult_Draw, DrawReason: &fiftyMoves},
		},
		"50-move rule -- should not trigger (capture)": {
			map[string]Piece{
				"a8": NewKing(PieceColor_White),
				"h1": NewKing(PieceColor_Black),
				"a7": NewRook(PieceColor_White),
				"h2": NewRook(PieceColor_Black),
				"b7": NewPawn(PieceColor_Black), // will be captured on move 1
			},
			ResultData{Result: GameResult_Active},
		},
		"50-move rule -- should not trigger (pawn move)": {
			map[string]Piece{
				"a8": NewKing(PieceColor_White),
				"h1": NewKing(PieceColor_Black),
				"b7": NewRook(PieceColor_White),
				"h2": NewRook(PieceColor_Black),
				"c7": NewPawn(PieceColor_White), // will be moved on move 1 and promoted to Queen
			},
			ResultData{Result: GameResult_Active},
		},
	}

	// threeFoldTests := map[string]ResultTestCase{
	// 	3-fold repetition
	// 	"3-fold repetition": {
	// 		map[string]Piece{
	// 			"a8": NewKing(PieceColor_White),
	// 			"h1": NewKing(PieceColor_Black),
	// 			"a7": NewRook(PieceColor_White),
	// 			"h2": NewRook(PieceColor_Black),
	// 		},
	// 		ResultData{Result: GameResult_Active},
	// 		// ResultData{Result: GameResult_Draw, DrawReason: &threeFold},
	// 	},
	// }

	testResult := func(title string, test ResultTestCase, g *Game) {
		result := g.computeResult(g)
		if result.Result != test.result.Result {
			t.Errorf("%s: got game result %v, want %v",
				title, result.Result, test.result.Result)
		}
		if test.result.DrawReason != nil {
			if *result.DrawReason != *test.result.DrawReason {
				t.Errorf("%s: got draw reason %v, want %v",
					title, *result.DrawReason, *test.result.DrawReason)
			}
		} else if result.DrawReason != nil {
			t.Errorf("%s: draw reason got: %v, want nil", title, *result.DrawReason)
		}
	}

	play50Moves := func(title string, test ResultTestCase, g *Game) {
		var err error

		// Make the moves
		type Direction int
		const (
			Direction_Up Direction = iota
			Direction_Down
			Direction_Left
			Direction_Right
		)
		directionWhite := Direction_Right
		directionBlack := Direction_Left
		whitePos := sq("a7")
		blackPos := sq("h2")
		nextPos := func(pos Square, direction Direction) Square {
			switch direction {
			case Direction_Up:
				return pos.Adding(Square{0, 1})
			case Direction_Down:
				return pos.Adding(Square{0, -1})
			case Direction_Left:
				return pos.Adding(Square{-1, 0})
			case Direction_Right:
				return pos.Adding(Square{1, 0})
			}
			unreachable(t)
			panic("unreachable") // prevent compiler from complaining about no return statement
		}
		for i := 0; i < 50; i++ {

			// White's move
			nextWhitePos := nextPos(whitePos, directionWhite)
			if i == 0 && strings.Contains(title, "pawn move") {
				// move the pawn
				err = g.Move(Move{From: sq("c7"), To: sq("c8"), Promotion: NewQueen(PieceColor_White)})
			} else {
				err = g.Move(Move{From: whitePos, To: nextWhitePos})
			}
			if err != nil {
				t.Errorf("%s: unexpected error: %v", title, err)
			}
			if g.computeResult(g).Result != GameResult_Active {
				t.Errorf("%s: triggered prematurely (White's move %d)", title, (i / 2))
			}

			// Black's move
			nextBlackPos := nextPos(blackPos, directionBlack)
			err = g.Move(Move{From: blackPos, To: nextBlackPos})
			if err != nil {
				t.Fatalf("%s: unexpected error: %v", title, err)
			}

			// adjust directions, based on where the pieces came from
			if nextWhitePos.File == 0 {
				if directionWhite == Direction_Left {
					directionWhite = Direction_Down
				} else if directionWhite == Direction_Down {
					directionWhite = Direction_Right
				} else {
					unreachable(t)
				}
			} else if nextWhitePos.File == 6 {
				directionWhite = Direction_Left
			}
			if nextBlackPos.File == 7 {
				if directionBlack == Direction_Right {
					directionBlack = Direction_Up
				} else if directionBlack == Direction_Up {
					directionBlack = Direction_Left
				} else {
					unreachable(t)
				}
			} else if nextBlackPos.File == 1 {
				directionBlack = Direction_Right
			}

			// discard the old positions
			whitePos = nextWhitePos
			blackPos = nextBlackPos

			// check that we're still in the game
			if i < (50-1) && g.computeResult(g).Result != GameResult_Active {
				t.Fatalf("%s: triggered prematurely (Black's move %d)", title, (i / 2))
			}
		}
	}

	for title, test := range tests {
		g := createPosition(test.pos, true)
		testResult(title, test, g)
	}

	for title, test := range fiftyMoveTests {
		g := createPosition(test.pos, false)
		play50Moves(title, test, g)
		testResult(title, test, g)
	}

}

func unreachable(t *testing.T) {
	_, file, line, _ := runtime.Caller(1)
	t.Fatalf("%s:%d: internal test error: reached an unreachable branch", file, line)
}
