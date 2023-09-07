package game

import "fmt"

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
	board, err := g.BoardAtMove(g.NumMoves())
	if err != nil {
		panic(fmt.Sprintf("game: initialized without a board: %v", err))
	}
	return board
}

func (g *Game) CanMove(from Square, to Square) bool {
	piece, exists := g.Board().GetPiece(from)
	return exists && piece.canMove(from, to, g)
}

func (g *Game) ComputeAttackedSquares(from Square) map[Square]bool {
	piece, exists := g.Board().GetPiece(from)
	if !exists {
		return make(map[Square]bool)
	}
	return piece.computeAttackedSquares(from, g)
}

func (g *Game) Move(from Square, to Square) error {
	piece, exists := g.Board().GetPiece(from)
	if !exists {
		return fmt.Errorf("game: no piece exists at %s", from)
	}
	if !g.CanMove(from, to) {
		return fmt.Errorf("game: invalid move - %s to %s", from, to)
	}
	result := piece.move(from, to, g)
	g.positions = append(g.positions, result)
	return nil
}

// Moves use a 1 based index because move 0 is a valid position
func (g *Game) BoardAtMove(move int) (*Board, error) {
	size := len(g.positions)
	if move < 0 || move > size {
		return nil, fmt.Errorf("game: no board found for move %d", move)
	}
	return g.positions[move], nil
}

func (g *Game) NumMoves() int {
	return len(g.positions) - 1
}
