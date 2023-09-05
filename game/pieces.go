package game

type PieceColor int

const (
	PieceColor_White PieceColor = iota
	PieceColor_Black
)

type Piece interface {
	mover
	Color() PieceColor
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

// TODO: castling
type King struct {
	laterallMover
	PieceColor PieceColor
}

func (k King) Color() PieceColor {
	return k.PieceColor
}
func (k King) String() string {
	return "K"
}
func (k King) CanMove(from Square, to Square, board *Board) bool {
	return k.laterallMover.canMove(from, to, royalDeltas, 1, board)
}
func (k King) ComputeAttackedSquares(sq Square, board *Board) map[Square]bool {
	return k.laterallMover.computeAttackedSquares(sq, royalDeltas, 1, board)
}

type Queen struct {
	laterallMover
	PieceColor PieceColor
}

func (q Queen) Color() PieceColor {
	return q.PieceColor
}
func (q Queen) String() string {
	return "Q"
}
func (q Queen) CanMove(from Square, to Square, board *Board) bool {
	return q.laterallMover.canMove(from, to, royalDeltas, 0, board)
}
func (q Queen) ComputeAttackedSquares(sq Square, board *Board) map[Square]bool {
	return q.laterallMover.computeAttackedSquares(sq, royalDeltas, 0, board)
}

type Rook struct {
	laterallMover
	PieceColor PieceColor
}

func (r Rook) Color() PieceColor {
	return r.PieceColor
}
func (r Rook) String() string {
	return "R"
}
func (r Rook) CanMove(from Square, to Square, board *Board) bool {
	return r.laterallMover.canMove(from, to, perpendicularDeltas, 0, board)
}
func (r Rook) ComputeAttackedSquares(sq Square, board *Board) map[Square]bool {
	return r.laterallMover.computeAttackedSquares(sq, perpendicularDeltas, 0, board)
}

type Bishop struct {
	laterallMover
	PieceColor PieceColor
}

func (b Bishop) Color() PieceColor {
	return b.PieceColor
}
func (b Bishop) String() string {
	return "B"
}
func (b Bishop) CanMove(from Square, to Square, board *Board) bool {
	return b.laterallMover.canMove(from, to, diagonalDeltas, 0, board)
}
func (b Bishop) ComputeAttackedSquares(sq Square, board *Board) map[Square]bool {
	return b.laterallMover.computeAttackedSquares(sq, diagonalDeltas, 0, board)
}

type Knight struct {
	laterallMover
	PieceColor PieceColor
}

func (k Knight) Color() PieceColor {
	return k.PieceColor
}
func (k Knight) String() string {
	return "N"
}
func (k Knight) CanMove(from Square, to Square, board *Board) bool {
	return k.laterallMover.canMove(from, to, knightDeltas, 1, board)
}
func (k Knight) ComputeAttackedSquares(sq Square, board *Board) map[Square]bool {
	return k.laterallMover.computeAttackedSquares(sq, knightDeltas, 1, board)
}
