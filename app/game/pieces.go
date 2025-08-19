package game

import (
	"math/rand"
)

type Piece interface {
	PieceMover
	Color() PieceColor
	Id() PieceId
	Type() PieceType
}

type PieceMover interface {
	// WithLocalMove Attempts a move without considering if it might leave the king in check.
	WithLocalMove(move Move, g *GameState) (*GameState, error)

	PlanPossibleMovesLocally(from Square, g *GameState) []MovePlan

	ComputeAttackedSquares(sq Square, g *GameState) map[Square]bool
}

type PieceId int

type PieceType int

const (
	PieceType_King PieceType = iota
	PieceType_Queen
	PieceType_Rook
	PieceType_Bishop
	PieceType_Knight
	PieceType_Pawn
)

type pieceProps struct {
	PieceColor PieceColor
	PieceId    PieceId
}

func newPieceProps(color PieceColor) pieceProps {
	return pieceProps{PieceColor: color, PieceId: PieceId(rand.Int())}
}
func (p pieceProps) Color() PieceColor {
	return p.PieceColor
}
func (p pieceProps) Id() PieceId {
	return p.PieceId
}

type PieceColor int

const (
	PieceColor_White PieceColor = iota
	PieceColor_Black
)

func (p PieceColor) Opponent() PieceColor {
	switch p {
	case PieceColor_White:
		return PieceColor_Black
	case PieceColor_Black:
		return PieceColor_White
	}
	return PieceColor_Black
}

var (
	perpendicularDeltas = []Square{
		{File: 0, Rank: 1}, {File: 0, Rank: -1},
		{File: 1, Rank: 0}, {File: -1, Rank: 0},
	}
	diagonalDeltas = []Square{
		{File: 1, Rank: 1}, {File: 1, Rank: -1},
		{File: -1, Rank: 1}, {File: -1, Rank: -1},
	}
	knightDeltas = []Square{
		{File: -2, Rank: 1}, {File: 2, Rank: 1},
		{File: -1, Rank: 2}, {File: 1, Rank: 2},
		{File: -2, Rank: -1}, {File: 2, Rank: -1},
		{File: -1, Rank: -2}, {File: 1, Rank: -2},
	}
	royalDeltas []Square
)

func init() {
	royalDeltas = append(royalDeltas, perpendicularDeltas...)
	royalDeltas = append(royalDeltas, diagonalDeltas...)
}

func NewQueen(color PieceColor) Queen {
	return Queen{pieceProps: newPieceProps(color)}
}

type Queen struct {
	deltaMover
	pieceProps
}

func (q Queen) String() string {
	return "Q"
}
func (q Queen) WithLocalMove(move Move, g *GameState) (*GameState, error) {
	return q.deltaMover.gameWithMove(move.From, move.To, royalDeltas, 0, g)
}
func (q Queen) PlanPossibleMovesLocally(from Square, g *GameState) []MovePlan {
	return q.deltaMover.planPossibleMoves(from, royalDeltas, 0, g)
}
func (q Queen) ComputeAttackedSquares(sq Square, g *GameState) map[Square]bool {
	return q.deltaMover.computeAttackedSquares(sq, royalDeltas, 0, g)
}
func (q Queen) Type() PieceType {
	return PieceType_Queen
}

func NewRook(color PieceColor) Rook {
	return Rook{pieceProps: newPieceProps(color)}
}

type Rook struct {
	deltaMover
	pieceProps
}

func (r Rook) String() string {
	return "R"
}
func (r Rook) WithLocalMove(move Move, g *GameState) (*GameState, error) {
	return r.deltaMover.gameWithMove(move.From, move.To, perpendicularDeltas, 0, g)
}
func (r Rook) PlanPossibleMovesLocally(from Square, g *GameState) []MovePlan {
	return r.deltaMover.planPossibleMoves(from, perpendicularDeltas, 0, g)
}
func (r Rook) ComputeAttackedSquares(sq Square, g *GameState) map[Square]bool {
	return r.deltaMover.computeAttackedSquares(sq, perpendicularDeltas, 0, g)
}
func (r Rook) Type() PieceType {
	return PieceType_Rook
}

func NewBishop(color PieceColor) Bishop {
	return Bishop{pieceProps: newPieceProps(color)}
}

type Bishop struct {
	deltaMover
	pieceProps
}

func (b Bishop) WithLocalMove(move Move, g *GameState) (*GameState, error) {
	return b.deltaMover.gameWithMove(move.From, move.To, diagonalDeltas, 0, g)
}
func (b Bishop) PlanPossibleMovesLocally(from Square, g *GameState) []MovePlan {
	return b.deltaMover.planPossibleMoves(from, diagonalDeltas, 0, g)
}
func (b Bishop) ComputeAttackedSquares(sq Square, g *GameState) map[Square]bool {
	return b.deltaMover.computeAttackedSquares(sq, diagonalDeltas, 0, g)
}
func (b Bishop) Type() PieceType {
	return PieceType_Bishop
}

func NewKnight(color PieceColor) Knight {
	return Knight{pieceProps: newPieceProps(color)}
}

type Knight struct {
	deltaMover
	pieceProps
}

func (k Knight) WithLocalMove(move Move, g *GameState) (*GameState, error) {
	return k.deltaMover.gameWithMove(move.From, move.To, knightDeltas, 1, g)
}
func (k Knight) PlanPossibleMovesLocally(from Square, g *GameState) []MovePlan {
	return k.deltaMover.planPossibleMoves(from, knightDeltas, 1, g)
}
func (k Knight) ComputeAttackedSquares(sq Square, g *GameState) map[Square]bool {
	return k.deltaMover.computeAttackedSquares(sq, knightDeltas, 1, g)
}
func (k Knight) Type() PieceType {
	return PieceType_Knight
}
