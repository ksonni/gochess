package game

type Pawn struct {
	PieceColor PieceColor
}

func (p Pawn) Color() PieceColor {
	return p.PieceColor
}
func (p Pawn) String() string {
	return "p"
}
func (p Pawn) move(from Square, to Square, g *Game) *Board {
	panic("pawn: not implemented")
}
func (p Pawn) canMove(from Square, to Square, g *Game) bool {
	panic("pawn: not implemented")
}
func (p Pawn) computeAttackedSquares(sq Square, g *Game) map[Square]bool {
	panic("pawn: not implemented")
}
