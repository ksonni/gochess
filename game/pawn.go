package game

// Conforms to promotablePiece
type Pawn struct {
	pieceProps
}

func (p Pawn) Color() PieceColor {
	return p.PieceColor
}
func (p Pawn) String() string {
	return "p"
}

// TODO
func (p Pawn) move(from Square, to Square, g *Game) (*Board, error) {
	panic("pawn: not implemented")
}

// TODO
func (p Pawn) canMove(from Square, to Square, g *Game) bool {
	panic("pawn: not implemented")
}

// TODO
func (p Pawn) computeAttackedSquares(sq Square, g *Game) map[Square]bool {
	panic("pawn: not implemented")
}

// TODO
func (p Pawn) moveAndPromote(from Square, to Square, promotion Piece, g *Game) (*Board, error) {
	panic("pawn: promotion not implemented")
}
