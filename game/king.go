package game

type King struct {
	deltaMover
	pieceProps
}

func (k King) String() string {
	return "K"
}

// TODO: castling
func (k King) PlanMoveLocally(from Square, to Square, g *Game) (*Board, error) {
	return k.deltaMover.planMove(from, to, royalDeltas, 1, g)
}
func (k King) ComputeAttackedSquares(sq Square, g *Game) map[Square]bool {
	return k.deltaMover.computeAttackedSquares(sq, royalDeltas, 1, g)
}
