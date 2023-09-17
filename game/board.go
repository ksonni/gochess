package game

const (
	boardNumFiles = 8
	boardNumRanks = 8
)

type Board struct {
	pieces map[Square]Piece
}

func NewBoard() *Board {
	return &Board{pieces: make(map[Square]Piece)}
}

func (board *Board) GetPiece(square Square) (Piece, bool) {
	piece, exists := board.pieces[square]
	return piece, exists
}

func (board *Board) NumRanks() int {
	return boardNumRanks
}

func (board *Board) NumFiles() int {
	return boardNumFiles
}

func (board *Board) Clone() *Board {
	copy := NewBoard()
	for k, v := range board.pieces {
		copy.pieces[k] = v
	}
	return copy
}

func (board *Board) SquareInRange(square Square) bool {
	return square.File >= 0 && square.File < board.NumFiles() &&
		square.Rank >= 0 && square.Rank < board.NumRanks()
}

func (board *Board) GetKingSquare(color PieceColor) (*Square, bool) {
	for square, piece := range board.pieces {
		if k, ok := piece.(King); ok && k.Color() == color {
			return &square, true
		}
	}
	return nil, false
}

func (board *Board) HasSamePiece(other *Board, square Square) bool {
	p1, _ := board.GetPiece(square)
	p2, _ := other.GetPiece(square)
	if p1 == nil && p2 == nil {
		return true
	}
	if p1 == nil || p2 == nil {
		return false
	}
	return p1.Id() == p2.Id()
}

func (board *Board) FindPiece(piece Piece) (*Square, bool) {
	for square, p := range board.pieces {
		if p.Id() == piece.Id() {
			return &square, true
		}
	}
	return nil, false
}

func (board *Board) setPiece(piece Piece, square Square) {
	board.pieces[square] = piece
}

func (board *Board) clearSquare(square Square) {
	delete(board.pieces, square)
}

func (board *Board) jumpPiece(start Square, end Square) {
	piece, exists := board.GetPiece(start)
	board.clearSquare(start)
	if exists {
		board.setPiece(piece, end)
	}
}
