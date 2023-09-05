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

// TODO: implement
func (p Pawn) Move(from Square, to Square, board *Board) {
	panic("pawn: not implemented")
}
func (p Pawn) CanMove(from Square, to Square, board *Board) bool {
	panic("pawn: not implemented")
}
func (p Pawn) ComputeAttackedSquares(sq Square, board *Board) map[Square]bool {
	panic("pawn: not implemented")
}
