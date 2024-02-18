package game

type Game struct {
	state            *GameState
	repititionHashes map[string]int
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
	GameResult_Active Result = iota
	GameResult_Draw
	GameResult_Checkmate
)

type DrawReason int

const (
	DrawReason_Stalemate DrawReason = iota
	DrawReason_3FoldRepetition
	DrawReason_InusfficientMaterial
	DrawReason_50Moves
)

type ResultData struct {
	Result     Result
	DrawReason *DrawReason
}

const (
	drawRepetitionMoveCount = 3

	// For draw purposes a 'move' consists of a player completing a turn followed by the opponent completing a turn
	drawMoveCount = 50 * 2
)

func NewGame() *Game {
	return &Game{
		state:            NewGameState(),
		repititionHashes: make(map[string]int),
	}
}

func (g *Game) Move(move Move) error {
	state, err := g.state.WithMove(move)
	if err != nil {
		return err
	}
	g.state = state
	hash := state.repititionHashString()
	g.repititionHashes[hash] = g.repititionHashes[hash] + 1
	return nil
}

func (game *Game) ComputeResult() ResultData {
	g := game.state
	side := g.MovingSide()
	result := GameResult_Active
	drawReason := DrawReason_Stalemate

	if possibleMoves := g.PlanPossibleMovesForSide(side); len(possibleMoves) == 0 {
		if g.IsSideInCheck(side) {
			result = GameResult_Checkmate
		} else {
			result = GameResult_Draw
			drawReason = DrawReason_Stalemate
		}
	} else if reason, isDraw := game.testForDraw(game, side); isDraw {
		result = GameResult_Draw
		drawReason = *reason
	}

	out := ResultData{Result: result}
	if result == GameResult_Draw {
		out.DrawReason = &drawReason
	}
	return out
}

// Helpers

func (game *Game) testForDraw(g *Game, color PieceColor) (*DrawReason, bool) {
	drawReason := DrawReason_InusfficientMaterial
	if game.hasInsufficientMaterial(g.state, color) {
		drawReason = DrawReason_InusfficientMaterial
	} else if game.hasReached3FoldRepetition(g, color) {
		drawReason = DrawReason_3FoldRepetition
	} else if game.qualifiesFor50MoveRule(g.state, color) {
		drawReason = DrawReason_50Moves
	} else {
		return nil, false
	}
	return &drawReason, true
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
