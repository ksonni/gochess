package game

type Position struct {
	board *Board
}

func (p *Position) Appending(board *Board) *Position {
	nextPos := Position{board: board}
	return &nextPos
}

func (p *Position) Setting(board *Board) *Position {
	nextPos := Position{board: board}
	return &nextPos
}

func (p *Position) Clone() Position {
	return Position{board: p.board}
}

func (p *Position) Board() (*Board, bool) {
	if p.board == nil {
		return nil, false
	}
	return p.board, true
}
