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
	board := game.Board
	piece := board.GetPiece(from)
	for current, i := from.Adding(delta), 1; board.SquareInRange(current) && (maxSteps == 0 || i <= maxSteps); current, i = current.Adding(delta), i+1 {
		if currentPiece := board.GetPiece(current); currentPiece != nil &&
			currentPiece.Color() == piece.Color() {
			return false
		}
		if current == to {
			return true
		}
	}
	return false
}

func (mover deltaMover) Move(from Square, to Square, game *Game) {
	game.Board.SetPiece(game.Board.GetPiece(from), to)
	game.Board.ClearSquare(from)
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

	board := game.Board
	piece := board.GetPiece(sq)
	for current, i := sq.Adding(delta), 1; board.SquareInRange(current) && (maxSteps == 0 || i <= maxSteps); current, i = current.Adding(delta), i+1 {
		if currentPiece := board.GetPiece(current); currentPiece != nil {
			if currentPiece.Color() != piece.Color() {
				attacked[current] = true
			}
			break
		}
		attacked[current] = true
	}
	return attacked
}
