package game

type Game struct {
	positions []*Board
	initializer
}

func NewGame() *Game {
	game := new(Game)
	board := make(Board)

	game.positions = []*Board{&board}
	game.initializePieces(game.Board())

	return game
}

func (g *Game) Board() *Board {
	return g.BoardAtMove(g.NumMoves())
}

func (g *Game) PreviousBoard() *Board {
	return g.BoardAtMove(g.NumMoves() - 1)
}

// Moves use a 1 based index because move 0 is a valid position
func (g *Game) BoardAtMove(move int) *Board {
	size := len(g.positions)
	if move < 0 || move > size {
		return nil
	}
	return g.positions[move]
}

func (g *Game) NumMoves() int {
	return len(g.positions) - 1
}
