package game

import (
	"fmt"
	"math"
)

type Game struct {
	board           *Board
	numMoves        int
	castlingSquares map[Square]SquareMovementStatus
	enpessantTarget *Square

	resultAnalyzer
}

type Move struct {
	From      Square
	To        Square
	Promotion Piece
}

type MovePlan struct {
	Move
	Game *Game
}

func (g *Game) Board() *Board {
	return g.board
}

func (g *Game) Move(move Move) error {
	result, err := g.WithMove(move)
	if err != nil {
		return err
	}
	g.numMoves = result.numMoves
	g.board = result.board
	return nil
}

func (g *Game) WithMove(move Move) (*Game, error) {
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

	// For tracking castling rights
	if _, tracked := g.castlingSquares[move.From]; tracked {
		g.castlingSquares[move.From] = SquareMovementStatus_Moved
	}
	if _, tracked := g.castlingSquares[move.To]; tracked {
		g.castlingSquares[move.To] = SquareMovementStatus_Moved
	}

	// For tracking en-passant rights
	moveDist := int(math.Abs(float64(move.From.Rank - move.To.Rank)))
	if _, isPawn := piece.(Pawn); isPawn && moveDist == 2 {
		g.enpessantTarget = &move.To
	} else {
		g.enpessantTarget = nil
	}

	return nextPos, nil
}

func (g *Game) ComputeSquaresAttackedBySide(color PieceColor) map[Square]bool {
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

func (g *Game) ComputeAttackedSquares(from Square) map[Square]bool {
	piece, exists := g.Board().GetPiece(from)
	if !exists {
		return make(map[Square]bool)
	}
	return piece.ComputeAttackedSquares(from, g)
}

func (g *Game) PlanPossibleMovesForSide(color PieceColor) []MovePlan {
	var out []MovePlan
	for sq, piece := range g.Board().pieces {
		if piece.Color() != color {
			continue
		}
		out = append(out, g.PlanPossibleMoves(sq)...)
	}
	return out
}

func (g *Game) PlanPossibleMoves(from Square) []MovePlan {
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

func (g *Game) IsSideInCheck(color PieceColor) bool {
	attackMap := g.ComputeSquaresAttackedBySide(color.Opponent())
	return g.resultAnalyzer.isKingInCheck(color, attackMap, g)
}

func (g *Game) NumMoves() int {
	return g.numMoves
}

func (g *Game) MovingSide() PieceColor {
	return PieceColor(g.NumMoves() % 2)
}

func (g *Game) CountPieces() map[PieceColor]map[PieceType]int {
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

func (g *Game) appendingPosition(board *Board) *Game {
	return &Game{
		board,
		g.numMoves + 1,
		g.castlingSquares,
		g.enpessantTarget,
		g.resultAnalyzer,
	}
}

func (g *Game) withPosition(board *Board) *Game {
	return &Game{
		board,
		g.numMoves,
		g.castlingSquares,
		g.enpessantTarget,
		g.resultAnalyzer,
	}
}
