package game

type boardAnalyzer struct{}

func (a boardAnalyzer) computeSquaresAttackedBySide(color PieceColor, g *Game) map[Square]bool {
	attacked := make(map[Square]bool)
	for square, piece := range g.Board().pieces {
		if piece.Color() != color {
			continue
		}
		pieceAttacked := piece.ComputeAttackedSquares(square, g)
		for pieceSquare, val := range pieceAttacked {
			attacked[pieceSquare] = val
		}
	}
	return attacked
}

func (a boardAnalyzer) isSideInCheck(color PieceColor, attackMap map[Square]bool, board *Board, g *Game) bool {
	square, ok := board.GetKingSquare(color)
	if !ok {
		return false
	}
	return attackMap[*square]
}
