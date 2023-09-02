package game

type PieceColor int

const (
	PieceColor_White PieceColor = iota
	PieceColor_Black
)

type Piece interface {
	Color() PieceColor
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

type Queen struct {
	PieceColor PieceColor
}

func (q Queen) Color() PieceColor {
	return q.PieceColor
}
func (q Queen) String() string {
	return "Q"
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

type Bishop struct {
	PieceColor PieceColor
}

func (b Bishop) Color() PieceColor {
	return b.PieceColor
}
func (b Bishop) String() string {
	return "B"
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

type Pawn struct {
	PieceColor PieceColor
}

func (p Pawn) Color() PieceColor {
	return p.PieceColor
}
func (p Pawn) String() string {
	return "p"
}
