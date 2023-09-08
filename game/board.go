package game

const (
	boardNumFiles = 8
	boardNumRanks = 8
)

type Board map[Square]Piece

func (board Board) GetPiece(square Square) (Piece, bool) {
	piece, exists := board[square]
	return piece, exists
}

func (board Board) SetPiece(piece Piece, square Square) {
	board[square] = piece
}

func (board Board) ClearSquare(square Square) {
	delete(board, square)
}

func (board Board) JumpPiece(start Square, end Square) {
	piece, exists := board.GetPiece(start)
	board.ClearSquare(start)
	if exists {
		board.SetPiece(piece, end)
	}
}

func (board Board) NumRanks() int {
	return boardNumRanks
}

func (board Board) NumFiles() int {
	return boardNumFiles
}

func (board Board) Clone() Board {
	copy := make(Board)
	for k, v := range board {
		copy[k] = v
	}
	return copy
}

func (board Board) SquareInRange(square Square) bool {
	return square.File >= 0 && square.File < board.NumFiles() &&
		square.Rank >= 0 && square.Rank < board.NumRanks()
}
