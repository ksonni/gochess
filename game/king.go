package game

import (
	"fmt"
	"math"
)

type King struct {
	PieceProps
	deltaMover deltaMover
}

type castleConfig struct {
	kingDelta     Square
	rookDelta     Square
	rookStartFile int
}

func (k King) String() string {
	return "K"
}

func (k King) PlanMoveLocally(move Move, g *Game) (*Board, error) {
	from, to := move.From, move.To
	// First check if normal movement possible
	result, err := k.deltaMover.planMove(from, to, royalDeltas, 1, g)
	if err == nil {
		return result, nil
	}

	// Attempt castling
	board := g.Board().Clone()
	position := g.Position()
	attackMap := g.ComputeSquaresAttackedBySide(k.Color().Opponent(), &board)

	// King must not be in check
	if attackMap[from] {
		return nil, fmt.Errorf("king: castling not possible when in check")
	}

	// Castling - king must never have moved
	if position.SquareHasEverChanged(from) {
		return nil, fmt.Errorf("king: castling not possible, king has moved in the past")
	}

	for _, config := range k.castleConfigs(&board) {
		kingTargetSquare := from.Adding(config.kingDelta)
		if kingTargetSquare != to {
			continue
		}
		rookSquare := Square{File: config.rookStartFile, Rank: from.Rank}
		// Rook must never have moved
		if position.SquareHasEverChanged(rookSquare) {
			return nil, fmt.Errorf("king: castling not possible, rook has moved in the past")
		}
		// Path must not be under attack
		if !k.castlePathSafe(from, kingTargetSquare, attackMap, &board) {
			return nil, fmt.Errorf("king: castling not possible, path is attacked/obstructed")
		}
		board.JumpPiece(from, kingTargetSquare)
		board.JumpPiece(rookSquare, rookSquare.Adding(config.rookDelta))
		return &board, nil
	}
	return nil, fmt.Errorf("king: not a valid castling move")
}

func (k King) ComputeAttackedSquares(sq Square, g *Game) map[Square]bool {
	return k.deltaMover.computeAttackedSquares(sq, royalDeltas, 1, g)
}

// Helpers

func (k King) castleConfigs(board *Board) []castleConfig {
	return []castleConfig{
		{
			kingDelta:     Square{File: 2, Rank: 0},
			rookDelta:     Square{File: -2, Rank: 0},
			rookStartFile: board.NumFiles() - 1,
		},
		{
			kingDelta:     Square{File: -2, Rank: 0},
			rookDelta:     Square{File: 3, Rank: 0},
			rookStartFile: 0,
		},
	}
}

func (k King) castlePathSafe(from Square, to Square, attackMap map[Square]bool, board *Board) bool {
	delta := Square{File: 1}
	if from.File >= to.File {
		delta.File = -1
	}
	steps := int(math.Abs(float64(from.File - to.File)))
	for step, sq := 0, from.Adding(delta); step < steps; step, sq = step+1, sq.Adding(delta) {
		if _, hasPiece := board.GetPiece(sq); hasPiece || attackMap[sq] {
			return false
		}
	}
	return true
}
