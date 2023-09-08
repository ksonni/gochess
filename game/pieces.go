package game

import "math/rand"

type Piece interface {
	pieceMover
	Color() PieceColor
	Id() PieceId
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

type pieceMover interface {
	move(from Square, to Square, g *Game) (*Board, error)

	// TODO: consider global rules before making decision
	canMove(from Square, to Square, g *Game) bool

	computeAttackedSquares(sq Square, g *Game) map[Square]bool
}

type PieceColor int

const (
	PieceColor_White PieceColor = iota
	PieceColor_Black
)

type promotablePiece interface {
	moveAndPromote(from Square, to Square, promotion Piece, g *Game) (*Board, error)
}

// TODO: maybe less dupe?
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

// TODO: castling
type King struct {
	deltaMover
	pieceProps
}

func (k King) String() string {
	return "K"
}
func (k King) canMove(from Square, to Square, g *Game) bool {
	return k.deltaMover.canMove(from, to, royalDeltas, 1, g)
}
func (k King) computeAttackedSquares(sq Square, g *Game) map[Square]bool {
	return k.deltaMover.computeAttackedSquares(sq, royalDeltas, 1, g)
}

type Queen struct {
	deltaMover
	pieceProps
}

func (q Queen) String() string {
	return "Q"
}
func (q Queen) canMove(from Square, to Square, g *Game) bool {
	return q.deltaMover.canMove(from, to, royalDeltas, 0, g)
}
func (q Queen) computeAttackedSquares(sq Square, g *Game) map[Square]bool {
	return q.deltaMover.computeAttackedSquares(sq, royalDeltas, 0, g)
}

type Rook struct {
	deltaMover
	pieceProps
}

func (r Rook) String() string {
	return "R"
}
func (r Rook) canMove(from Square, to Square, g *Game) bool {
	return r.deltaMover.canMove(from, to, perpendicularDeltas, 0, g)
}
func (r Rook) computeAttackedSquares(sq Square, g *Game) map[Square]bool {
	return r.deltaMover.computeAttackedSquares(sq, perpendicularDeltas, 0, g)
}

type Bishop struct {
	deltaMover
	pieceProps
}

func (b Bishop) canMove(from Square, to Square, g *Game) bool {
	return b.deltaMover.canMove(from, to, diagonalDeltas, 0, g)
}
func (b Bishop) computeAttackedSquares(sq Square, g *Game) map[Square]bool {
	return b.deltaMover.computeAttackedSquares(sq, diagonalDeltas, 0, g)
}

type Knight struct {
	deltaMover
	pieceProps
}

func (k Knight) canMove(from Square, to Square, g *Game) bool {
	return k.deltaMover.canMove(from, to, knightDeltas, 1, g)
}
func (k Knight) computeAttackedSquares(sq Square, g *Game) map[Square]bool {
	return k.deltaMover.computeAttackedSquares(sq, knightDeltas, 1, g)
}
