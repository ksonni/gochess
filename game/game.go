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

func (g *Game) Move(from Square, to Square) error {
	return g.move(from, to, nil)
}

func (g *Game) MoveAndPromote(from Square, to Square, promotion Piece) error {
	return g.move(from, to, promotion)
}

func (g *Game) ComputeAttackedSquares(from Square) map[Square]bool {
	piece, exists := g.Board().GetPiece(from)
	if !exists {
		return make(map[Square]bool)
	}
	return piece.computeAttackedSquares(from, g)
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

func (g *Game) move(from Square, to Square, promotion Piece) error {
	piece, exists := g.Board().GetPiece(from)
	if !exists {
		return fmt.Errorf("game: no piece exists at %s", from)
	}
	if !piece.canMove(from, to, g) {
		return fmt.Errorf("game: invalid move - %s to %s", from, to)
	}

	var nextPos *Board
	var err error

	if promotion == nil {
		nextPos, err = piece.move(from, to, g)
	} else if promoPiece, ok := piece.(promotablePiece); ok {
		nextPos, err = promoPiece.moveAndPromote(from, to, promotion, g)
	} else {
		return fmt.Errorf("game: piece at %s doesn't support promotion", from)
	}
	if err != nil {
		return fmt.Errorf("game: move failed: %v", err)
	}

	g.positions = append(g.positions, nextPos)

	//  TODO: compute game result

	return nil
}
