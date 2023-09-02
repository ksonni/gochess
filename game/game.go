package game

type Game struct {
	Board Board
}

func NewGame() *Game {
	var game Game
	var initializer initializer
	initializer.initializePieces(&game.Board)
	return &game
}
