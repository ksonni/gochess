package game

import (
	"crypto/sha1"
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
)

type GameState struct {
	board           *Board
	numMoves        int
	castlingSquares map[Square]SquareMovementStatus
	enpassantTarget *Square
	lastCaptureMove int
	lastPawnMove    int
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
	square, ok := g.Board().getKingSquare(color)
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
	pawnMove        bool
}

func (g *GameState) appendingPosition(board *Board, move Move, params AppendPosParams) *GameState {
	numMoves := g.numMoves + 1

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

	lastPawnMove := g.lastPawnMove
	if params.pawnMove {
		lastPawnMove = numMoves
	}

	lastCaptureMove := g.lastCaptureMove
	if _, exists := g.board.GetPiece(move.To); exists ||
		len(board.pieces) < len(g.board.pieces) {
		lastCaptureMove = numMoves
	}

	return &GameState{
		board,
		numMoves,
		castlingSquares,
		enpassantTarget,
		lastCaptureMove,
		lastPawnMove,
	}
}

func (g *GameState) repititionHashString() string {
	h := sha1.New()
	h.Write([]byte(g.repititionHashableString()))
	bs := h.Sum(nil)
	sh := string(fmt.Sprintf("%x", bs))
	return sh
}

func (g *GameState) repititionHashableString() string {
	hashes := []string{}

	// Ensure all the pieces and their colors are the same
	for square, piece := range g.board.pieces {
		hash := fmt.Sprintf("%s=%d_%d", square.String(), piece.Color(), piece.Type())
		hashes = append(hashes, hash)
	}
	slices.Sort(hashes)

	// Ensure castling rights are unchanged
	castlingHashes := []string{}
	for square, status := range g.castlingSquares {
		hash := fmt.Sprintf("%s=%d", square.String(), status)
		castlingHashes = append(castlingHashes, hash)
	}
	slices.Sort(castlingHashes)
	hashes = append(hashes, castlingHashes...)

	// Must be the move of the same player
	hashes = append(hashes, fmt.Sprint(g.numMoves%2))

	// Ensure en-passant rights are unchanged
	if g.enpassantTarget != nil {
		hashes = append(hashes, g.enpassantTarget.String())
	}

	return strings.Join(hashes, ",")
}
