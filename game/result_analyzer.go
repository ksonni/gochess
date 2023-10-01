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

func (a *resultAnalyzer) hasInsufficientMaterial(g *Game, color PieceColor) bool {
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

// TODO/implement
func (a *resultAnalyzer) hasReached3FoldRepetition(g *Game, color PieceColor) bool {
	return false
}

// TODO/implement
func (a *resultAnalyzer) qualifiesFor50MoveRule(g *Game, color PieceColor) bool {
	return false
}
