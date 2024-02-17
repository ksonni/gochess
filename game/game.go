package game

type Game struct {
	GameState
	resultAnalyzer
}

type Move struct {
	From      Square
	To        Square
	Promotion Piece
}

type MovePlan struct {
	Move
	Game *GameState
}

func NewGame() *Game {
	return &Game{
		GameState: *NewGameState(),
	}
}
