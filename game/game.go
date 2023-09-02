package game

type Game struct {
	Board Board
}

func NewGame() *Game {
	game := new(Game)
	game.Board = NewBoard()
	var initializer initializer
	initializer.initializePieces(game.Board)
	return game
}
