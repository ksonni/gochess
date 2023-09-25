package game

type resultAnalyzer struct{}

func (a *resultAnalyzer) isSideInCheck(color PieceColor, attackMap map[Square]bool, board *Board, g *Game) bool {
	square, ok := board.GetKingSquare(color)
	if !ok {
		return false
	}
	return attackMap[*square]
}
