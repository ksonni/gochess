package game

import (
	"runtime"
	"testing"
)

type ResultTestCase struct {
	pos    map[string]Piece
	result ResultData
	/** optional function that will be called to play the position before a result is computed */
	eval func(title string, test ResultTestCase, g *Game) *Game
	/** optional flags used by the `eval` function */
	flags map[int]bool
}

func (test *ResultTestCase) hasFlag(flag int) bool {
	if test.flags == nil {
		return false
	}
	return test.flags[flag]
}

const (
	Flag_PawnMove int = iota
)

func TestResults(t *testing.T) {
	stalemate := DrawReason_Stalemate
	insufficientMat := DrawReason_InusfficientMaterial
	threeFold := DrawReason_3FoldRepetition
	fiftyMoves := DrawReason_50Moves

	must := func(err error) {
		if err != nil {
			_, file, line, _ := runtime.Caller(1)
			t.Fatalf("%s:%d: internal test error: %v", file, line, err)
		}
	}

	play50Moves := func(title string, test ResultTestCase, g *Game) *Game {
		g = createPosition(test.pos, false)

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
			if i == 0 && test.hasFlag(Flag_PawnMove) {
				// move the pawn
				must(g.Move(Move{From: sq("c7"), To: sq("c8"), Promotion: NewQueen(PieceColor_White)}))
			} else {
				must(g.Move(Move{From: whitePos, To: nextWhitePos}))
			}
			if g.computeResult(g).Result != GameResult_Active {
				t.Errorf("%s: triggered prematurely (White's move %d)", title, (i / 2))
			}

			// Black's move
			nextBlackPos := nextPos(blackPos, directionBlack)
			must(g.Move(Move{From: blackPos, To: nextBlackPos}))

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
		return g
	}

	play3Fold := func(title string, test ResultTestCase, g *Game) *Game {
		g = createPosition(test.pos, false)

		for i := 0; i < 2; i++ {
			// zug...
			must(g.Move(Move{From: sq("e3"), To: sq("f3")}))
			must(g.Move(Move{From: sq("e6"), To: sq("d5")}))
			// ...zwang
			must(g.Move(Move{From: sq("f3"), To: sq("e3")}))
			must(g.Move(Move{From: sq("d5"), To: sq("e6")}))
		}
		return g

	}
	tests := map[string]ResultTestCase{
		// Checkmate
		"Sole king checkmate": {
			pos: map[string]Piece{
				"c8": NewKing(PieceColor_Black),
				"h8": NewRook(PieceColor_White),
				"g7": NewRook(PieceColor_White),
				"c1": NewKing(PieceColor_White),
			},
			result: ResultData{Result: GameResult_Checkmate},
		},
		"Checkmate with an unhelpful friendly piece present": {
			pos: map[string]Piece{
				"c8": NewKing(PieceColor_Black),
				"a1": NewBishop(PieceColor_Black),
				"h8": NewRook(PieceColor_White),
				"g7": NewRook(PieceColor_White),
				"c1": NewKing(PieceColor_White),
			},
			result: ResultData{Result: GameResult_Checkmate},
		},
		"Checkmate evadable if friendly piece can capture": {
			pos: map[string]Piece{
				"c8": NewKing(PieceColor_Black),
				"h1": NewRook(PieceColor_Black),
				"h8": NewRook(PieceColor_White),
				"g7": NewRook(PieceColor_White),
				"c1": NewKing(PieceColor_White),
			},
			result: ResultData{Result: GameResult_Active},
		},
		"Checkmate evadable if friendly piece can block": {
			pos: map[string]Piece{
				"c8": NewKing(PieceColor_Black),
				"a5": NewBishop(PieceColor_Black),
				"h8": NewRook(PieceColor_White),
				"g7": NewRook(PieceColor_White),
				"c1": NewKing(PieceColor_White),
			},
			result: ResultData{Result: GameResult_Active},
		},
		"Checkmate evadable if sole king can move": {
			pos: map[string]Piece{
				"c8": NewKing(PieceColor_Black),
				"h8": NewRook(PieceColor_White),
				"g6": NewRook(PieceColor_White),
				"c1": NewKing(PieceColor_White),
			},
			result: ResultData{Result: GameResult_Active},
		},

		// Stalemate
		"Stalemate with a sole king": {
			pos: map[string]Piece{
				"a8": NewKing(PieceColor_Black),
				"a7": NewPawn(PieceColor_White),
				"a6": NewKing(PieceColor_White),
			},
			result: ResultData{Result: GameResult_Draw, DrawReason: &stalemate},
		},
		"Stalemate with other pieces on the board": {
			pos: map[string]Piece{
				"a8": NewKing(PieceColor_Black),
				"a7": NewPawn(PieceColor_White),
				"a6": NewKing(PieceColor_White),
				"d5": NewPawn(PieceColor_Black),
				"d4": NewPawn(PieceColor_White),
			},
			result: ResultData{Result: GameResult_Draw, DrawReason: &stalemate},
		},
		"Not stalemate if king blocked but other pieces movable": {
			pos: map[string]Piece{
				"a8": NewKing(PieceColor_Black),
				"a7": NewPawn(PieceColor_White),
				"a6": NewKing(PieceColor_White),
				"d5": NewPawn(PieceColor_Black),
			},
			result: ResultData{Result: GameResult_Active},
		},

		// Insufficient material
		"King vs King insufficient material": {
			pos: map[string]Piece{
				"a8": NewKing(PieceColor_Black),
				"a6": NewKing(PieceColor_White),
			},
			result: ResultData{Result: GameResult_Draw, DrawReason: &insufficientMat},
		},
		"King vs King & white knight insufficient material": {
			pos: map[string]Piece{
				"a8": NewKing(PieceColor_Black),
				"a6": NewKing(PieceColor_White),
				"a1": NewKnight(PieceColor_White),
			},
			result: ResultData{Result: GameResult_Draw, DrawReason: &insufficientMat},
		},
		"King vs King & black knight insufficient material": {
			pos: map[string]Piece{
				"a8": NewKing(PieceColor_Black),
				"a6": NewKing(PieceColor_White),
				"a1": NewKnight(PieceColor_Black),
			},
			result: ResultData{Result: GameResult_Draw, DrawReason: &insufficientMat},
		},
		"King vs King & white bishop insufficient material": {
			pos: map[string]Piece{
				"a8": NewKing(PieceColor_Black),
				"a6": NewKing(PieceColor_White),
				"a1": NewBishop(PieceColor_White),
			},
			result: ResultData{Result: GameResult_Draw, DrawReason: &insufficientMat},
		},
		"King vs King & black bishop insufficient material": {
			pos: map[string]Piece{
				"a8": NewKing(PieceColor_Black),
				"a6": NewKing(PieceColor_White),
				"a1": NewBishop(PieceColor_Black),
			},
			result: ResultData{Result: GameResult_Draw, DrawReason: &insufficientMat},
		},
		"King vs King & opposing square colour bishops insufficient material": {
			pos: map[string]Piece{
				"a8": NewKing(PieceColor_Black),
				"a1": NewBishop(PieceColor_Black),
				"a6": NewKing(PieceColor_White),
				"b1": NewBishop(PieceColor_White),
			},
			result: ResultData{Result: GameResult_Draw, DrawReason: &insufficientMat},
		},
		"King vs King & same square colour bishops, not insufficient material": {
			pos: map[string]Piece{
				"a8": NewKing(PieceColor_Black),
				"a1": NewBishop(PieceColor_Black),
				"a6": NewKing(PieceColor_White),
				"a3": NewBishop(PieceColor_White),
			},
			result: ResultData{Result: GameResult_Active},
		},

		// 50-move rule
		"50-move rule": {
			pos: map[string]Piece{
				"a8": NewKing(PieceColor_White),
				"h1": NewKing(PieceColor_Black),
				"a7": NewRook(PieceColor_White),
				"h2": NewRook(PieceColor_Black),
			},
			result: ResultData{Result: GameResult_Draw, DrawReason: &fiftyMoves},
			eval:   play50Moves,
		},
		"50-move rule -- should not trigger (capture)": {
			pos: map[string]Piece{
				"a8": NewKing(PieceColor_White),
				"h1": NewKing(PieceColor_Black),
				"a7": NewRook(PieceColor_White),
				"h2": NewRook(PieceColor_Black),
				"b7": NewPawn(PieceColor_Black), // will be captured on move 1
			},
			result: ResultData{Result: GameResult_Active},
			eval:   play50Moves,
		},
		"50-move rule -- should not trigger (pawn move)": {
			pos: map[string]Piece{
				"a8": NewKing(PieceColor_White),
				"h1": NewKing(PieceColor_Black),
				"b7": NewRook(PieceColor_White),
				"h2": NewRook(PieceColor_Black),
				"c7": NewPawn(PieceColor_White), // will be moved on move 1 and promoted to Queen
			},
			result: ResultData{Result: GameResult_Active},
			eval:   play50Moves,
			flags:  map[int]bool{Flag_PawnMove: true},
		},

		// 3-fold repetition
		"3-fold repetition": {
			pos: map[string]Piece{
				"e6": NewKing(PieceColor_Black),
				"e5": NewPawn(PieceColor_Black),
				"e3": NewKing(PieceColor_White),
			},
			result: ResultData{Result: GameResult_Draw, DrawReason: &threeFold},
			eval:   play3Fold,
		},
	}

	testResult := func(title string, test ResultTestCase, g *Game) {
		valueOrNil := func(v *DrawReason) string {
			if v != nil {
				return v.String()
			}
			return "nil"
		}
		result := g.computeResult(g)
		if result.Result != test.result.Result {
			t.Errorf("%s: got game result %v, want %v",
				title, result.Result, test.result.Result)
		} else {
			expected := valueOrNil(test.result.DrawReason)
			actual := valueOrNil(result.DrawReason)
			if expected != actual {
				t.Errorf("%s: draw reason got %s, want %s", title, actual, expected)
			}
		}
	}

	for title, test := range tests {
		g := createPosition(test.pos, true)
		if test.eval != nil {
			g = test.eval(title, test, g)
		}
		testResult(title, test, g)
	}
}

func unreachable(t *testing.T) {
	_, file, line, _ := runtime.Caller(1)
	t.Fatalf("%s:%d: internal test error: reached an unreachable branch", file, line)
}
