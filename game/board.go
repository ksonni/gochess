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

func (board *Board) ContainsSquare(square Square) bool {
	return square.File >= 0 && square.File < board.NumFiles() &&
		square.Rank >= 0 && square.Rank < board.NumRanks()
}

// Helpers

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

func (board *Board) getKingSquare(color PieceColor) (*Square, bool) {
	for square, piece := range board.pieces {
		if k, ok := piece.(King); ok && k.Color() == color {
			return &square, true
		}
	}
	return nil, false
}
