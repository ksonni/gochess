package game

import "math/rand"

type Piece interface {
	PieceMover
	Color() PieceColor
	Id() PieceId
}

type PieceMover interface {
	// Plans a move without considering if it might leave the king in check.
	PlanMoveLocally(from Square, to Square, g *Game) (*Board, error)

	ComputeAttackedSquares(sq Square, g *Game) map[Square]bool
}

type PieceId int

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

type Queen struct {
	deltaMover
	pieceProps
}

func (q Queen) String() string {
	return "Q"
}
func (q Queen) PlanMoveLocally(from Square, to Square, g *Game) (*Board, error) {
	return q.deltaMover.planMove(from, to, royalDeltas, 0, g)
}
func (q Queen) ComputeAttackedSquares(sq Square, g *Game) map[Square]bool {
	return q.deltaMover.computeAttackedSquares(sq, royalDeltas, 0, g)
}

type Rook struct {
	deltaMover
	pieceProps
}

func (r Rook) String() string {
	return "R"
}
func (r Rook) PlanMoveLocally(from Square, to Square, g *Game) (*Board, error) {
	return r.deltaMover.planMove(from, to, perpendicularDeltas, 0, g)
}
func (r Rook) ComputeAttackedSquares(sq Square, g *Game) map[Square]bool {
	return r.deltaMover.computeAttackedSquares(sq, perpendicularDeltas, 0, g)
}

type Bishop struct {
	deltaMover
	pieceProps
}

func (b Bishop) PlanMoveLocally(from Square, to Square, g *Game) (*Board, error) {
	return b.deltaMover.planMove(from, to, diagonalDeltas, 0, g)
}
func (b Bishop) ComputeAttackedSquares(sq Square, g *Game) map[Square]bool {
	return b.deltaMover.computeAttackedSquares(sq, diagonalDeltas, 0, g)
}

type Knight struct {
	deltaMover
	pieceProps
}

func (k Knight) PlanMoveLocally(from Square, to Square, g *Game) (*Board, error) {
	return k.deltaMover.planMove(from, to, knightDeltas, 1, g)
}
func (k Knight) ComputeAttackedSquares(sq Square, g *Game) map[Square]bool {
	return k.deltaMover.computeAttackedSquares(sq, knightDeltas, 1, g)
}
