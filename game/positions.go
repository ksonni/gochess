package game

type Position struct {
	board    *Board
	previous *Position
}

func (p *Position) Appending(board *Board) *Position {
	nextPos := Position{board: board, previous: p}
	return &nextPos
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
