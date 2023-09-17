package game

const (
	boardNumFiles = 8
	boardNumRanks = 8
)

type Board map[Square]Piece

func (board *Board) GetPiece(square Square) (Piece, bool) {
	piece, exists := (*board)[square]
	return piece, exists
}

func (board *Board) NumRanks() int {
	return boardNumRanks
}

func (board *Board) NumFiles() int {
	return boardNumFiles
}

func (board *Board) Clone() Board {
	copy := make(Board)
	for k, v := range *board {
		copy[k] = v
	}
	return copy
}

func (board *Board) SquareInRange(square Square) bool {
	return square.File >= 0 && square.File < board.NumFiles() &&
		square.Rank >= 0 && square.Rank < board.NumRanks()
}

func (board *Board) GetKingSquare(color PieceColor) (*Square, bool) {
	for square, piece := range *board {
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
	for square, p := range *board {
		if p.Id() == piece.Id() {
			return &square, true
		}
	}
	return nil, false
}

func (board *Board) setPiece(piece Piece, square Square) {
	(*board)[square] = piece
}

func (board *Board) clearSquare(square Square) {
	delete(*board, square)
}

func (board *Board) jumpPiece(start Square, end Square) {
	piece, exists := board.GetPiece(start)
	board.clearSquare(start)
	if exists {
		board.setPiece(piece, end)
	}
}
