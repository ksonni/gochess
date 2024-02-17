package game

import (
	"fmt"
)

type GameState struct {
	board           *Board
	numMoves        int
	castlingSquares map[Square]SquareMovementStatus
	enpassantTarget *Square
}

func (g *GameState) Board() *Board {
	return g.board
}

func (g *GameState) WithMove(move Move) (*GameState, error) {
	piece, exists := g.Board().GetPiece(move.From)
	if !exists {
		return nil, fmt.Errorf("game: no piece exists at %s", move.From)
	}
	if g.MovingSide() != piece.Color() {
		return nil, fmt.Errorf("game: attempted to move piece out of turn")
	}

	nextPos, err := piece.WithLocalMove(move, g)
	if err != nil {
		return nil, fmt.Errorf("game: move failed: %v", err)
	}
	if nextPos.IsSideInCheck(piece.Color()) {
		return nil, fmt.Errorf("game: move failed: violates king integrity")
	}

	return nextPos, nil
}

func (g *GameState) ComputeSquaresAttackedBySide(color PieceColor) map[Square]bool {
	attacked := make(map[Square]bool)
	for square, piece := range g.Board().pieces {
		if piece.Color() != color {
			continue
		}
		pieceAttacked := piece.ComputeAttackedSquares(square, g)
		for pieceSquare, val := range pieceAttacked {
			attacked[pieceSquare] = val
		}
	}
	return attacked
}

func (g *GameState) ComputeAttackedSquares(from Square) map[Square]bool {
	piece, exists := g.Board().GetPiece(from)
	if !exists {
		return make(map[Square]bool)
	}
	return piece.ComputeAttackedSquares(from, g)
}

func (g *GameState) PlanPossibleMovesForSide(color PieceColor) []MovePlan {
	var out []MovePlan
	for sq, piece := range g.Board().pieces {
		if piece.Color() != color {
			continue
		}
		out = append(out, g.PlanPossibleMoves(sq)...)
	}
	return out
}

func (g *GameState) PlanPossibleMoves(from Square) []MovePlan {
	piece, exists := g.Board().GetPiece(from)
	if !exists {
		return []MovePlan{}
	}
	moves := piece.PlanPossibleMovesLocally(from, g)
	var out []MovePlan
	color := piece.Color()
	for _, move := range moves {
		if move.Game.IsSideInCheck(color) {
			continue
		}
		out = append(out, move)
	}
	return out
}

func (g *GameState) IsSideInCheck(color PieceColor) bool {
	attackMap := g.ComputeSquaresAttackedBySide(color.Opponent())
	square, ok := g.Board().GetKingSquare(color)
	if !ok {
		return false
	}
	return attackMap[*square]
}

func (g *GameState) NumMoves() int {
	return g.numMoves
}

func (g *GameState) MovingSide() PieceColor {
	return PieceColor(g.NumMoves() % 2)
}

func (g *GameState) CountPieces() map[PieceColor]map[PieceType]int {
	out := map[PieceColor]map[PieceType]int{
		PieceColor_White: {},
		PieceColor_Black: {},
	}
	for _, piece := range g.Board().pieces {
		out[piece.Color()][piece.Type()] += 1
	}
	return out
}

// Helpers

type AppendPosParams struct {
	enpassantTarget bool
}

func (g *GameState) appendingPosition(board *Board, move Move, params AppendPosParams) *GameState {
	numMoves := g.numMoves + 1

	// En-passant rights
	var enpassantTarget *Square
	if params.enpassantTarget {
		enpassantTarget = &move.To
	}

	castlingSquares := make(map[Square]SquareMovementStatus)
	for square, status := range g.castlingSquares {
		if square == move.From || square == move.To {
			castlingSquares[square] = SquareMovementStatus_Moved
		} else {
			castlingSquares[square] = status
		}
	}

	return &GameState{
		board,
		numMoves,
		castlingSquares,
		enpassantTarget,
	}
}
