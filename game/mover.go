package game

// Conforms to PieceMover
type deltaMover struct{}

func (mover deltaMover) canMove(from Square, to Square, deltas []Square,
	maxSteps int, game *Game) bool {
	for _, delta := range deltas {
		if mover.canMoveWithDelta(from, to, delta, maxSteps, game) {
			return true
		}
	}
	return false
}

func (mover deltaMover) canMoveWithDelta(from Square, to Square, delta Square, maxSteps int, game *Game) bool {
	board := game.Board()
	piece, exists := board.GetPiece(from)
	if !exists {
		return false
	}
	for current, i := from.Adding(delta), 1; board.SquareInRange(current) && (maxSteps == 0 || i <= maxSteps); current, i = current.Adding(delta), i+1 {
		if currentPiece, exists := board.GetPiece(current); exists && currentPiece.Color() == piece.Color() {
			return false
		}
		if current == to {
			return true
		}
	}
	return false
}

func (mover deltaMover) move(from Square, to Square, game *Game) *Board {
	b := game.Board().Clone()
	b.JumpPiece(from, to)
	return &b
}

func (mover deltaMover) computeAttackedSquares(sq Square, deltas []Square, maxSteps int, game *Game) map[Square]bool {
	attacked := make(map[Square]bool)
	for _, delta := range deltas {
		partAttacked := mover.computeAttackedSquaresWithDelta(sq, delta, maxSteps, game)
		for k, v := range partAttacked {
			attacked[k] = v
		}
	}
	return attacked
}

func (mover deltaMover) computeAttackedSquaresWithDelta(
	sq Square, delta Square, maxSteps int, game *Game) map[Square]bool {
	attacked := make(map[Square]bool)

	board := game.Board()
	piece, exists := board.GetPiece(sq)
	if !exists {
		return make(map[Square]bool)
	}
	for current, i := sq.Adding(delta), 1; board.SquareInRange(current) && (maxSteps == 0 || i <= maxSteps); current, i = current.Adding(delta), i+1 {
		if currentPiece, exists := board.GetPiece(current); exists {
			if currentPiece.Color() != piece.Color() {
				attacked[current] = true
			}
			break
		}
		attacked[current] = true
	}
	return attacked
}
