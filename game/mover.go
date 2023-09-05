package game

type mover interface {
	Move(from Square, to Square, board *Board)

	CanMove(from Square, to Square, board *Board) bool

	ComputeAttackedSquares(sq Square, board *Board) map[Square]bool
}

type laterallMover struct{}

func (mover laterallMover) canMove(from Square, to Square, deltas []Square,
	maxSteps int, board *Board) bool {
	for _, delta := range deltas {
		if mover.canMoveWithDelta(from, to, delta, maxSteps, board) {
			return true
		}
	}
	return false
}

func (mover laterallMover) canMoveWithDelta(from Square, to Square, delta Square, maxSteps int, board *Board) bool {
	piece := board.GetPiece(from)
	for current, i := from.Adding(delta), 1; board.HasSquare(current) && (maxSteps == 0 || i <= maxSteps); current, i = current.Adding(delta), i+1 {
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

func (mover laterallMover) Move(from Square, to Square, board *Board) {
	board.SetPiece(board.GetPiece(from), to)
	board.ClearSquare(from)
}

func (mover laterallMover) computeAttackedSquares(sq Square, deltas []Square, maxSteps int, board *Board) map[Square]bool {
	attacked := make(map[Square]bool)
	for _, delta := range deltas {
		partAttacked := mover.computeAttackedSquaresWithDelta(sq, delta, maxSteps, board)
		for k, v := range partAttacked {
			attacked[k] = v
		}
	}
	return attacked
}

func (mover laterallMover) computeAttackedSquaresWithDelta(
	sq Square, delta Square, maxSteps int, board *Board) map[Square]bool {
	attacked := make(map[Square]bool)
	piece := board.GetPiece(sq)
	for current, i := sq.Adding(delta), 1; board.HasSquare(current) && (maxSteps == 0 || i <= maxSteps); current, i = current.Adding(delta), i+1 {
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
