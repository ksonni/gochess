package game

import "fmt"

type Game struct {
	state            *GameState
	repititionHashes map[string]int
	control          TimeControl
	clocks           map[PieceColor]*Clock
	result           *ResultData
}

type Move struct {
	From      Square
	To        Square
	Promotion Piece
}

type MovePlan struct {
	Move
	Game *GameState
}

type Result int

const (
	GameResult_Draw Result = iota
	GameResult_Checkmate
	GameResult_Timeout
)

type DrawReason int

const (
	DrawReason_None DrawReason = iota
	DrawReason_Stalemate
	DrawReason_3FoldRepetition
	DrawReason_InusfficientMaterial
	DrawReason_50Moves
	DrawReason_InusfficientMaterialTimeout
)

type ResultData struct {
	Result     Result
	DrawReason DrawReason
	Winner     *PieceColor
}

const (
	drawRepetitionMoveCount = 3

	// For draw purposes a 'move' consists of a player completing a turn followed by the opponent completing a turn
	drawMoveCount = 50 * 2
)

func NewGame(control TimeControl) *Game {
	return &Game{
		state:            NewGameState(),
		repititionHashes: make(map[string]int),
		control:          control,
		clocks: map[PieceColor]*Clock{
			PieceColor_White: NewClock(control.Total),
			PieceColor_Black: NewClock(control.Total),
		},
	}
}

func (g *Game) Start() {
	if g.HasEnded() || g.state.NumMoves() > 0 {
		return
	}
	g.clocks[g.state.MovingSide()].Start()
}

func (g *Game) Move(move Move) error {
	if g.HasEnded() {
		return fmt.Errorf("game: game has ended, can not move")
	}
	state, err := g.state.WithMove(move)
	if err != nil {
		return err
	}
	if err = g.toggleClocks(); err != nil {
		return err
	}
	g.state = state
	hash := state.repititionHashString()
	g.repititionHashes[hash] = g.repititionHashes[hash] + 1

	g.result, _ = g.computeResult()

	return nil
}

func (game *Game) Result() (*ResultData, bool) {
	if game.result != nil {
		game.result, _ = game.computeResult()
	}
	return game.result, game.result != nil
}

func (game *Game) HasEnded() bool {
	_, ended := game.Result()
	return ended
}

// Helpers

func (game *Game) computeResult() (*ResultData, bool) {
	g := game.state
	side := g.MovingSide()
	opponent := side.Opponent()

	if hasTime := game.clocks[side].RemainingTime() > 0; !hasTime {
		if game.sideHasMaterialForCheckmate(g, opponent) {
			return &ResultData{Result: GameResult_Timeout, Winner: &opponent}, true
		} else {
			return &ResultData{
				Result:     GameResult_Draw,
				DrawReason: DrawReason_InusfficientMaterialTimeout,
			}, true
		}
	} else if possibleMoves := g.PlanPossibleMovesForSide(side); len(possibleMoves) == 0 {
		if g.IsSideInCheck(side) {
			return &ResultData{Result: GameResult_Checkmate, Winner: &opponent}, true
		} else {
			return &ResultData{
				Result:     GameResult_Draw,
				DrawReason: DrawReason_Stalemate,
			}, true
		}
	} else if reason, isDraw := game.testForDraw(game, side); isDraw {
		return &ResultData{
			Result:     GameResult_Draw,
			DrawReason: reason,
		}, true
	}
	return nil, false
}

func (game *Game) toggleClocks() error {
	side := game.state.MovingSide()

	if clock, _ := game.clocks[side]; clock.Running() {
		clock.Stop()
		if clock.RemainingTime() <= 0 {
			return fmt.Errorf("game: clock ran out of time")
		}
		clock.Increment(game.control.Increment)
	}
	if otherClock, _ := game.clocks[side.Opponent()]; otherClock.RemainingTime() > 0 {
		otherClock.Start()
	}
	return nil
}

func (game *Game) testForDraw(g *Game, color PieceColor) (DrawReason, bool) {
	if game.hasInsufficientMaterial(g.state, color) {
		return DrawReason_InusfficientMaterial, true
	} else if game.hasReached3FoldRepetition(g, color) {
		return DrawReason_3FoldRepetition, true
	} else if game.qualifiesFor50MoveRule(g.state, color) {
		return DrawReason_50Moves, true
	}
	return DrawReason_None, false
}

func (game *Game) sideHasMaterialForCheckmate(g *GameState, color PieceColor) bool {
	counts := g.CountPieces()[color]
	totalPieces := len(counts)
	// Lone/no king
	if totalPieces <= 1 || counts[PieceType_King] <= 0 {
		return false
	}
	// King + Knight/Bishop
	if totalPieces == 2 && (counts[PieceType_Bishop] == 1 || counts[PieceType_Knight] == 1) {
		return false
	}
	return true
}

func (game *Game) hasInsufficientMaterial(g *GameState, color PieceColor) bool {
	counts := g.CountPieces()
	black, white := counts[PieceColor_Black], counts[PieceColor_White]
	nBlack, nWhite := len(black), len(white)
	// Sanity check
	if black[PieceType_King] != 1 || white[PieceType_King] != 1 {
		return false
	}
	// King vs king
	if nBlack == 1 && nWhite == 1 {
		return true
	}
	// King vs King & (Knight/Bishop)
	if nBlack+nWhite == 3 {
		return black[PieceType_Knight] == 1 || black[PieceType_Bishop] == 1 ||
			white[PieceType_Knight] == 1 || white[PieceType_Bishop] == 1
	}
	// King vs King with both sides having an opposite color bishop
	if nBlack == 2 && nWhite == 2 {
		var bishopSqs []Square
		for square, piece := range g.Board().pieces {
			if piece.Type() == PieceType_Bishop {
				bishopSqs = append(bishopSqs, square)
			}
		}
		return len(bishopSqs) == 2 && bishopSqs[0].Color() != bishopSqs[1].Color()
	}
	return false
}

func (game *Game) hasReached3FoldRepetition(g *Game, color PieceColor) bool {
	return g.repititionHashes[g.state.repititionHashString()] >= drawRepetitionMoveCount
}

func (game *Game) qualifiesFor50MoveRule(g *GameState, color PieceColor) bool {
	return g.numMoves-g.lastCaptureMove >= drawMoveCount &&
		g.numMoves-g.lastPawnMove >= drawMoveCount
}
