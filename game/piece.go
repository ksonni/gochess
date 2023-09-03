package game

type PieceColor int

const (
	PieceColor_White PieceColor = iota
	PieceColor_Black
)

type Piece interface {
	Color() PieceColor

	move(from Square, to Square, board *Board) error

	computeAttackedSquares(sq Square, board *Board) map[Square]bool
}

type King struct {
	PieceColor PieceColor
}

func (k King) Color() PieceColor {
	return k.PieceColor
}
func (k King) String() string {
	return "K"
}
func (k King) move(from Square, to Square, board *Board) error {
	panic("Not implemented")
}
func (k King) computeAttackedSquares(sq Square, board *Board) map[Square]bool {
	panic("Not implemented")
}

type Queen struct {
	PieceColor PieceColor
}

func (q Queen) Color() PieceColor {
	return q.PieceColor
}
func (q Queen) String() string {
	return "Q"
}
func (q Queen) move(from Square, to Square, board *Board) error {
	panic("Not implemented")
}
func (q Queen) computeAttackedSquares(sq Square, board *Board) map[Square]bool {
	panic("Not implemented")
}

type Rook struct {
	PieceColor PieceColor
}

func (r Rook) Color() PieceColor {
	return r.PieceColor
}
func (r Rook) String() string {
	return "R"
}
func (r Rook) move(from Square, to Square, board *Board) error {
	panic("Not implemented")
}
func (r Rook) computeAttackedSquares(sq Square, board *Board) map[Square]bool {
	panic("Not implemented")
}

type Bishop struct {
	PieceColor PieceColor
}

func (b Bishop) Color() PieceColor {
	return b.PieceColor
}
func (b Bishop) String() string {
	return "B"
}
func (b Bishop) move(from Square, to Square, board *Board) error {
	panic("Not implemented")
}
func (b Bishop) computeAttackedSquares(sq Square, board *Board) map[Square]bool {
	panic("Not implemented")
}

type Knight struct {
	PieceColor PieceColor
}

func (k Knight) Color() PieceColor {
	return k.PieceColor
}
func (k Knight) String() string {
	return "N"
}
func (k Knight) move(from Square, to Square, board *Board) error {
	panic("Not implemented")
}
func (k Knight) computeAttackedSquares(sq Square, board *Board) map[Square]bool {
	panic("Not implemented")
}

type Pawn struct {
	PieceColor PieceColor
}

func (p Pawn) Color() PieceColor {
	return p.PieceColor
}
func (p Pawn) String() string {
	return "p"
}
func (p Pawn) move(from Square, to Square, board *Board) error {
	panic("Not implemented")
}
func (p Pawn) computeAttackedSquares(sq Square, board *Board) map[Square]bool {
	panic("Not implemented")
}
