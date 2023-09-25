package game

type resultAnalyzer struct{}

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

func (a *resultAnalyzer) isKingInCheck(color PieceColor, attackMap map[Square]bool, g *Game) bool {
	square, ok := g.Board().GetKingSquare(color)
	if !ok {
		return false
	}
	return attackMap[*square]
}

func (a *resultAnalyzer) computeResult(g *Game) ResultData {
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
	} else if reason, isDraw := a.testForDraw(g, side); isDraw {
		result = GameResult_Draw
		drawReason = *reason
	}

	out := ResultData{Result: result}
	if result == GameResult_Draw {
		out.DrawReason = &drawReason
	}
	return out
}

func (a *resultAnalyzer) testForDraw(g *Game, color PieceColor) (*DrawReason, bool) {
	drawReason := DrawReason_InusfficientMaterial
	if a.hasInsufficientMaterial(g, color) {
		drawReason = DrawReason_InusfficientMaterial
	} else if a.hasReached3FoldRepetition(g, color) {
		drawReason = DrawReason_3FoldRepetition
	} else if a.qualifiesFor50MoveRule(g, color) {
		drawReason = DrawReason_50Moves
	} else {
		return nil, false
	}
	return &drawReason, true
}

// TODO/implement
func (a *resultAnalyzer) hasInsufficientMaterial(g *Game, color PieceColor) bool {
	return false
}

// TODO/implement
func (a *resultAnalyzer) hasReached3FoldRepetition(g *Game, color PieceColor) bool {
	return false
}

// TODO/implement
func (a *resultAnalyzer) qualifiesFor50MoveRule(g *Game, color PieceColor) bool {
	return false
}
