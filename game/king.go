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

func (k King) WithLocalMove(move Move, g *Game) (*Game, error) {
	from, to := move.From, move.To
	// First check if normal movement possible
	result, err := k.deltaMover.planMove(from, to, royalDeltas, 1, g)
	if err == nil {
		return g.appendingPosition(result), nil
	}

	// Attempt castling
	if moves := k.planPossibleCastlingMoves(from, &to, g); len(moves) > 0 {
		return moves[0].Game, nil
	}
	return nil, fmt.Errorf("king: not a valid move")
}

func (k King) ComputeAttackedSquares(sq Square, g *Game) map[Square]bool {
	return k.deltaMover.computeAttackedSquares(sq, royalDeltas, 1, g)
}

func (k King) PlanPossibleMovesLocally(from Square, g *Game) []MovePlan {
	moves := k.deltaMover.planPossibleMoves(from, royalDeltas, 1, g)
	moves = append(moves, k.planPossibleCastlingMoves(from, nil, g)...)
	return moves
}

// Helpers

func (k King) planPossibleCastlingMoves(from Square, to *Square, g *Game) []MovePlan {
	moves := []MovePlan{}

	// Computed attack map
	attackMap := g.ComputeSquaresAttackedBySide(k.Color().Opponent())

	position := g.Position()

	// King must not be in check
	if attackMap[from] {
		return moves
	}
	// Castling - king must never have moved
	if position.SquareHasEverChanged(from) {
		return moves
	}

	board := g.Board()
	for _, config := range k.castleConfigs(board) {
		kingTargetSquare := from.Adding(config.kingDelta)
		if to != nil && kingTargetSquare != *to {
			continue
		}
		rookSquare := Square{File: config.rookStartFile, Rank: from.Rank}
		// Rook must never have moved
		if position.SquareHasEverChanged(rookSquare) {
			continue
		}
		// Path must not be under attack
		if !k.castlePathSafe(from, kingTargetSquare, attackMap, board) {
			continue
		}
		out := board.Clone()
		out.jumpPiece(from, kingTargetSquare)
		out.jumpPiece(rookSquare, rookSquare.Adding(config.rookDelta))
		moves = append(moves, MovePlan{
			Move: Move{From: from, To: kingTargetSquare},
			Game: g.appendingPosition(out),
		})
	}

	return moves
}

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
