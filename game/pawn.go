package game

type Pawn struct {
	PieceColor PieceColor
	deltaMover
}

func (p Pawn) Color() PieceColor {
	return p.PieceColor
}
func (p Pawn) String() string {
	return "p"
}

func (p Pawn) CanMove(from Square, to Square, g *Game) bool {
	panic("pawn: not implemented")
}
func (p Pawn) ComputeAttackedSquares(sq Square, g *Game) map[Square]bool {
	panic("pawn: not implemented")
}
