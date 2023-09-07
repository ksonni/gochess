package game

type Game struct {
	Board *Board
	initializer
}

func NewGame() *Game {
	game := new(Game)
	board := make(Board)

	game.Board = &board
	game.initializePieces(game.Board)

	return game
}
