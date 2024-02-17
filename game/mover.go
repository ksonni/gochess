package game

import (
	"fmt"
)

// Conforms to PieceMover
type deltaMover struct{}

func (mover deltaMover) canMove(from Square, to Square, deltas []Square,
	maxSteps int, game *GameState) bool {
	for _, delta := range deltas {
		if mover.canMoveWithDelta(from, to, delta, maxSteps, game) {
			return true
		}
	}
	return false
}

func (mover deltaMover) canMoveWithDelta(from Square, to Square, delta Square, maxSteps int, game *GameState) bool {
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

func (mover deltaMover) gameWithMove(from Square, to Square, deltas []Square,
	maxSteps int, game *GameState) (*GameState, error) {
	b, err := mover.planMove(from, to, deltas, maxSteps, game)
	if err != nil {
		return nil, err
	}
	return game.appendingPosition(b, Move{From: from, To: to}, AppendPosParams{}), nil
}

func (mover deltaMover) planMove(from Square, to Square, deltas []Square,
	maxSteps int, game *GameState) (*Board, error) {
	if !mover.canMove(from, to, deltas, maxSteps, game) {
		return nil, fmt.Errorf("board: invalid move")
	}
	b := game.Board().Clone()
	b.jumpPiece(from, to)
	return b, nil
}

func (mover deltaMover) planPossibleMoves(from Square, deltas []Square, maxSteps int, game *GameState) []MovePlan {
	attacked := mover.computeAttackedSquares(from, deltas, maxSteps, game)
	var moves []MovePlan
	for to := range attacked {
		b := game.Board().Clone()
		b.jumpPiece(from, to)
		move := Move{From: from, To: to}
		moves = append(moves, MovePlan{
			Move: move,
			Game: game.appendingPosition(b, move, AppendPosParams{}),
		})
	}
	return moves
}

func (mover deltaMover) computeAttackedSquares(sq Square, deltas []Square, maxSteps int, game *GameState) map[Square]bool {
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
	sq Square, delta Square, maxSteps int, game *GameState) map[Square]bool {
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
