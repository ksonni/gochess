package game

type Position struct {
	board    *Board
	previous *Position
}

func (p *Position) Appending(board *Board) *Position {
	nextPos := Position{board: board, previous: p}
	return &nextPos
}

func (p *Position) Setting(board *Board) *Position {
	nextPos := Position{board: board, previous: p.previous}
	return &nextPos
}

func (p *Position) Clone() Position {
	return Position{board: p.board, previous: p.previous}
}

func (p *Position) HasPrevious() bool {
	return p.previous != nil
}

func (p *Position) Previous() (*Position, bool) {
	if p.previous == nil {
		return nil, false
	}
	return p.previous, true
}

func (p *Position) Board() (*Board, bool) {
	if p.board == nil {
		return nil, false
	}
	return p.board, true
}

func (p *Position) PreviousBoard() (*Board, bool) {
	pos, exists := p.Previous()
	if !exists {
		return nil, false
	}
	return pos.Board()
}

func (p *Position) PieceAtPreviousMove(square Square) (Piece, bool) {
	board, ok := p.PreviousBoard()
	if !ok {
		return nil, false
	}
	return board.GetPiece(square)
}
