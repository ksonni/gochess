package game

type Game struct {
	state      *GameState
	gameHashes map[string]int
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
		state:      NewGameState(),
		gameHashes: make(map[string]int),
	}
}

func (g *Game) Move(move Move) error {
	state, err := g.state.WithMove(move)
	if err != nil {
		return err
	}
	g.state = state
	hash := state.repititionHashString()
	g.gameHashes[hash] = g.gameHashes[hash] + 1
	return nil
}

func (g *Game) ComputeResult() ResultData {
	return g.resultAnalyzer.computeResult(g.state)
}
