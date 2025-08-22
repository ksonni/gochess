package game

import (
	"testing"
	"testing/synctest"
	"time"
)

func TestResults(t *testing.T) {
	white := PieceColor_White

	tests := map[string]struct {
		pos    map[string]Piece
		result *ResultData
		clocks map[PieceColor]*Clock
	}{
		// Timeout
		"Reports timeout": {
			pos: map[string]Piece{
				"e4": NewPawn(PieceColor_White),
				"e5": NewKing(PieceColor_White),
				"a4": NewKing(PieceColor_Black),
				"a5": NewKing(PieceColor_Black),
			},
			result: &ResultData{Result: GameResult_Timeout, Winner: &white},
			clocks: map[PieceColor]*Clock{
				PieceColor_White: NewClock(time.Minute),
				PieceColor_Black: NewClock(0),
			},
		},

		// Checkmate
		"Sole king checkmate": {
			pos: map[string]Piece{
				"c8": NewKing(PieceColor_Black),
				"h8": NewRook(PieceColor_White),
				"g7": NewRook(PieceColor_White),
				"c1": NewKing(PieceColor_White),
			},
			result: &ResultData{Result: GameResult_Checkmate, Winner: &white},
		},
		"Checkmate with an unhelpful friendly piece present": {
			pos: map[string]Piece{
				"c8": NewKing(PieceColor_Black),
				"a1": NewBishop(PieceColor_Black),
				"h8": NewRook(PieceColor_White),
				"g7": NewRook(PieceColor_White),
				"c1": NewKing(PieceColor_White),
			},
			result: &ResultData{Result: GameResult_Checkmate, Winner: &white},
		},
		"Checkmate evadable if friendly piece can capture": {
			pos: map[string]Piece{
				"c8": NewKing(PieceColor_Black),
				"h1": NewRook(PieceColor_Black),
				"h8": NewRook(PieceColor_White),
				"g7": NewRook(PieceColor_White),
				"c1": NewKing(PieceColor_White),
			},
			result: nil,
		},
		"Checkmate evadable if friendly piece can block": {
			pos: map[string]Piece{
				"c8": NewKing(PieceColor_Black),
				"a5": NewBishop(PieceColor_Black),
				"h8": NewRook(PieceColor_White),
				"g7": NewRook(PieceColor_White),
				"c1": NewKing(PieceColor_White),
			},
			result: nil,
		},
		"Checkmate evadable if sole king can move": {
			pos: map[string]Piece{
				"c8": NewKing(PieceColor_Black),
				"h8": NewRook(PieceColor_White),
				"g6": NewRook(PieceColor_White),
				"c1": NewKing(PieceColor_White),
			},
			result: nil,
		},

		// Stalemate
		"Stalemate with a sole king": {
			pos: map[string]Piece{
				"a8": NewKing(PieceColor_Black),
				"a7": NewPawn(PieceColor_White),
				"a6": NewKing(PieceColor_White),
			},
			result: &ResultData{Result: GameResult_Draw, DrawReason: DrawReason_Stalemate},
		},
		"Stalemate with other pieces on the board": {
			pos: map[string]Piece{
				"a8": NewKing(PieceColor_Black),
				"a7": NewPawn(PieceColor_White),
				"a6": NewKing(PieceColor_White),
				"d5": NewPawn(PieceColor_Black),
				"d4": NewPawn(PieceColor_White),
			},
			result: &ResultData{Result: GameResult_Draw, DrawReason: DrawReason_Stalemate},
		},
		"Not stalemate if king blocked but other pieces movable": {
			pos: map[string]Piece{
				"a8": NewKing(PieceColor_Black),
				"a7": NewPawn(PieceColor_White),
				"a6": NewKing(PieceColor_White),
				"d5": NewPawn(PieceColor_Black),
			},
			result: nil,
		},

		// Insufficient material
		"King vs King insufficient material": {
			pos: map[string]Piece{
				"a8": NewKing(PieceColor_Black),
				"a6": NewKing(PieceColor_White),
			},
			result: &ResultData{Result: GameResult_Draw, DrawReason: DrawReason_InusfficientMaterial},
		},
		"King vs King & white knight insufficient material": {
			pos: map[string]Piece{
				"a8": NewKing(PieceColor_Black),
				"a6": NewKing(PieceColor_White),
				"a1": NewKnight(PieceColor_White),
			},
			result: &ResultData{Result: GameResult_Draw, DrawReason: DrawReason_InusfficientMaterial},
		},
		"King vs King & black knight insufficient material": {
			pos: map[string]Piece{
				"a8": NewKing(PieceColor_Black),
				"a6": NewKing(PieceColor_White),
				"a1": NewKnight(PieceColor_Black),
			},
			result: &ResultData{Result: GameResult_Draw, DrawReason: DrawReason_InusfficientMaterial},
		},
		"King vs King & white bishop insufficient material": {
			pos: map[string]Piece{
				"a8": NewKing(PieceColor_Black),
				"a6": NewKing(PieceColor_White),
				"a1": NewBishop(PieceColor_White),
			},
			result: &ResultData{Result: GameResult_Draw, DrawReason: DrawReason_InusfficientMaterial},
		},
		"King vs King & black bishop insufficient material": {
			pos: map[string]Piece{
				"a8": NewKing(PieceColor_Black),
				"a6": NewKing(PieceColor_White),
				"a1": NewBishop(PieceColor_Black),
			},
			result: &ResultData{Result: GameResult_Draw, DrawReason: DrawReason_InusfficientMaterial},
		},
		"King vs King & opposing square colour bishops insufficient material": {
			pos: map[string]Piece{
				"a8": NewKing(PieceColor_Black),
				"a1": NewBishop(PieceColor_Black),
				"a6": NewKing(PieceColor_White),
				"b1": NewBishop(PieceColor_White),
			},
			result: &ResultData{Result: GameResult_Draw, DrawReason: DrawReason_InusfficientMaterial},
		},
		"King vs King & same square colour bishops, not insufficient material": {
			pos: map[string]Piece{
				"a8": NewKing(PieceColor_Black),
				"a1": NewBishop(PieceColor_Black),
				"a6": NewKing(PieceColor_White),
				"a3": NewBishop(PieceColor_White),
			},
			result: nil,
		},
		"Combination of timeout and opponent having insufficient material": {
			pos: map[string]Piece{
				"e5": NewKing(PieceColor_White),
				"a4": NewKing(PieceColor_Black),
				"a5": NewKing(PieceColor_Black),
			},
			result: &ResultData{Result: GameResult_Draw, DrawReason: DrawReason_InusfficientMaterialTimeout},
			clocks: map[PieceColor]*Clock{
				PieceColor_White: NewClock(time.Minute),
				PieceColor_Black: NewClock(0),
			},
		},
	}

	for title, pos := range tests {
		g := NewGame(TimeControl_Thirty)

		g.state = simulatePosition(pos.pos, true)

		if pos.clocks != nil {
			g.clocks = pos.clocks
		}

		g.Start()

		result, _ := g.computeResult()

		if (pos.result != nil) != (result != nil) {
			t.Errorf("%s: got has result data %v, want %v", title, result != nil, pos.result != nil)
		} else if result != nil && pos.result != nil {
			if result.Result != pos.result.Result {
				t.Errorf("%s: got game result %v, want %v", title, result.Result, pos.result.Result)
			}
			if result.DrawReason != pos.result.DrawReason {
				t.Errorf("%s: got draw reason %v, want %v", title, result.DrawReason, pos.result.DrawReason)
			}
			if (result.Winner != nil) != (pos.result.Winner != nil) {
				t.Errorf("%s: got has game winner %v, want %v", title, result.Winner != nil, pos.result.Winner != nil)
			} else if result.Winner != nil && pos.result.Winner != nil && *result.Winner != *pos.result.Winner {
				t.Errorf("%s: got game winner %v, want %v", title, *result.Winner, *pos.result.Winner)
			}
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
		g := NewGame(TimeControl_Thirty)
		g.Start()
		moves, mustDraw := params.moves, params.draw
		for i, move := range moves {
			err := g.Move(move.Move())
			if err != nil {
				t.Errorf("%s: move %d: failed before testing 3-fold repetition, %v", title, i, err)
				break
			}
			result, ok := g.computeResult()
			if i == len(moves)-1 && mustDraw {
				if !ok || result.Result != GameResult_Draw {
					t.Errorf("%s: move %d: got game result %v but want draw", title, i, result)
				}
				if result.DrawReason != DrawReason_3FoldRepetition {
					t.Errorf("%s: move %d: got draw reason %d but want %d", title, i, result.DrawReason, DrawReason_3FoldRepetition)
				}
			} else if result != nil {
				t.Errorf("%s: move %d: got game result %v when in progress but want nil", title, i, result)
			}
		}
	}

	for title, params := range tests {
		assertDrawAtEnd(title, params)
	}
}

func TestTracksCapturesAndPawnMoves(t *testing.T) {
	type CaptureParamsMove struct {
		move         testMove
		lastPawnMove int
		lastCapture  int
	}

	type CapturesParams struct {
		moves []CaptureParamsMove
	}

	tests := map[string]CapturesParams{
		"Tracks normal pawn capture & movement": {
			moves: []CaptureParamsMove{
				{
					move:         testMove{from: "e2", to: "e4"},
					lastPawnMove: 1,
					lastCapture:  0,
				},
				{
					move:         testMove{from: "d7", to: "d5"},
					lastPawnMove: 2,
					lastCapture:  0,
				},
				{
					move:         testMove{from: "e4", to: "d5"},
					lastPawnMove: 3,
					lastCapture:  3,
				},
			},
		},
		"Tracks en-passant pawn capture & movement": {
			moves: []CaptureParamsMove{
				{
					move:         testMove{from: "e2", to: "e4"},
					lastPawnMove: 1,
					lastCapture:  0,
				},
				{
					move:         testMove{from: "g8", to: "h6"},
					lastPawnMove: 1,
					lastCapture:  0,
				},
				{
					move:         testMove{from: "e4", to: "e5"},
					lastPawnMove: 3,
					lastCapture:  0,
				},
				{
					move:         testMove{from: "d7", to: "d5"},
					lastPawnMove: 4,
					lastCapture:  0,
				},
				{
					move:         testMove{from: "e5", to: "d6"},
					lastPawnMove: 5,
					lastCapture:  5,
				},
			},
		},
		"Tracks capture by piece": {
			moves: []CaptureParamsMove{
				{
					move: testMove{from: "b1", to: "c3"},
				},
				{
					move:         testMove{from: "d7", to: "d5"},
					lastPawnMove: 2,
				},
				{
					move:         testMove{from: "c3", to: "d5"},
					lastPawnMove: 2,
					lastCapture:  3,
				},
			},
		},
	}

	assertTracking := func(title string, params CapturesParams) {
		g := NewGame(TimeControl_Thirty)
		g.Start()
		for i, move := range params.moves {
			err := g.Move(move.move.Move())
			if err != nil {
				t.Errorf("%s: move %d: failed before testing pawn/capture tracking, %v", title, i, err)
				break
			}
			if move.lastCapture != g.state.lastCaptureMove {
				t.Errorf("%s: move %d: got lastCaptureMove: %d, want: %d", title, i,
					g.state.lastCaptureMove, move.lastCapture)
			}
			if move.lastPawnMove != g.state.lastPawnMove {
				t.Errorf("%s: move %d: got lastPawnMove: %d, want: %d", title, i,
					g.state.lastPawnMove, move.lastPawnMove)
			}
		}
	}

	for title, params := range tests {
		assertTracking(title, params)
	}
}

func Test50MoveDrawCondition(t *testing.T) {
	rookSq := sq("h1")

	leftExtreme := sq("a2").File
	rightExtreme := sq("g2").File

	delta := []Square{{Rank: 0, File: 1}, {Rank: 0, File: -1}}
	activeSquare := []Square{sq("g2"), sq("g8")}

	state := simulatePosition(map[string]Piece{
		activeSquare[0].String(): NewKing(PieceColor_White),
		rookSq.String():          NewRook(PieceColor_White),
		activeSquare[1].String(): NewKing(PieceColor_Black),
	}, false)

	g := NewGame(TimeControl_Thirty)
	g.Start()
	g.state = state

	for nMove := 0; nMove < 50; nMove++ {
		for nSide := 0; nSide < 2; nSide++ {
			if nSide == 0 && nMove > 0 && nMove%10 == 0 {
				oldSq := rookSq
				rookSq = rookSq.Adding(Square{Rank: 1})
				if err := g.Move(Move{From: oldSq, To: rookSq}); err != nil {
					t.Errorf("Move %d Repition avoidance move failed %v", nMove, err)
				}
				continue
			}
			if activeSquare[nSide].File == leftExtreme {
				delta[nSide].File = 1
			} else if activeSquare[nSide].File == rightExtreme {
				delta[nSide].File = -1
			}
			oldSq := activeSquare[nSide]
			activeSquare[nSide] = activeSquare[nSide].Adding(delta[nSide])
			if err := g.Move(Move{From: oldSq, To: activeSquare[nSide]}); err != nil {
				t.Errorf("Move %d side %d movement failed: from %s to %s: %v",
					nMove, nSide, oldSq.String(), activeSquare[nSide].String(), err)
			}
		}
		result, ok := g.computeResult()
		if nMove == 49 {
			if !ok || result.Result != GameResult_Draw {
				t.Errorf("Iteration %d: Got game result : %v, want draw", nMove, result)
			}
			if result.DrawReason != DrawReason_50Moves {
				t.Errorf("Iteration %d: Got draw reason %d, want: %d", nMove, result.DrawReason, DrawReason_50Moves)
			}
		} else if result != nil {
			t.Errorf("Iteration %d. Got game result : %d, want nil", nMove, result)
		}
	}
}

func TestTracksTime(t *testing.T) {
	g := NewGame(TimeControl_Thirty)

	base := 30 * time.Minute

	params := []timeTrackingParm{
		{
			move:      Move{From: sq("e2"), To: sq("e4")},
			timeTaken: time.Second,
			wantRemaining: map[PieceColor]time.Duration{
				PieceColor_White: base - time.Second,
				PieceColor_Black: base,
			},
		},
		{
			move:      Move{From: sq("e7"), To: sq("e5")},
			timeTaken: 2 * time.Second,
			wantRemaining: map[PieceColor]time.Duration{
				PieceColor_White: base - time.Second,
				PieceColor_Black: base - 2*time.Second,
			},
		},
		{
			move:      Move{From: sq("d2"), To: sq("d4")},
			timeTaken: 1 * time.Second,
			wantRemaining: map[PieceColor]time.Duration{
				PieceColor_White: base - 2*time.Second,
				PieceColor_Black: base - 2*time.Second,
			},
		},
		{
			move:      Move{From: sq("d7"), To: sq("d5")},
			timeTaken: 3 * time.Second,
			wantRemaining: map[PieceColor]time.Duration{
				PieceColor_White: base - 2*time.Second,
				PieceColor_Black: base - 5*time.Second,
			},
		},
	}

	testTimeTracking(t, g, params)
}

func TestTracksTimeWithIncrement(t *testing.T) {
	g := NewGame(TimeControl_TwoOne)

	base := 2 * time.Minute

	params := []timeTrackingParm{
		{
			move:      Move{From: sq("e2"), To: sq("e4")},
			timeTaken: 2 * time.Second,
			wantRemaining: map[PieceColor]time.Duration{
				PieceColor_White: base - time.Second,
				PieceColor_Black: base,
			},
		},
		{
			move:      Move{From: sq("e7"), To: sq("e5")},
			timeTaken: 5 * time.Second,
			wantRemaining: map[PieceColor]time.Duration{
				PieceColor_White: base - time.Second,
				PieceColor_Black: base - 4*time.Second,
			},
		},
		{
			move:      Move{From: sq("d2"), To: sq("d4")},
			timeTaken: 1 * time.Second,
			wantRemaining: map[PieceColor]time.Duration{
				PieceColor_White: base - time.Second,
				PieceColor_Black: base - 4*time.Second,
			},
		},
		{
			move:      Move{From: sq("d7"), To: sq("d5")},
			timeTaken: 3 * time.Second,
			wantRemaining: map[PieceColor]time.Duration{
				PieceColor_White: base - time.Second,
				PieceColor_Black: base - 6*time.Second,
			},
		},
	}
	testTimeTracking(t, g, params)
}

func TestRejectsInvalidMoves(t *testing.T) {
	g := NewGame(TimeControl_Thirty)
	g.Start()
	err := g.Move(Move{sq("a1"), sq("a2"), nil})
	if err == nil {
		t.Errorf("accepted invalid move")
	}
}

func TestRejectsMovesWhenOutOfTime(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		g := NewGame(TimeControl_Thirty)
		g.Start()

		time.Sleep(32 * time.Minute)
		synctest.Wait()

		err := g.Move(Move{sq("e2"), sq("e4"), nil})
		if err == nil {
			t.Errorf("accepted first move when out of time")
		}
	})

	synctest.Test(t, func(t *testing.T) {
		g := NewGame(TimeControl_Thirty)
		g.Start()

		time.Sleep(1 * time.Minute)
		synctest.Wait()

		err := g.Move(Move{sq("e2"), sq("e4"), nil})
		if err != nil {
			t.Errorf("first move failed unexpectedly")
		}
		time.Sleep(32 * time.Minute)
		err = g.Move(Move{sq("e7"), sq("e5"), nil})
		if err == nil {
			t.Errorf("second move did not fail despite running out of time")
		}
	})
}

func TestDrawByAgreement(t *testing.T) {
	g := NewGame(TimeControl_Thirty)
	if err := g.AgreeDraw(); err != nil {
		t.Fatalf("AgreeDraw() errored: %v", err)
	}

	if result, ended := g.Result(); ended {
		if result.DrawReason != DrawReason_Agreement {
			t.Errorf("Result() draw reason got %v, want %v", result.DrawReason, DrawReason_Agreement)
		}
		if result.Result != GameResult_Draw {
			t.Errorf("Result() result got %v, want %v", result.Result, GameResult_Draw)
		}
		if result.Winner != nil {
			t.Errorf("Result() winner got %v, want nil", *result.Winner)
		}
	} else {
		t.Errorf("Result() ended after agreeing draw got %v want %v", false, true)
	}

	move := Move{sq("e2"), sq("e4"), nil}
	if err := g.Move(move); err == nil {
		t.Errorf("Move(%v) was allowed even after agreeing a draw", move)
	}

	if err := g.AgreeDraw(); err == nil {
		t.Errorf("AgreeDraw() was allowed even after the game has ended")
	}
}

func TestResign(t *testing.T) {
	g := NewGame(TimeControl_Thirty)
	if err := g.Resign(PieceColor_White); err != nil {
		t.Fatalf("Resign() errored: %v", err)
	}

	if result, ended := g.Result(); ended {
		if result.Result != GameResult_Resigned {
			t.Errorf("Result() result got %v, want %v", result.Result, GameResult_Resigned)
		}
		if result.Winner == nil {
			t.Errorf("Result() winner got nil, want %v", PieceColor_Black)
		} else if *result.Winner != PieceColor_Black {
			t.Errorf("Result() winner got %v, want %v", *result.Winner, PieceColor_Black)
		}
		if result.DrawReason != DrawReason_None {
			t.Errorf("Result() draw reason got %v, want %v", result.DrawReason, DrawReason_None)
		}
	} else {
		t.Errorf("Result() ended after resigning got %v want %v", false, true)
	}

	move := Move{sq("e2"), sq("e4"), nil}
	if err := g.Move(move); err == nil {
		t.Errorf("Move(%v) was allowed even after resigning", move)
	}

	if err := g.Resign(PieceColor_Black); err == nil {
		t.Errorf("Resign() was allowed even after the game has ended")
	}
}

func TestMoveLogging(t *testing.T) {
	g := NewGame(TimeControl_Thirty)

	g.Start()

	move := Move{From: sq("e2"), To: sq("e4")}
	if err := g.Move(move); err != nil {
		t.Fatalf("Move(%s) failed unexpectedly: %v", move, err)
	}
	if len(g.moves) != 1 {
		t.Errorf("Incorrect moves length after first move got %d, want %d", len(g.moves), 1)
	}
	if actual := g.moves[0]; actual != move {
		t.Errorf("Incorrect first move got %v, want %v", actual, move)
	}

	move = Move{From: sq("e7"), To: sq("e5")}
	if err := g.Move(move); err != nil {
		t.Fatalf("Move(%s) failed unexpectedly: %v", move, err)
	}
	if len(g.moves) != 2 {
		t.Errorf("Incorrect moves length after second move got %d, want %d", len(g.moves), 1)
	}
	if actual := g.moves[1]; actual != move {
		t.Errorf("Incorrect second move got %v, want %v", actual, move)
	}
}

// Helper

type timeTrackingParm struct {
	move          Move
	timeTaken     time.Duration
	wantRemaining map[PieceColor]time.Duration
}

func testTimeTracking(t *testing.T, g *Game, params []timeTrackingParm) {
	synctest.Test(t, func(t *testing.T) {

		whiteClock := g.clocks[PieceColor_White]
		blackClock := g.clocks[PieceColor_Black]

		g.Start()

		for _, p := range params {
			time.Sleep(p.timeTaken)
			synctest.Wait()

			if error := g.Move(p.move); error != nil {
				t.Fatalf("Move(%v) unexpted error occured: %v", p.move, error)
			}

			want := p.wantRemaining[PieceColor_White]
			got := whiteClock.RemainingTime()

			if got != want {
				t.Errorf("Move(%v) got remaining white time %s, want %s", p.move, got, want)
			}
			want = p.wantRemaining[PieceColor_Black]
			got = blackClock.RemainingTime()
			if got != want {
				t.Errorf("Move(%v) got remaining black time %s, want %s", p.move, got, want)
			}
		}
	})
}
