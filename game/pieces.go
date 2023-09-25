package game

import "math/rand"

type Piece interface {
	PieceMover
	Color() PieceColor
	Id() PieceId
	Type() PieceType
}

type PieceMover interface {
	// WithLocalMove Attempts a move without considering if it might leave the king in check.
	WithLocalMove(move Move, g *Game) (*Game, error)

	PlanPossibleMovesLocally(from Square, g *Game) []MovePlan

	ComputeAttackedSquares(sq Square, g *Game) map[Square]bool
}

type PieceId int

type PieceType int

const (
	PieceType_King PieceType = iota
	PieceType_Queen
	PieceType_Rook_
	PieceType_Bishop_
	PieceType_Knight
	PieceType_Pawn_
)

type PieceProps struct {
	PieceColor PieceColor
	PieceId    PieceId
}

func NewPieceProps(color PieceColor) PieceProps {
	return PieceProps{PieceColor: color, PieceId: PieceId(rand.Int())}
}
func (p PieceProps) Color() PieceColor {
	return p.PieceColor
}
func (p PieceProps) Id() PieceId {
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

type Queen struct {
	deltaMover
	PieceProps
}

func (q Queen) String() string {
	return "Q"
}
func (q Queen) WithLocalMove(move Move, g *Game) (*Game, error) {
	return q.deltaMover.gameWithMove(move.From, move.To, royalDeltas, 0, g)
}
func (q Queen) PlanPossibleMovesLocally(from Square, g *Game) []MovePlan {
	return q.deltaMover.planPossibleMoves(from, royalDeltas, 0, g)
}
func (q Queen) ComputeAttackedSquares(sq Square, g *Game) map[Square]bool {
	return q.deltaMover.computeAttackedSquares(sq, royalDeltas, 0, g)
}
func (q Queen) Type() PieceType {
	return PieceType_Queen
}

type Rook struct {
	deltaMover
	PieceProps
}

func (r Rook) String() string {
	return "R"
}
func (r Rook) WithLocalMove(move Move, g *Game) (*Game, error) {
	return r.deltaMover.gameWithMove(move.From, move.To, perpendicularDeltas, 0, g)
}
func (r Rook) PlanPossibleMovesLocally(from Square, g *Game) []MovePlan {
	return r.deltaMover.planPossibleMoves(from, perpendicularDeltas, 0, g)
}
func (r Rook) ComputeAttackedSquares(sq Square, g *Game) map[Square]bool {
	return r.deltaMover.computeAttackedSquares(sq, perpendicularDeltas, 0, g)
}
func (r Rook) Type() PieceType {
	return PieceType_Rook_
}

type Bishop struct {
	deltaMover
	PieceProps
}

func (b Bishop) WithLocalMove(move Move, g *Game) (*Game, error) {
	return b.deltaMover.gameWithMove(move.From, move.To, diagonalDeltas, 0, g)
}
func (b Bishop) PlanPossibleMovesLocally(from Square, g *Game) []MovePlan {
	return b.deltaMover.planPossibleMoves(from, diagonalDeltas, 0, g)
}
func (b Bishop) ComputeAttackedSquares(sq Square, g *Game) map[Square]bool {
	return b.deltaMover.computeAttackedSquares(sq, diagonalDeltas, 0, g)
}
func (b Bishop) Type() PieceType {
	return PieceType_Bishop_
}

type Knight struct {
	deltaMover
	PieceProps
}

func (k Knight) WithLocalMove(move Move, g *Game) (*Game, error) {
	return k.deltaMover.gameWithMove(move.From, move.To, knightDeltas, 1, g)
}
func (k Knight) PlanPossibleMovesLocally(from Square, g *Game) []MovePlan {
	return k.deltaMover.planPossibleMoves(from, knightDeltas, 1, g)
}
func (k Knight) ComputeAttackedSquares(sq Square, g *Game) map[Square]bool {
	return k.deltaMover.computeAttackedSquares(sq, knightDeltas, 1, g)
}
func (k Knight) Type() PieceType {
	return PieceType_Knight
}
