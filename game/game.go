package game

import "fmt"

type Game struct {
	positions []*Board
	initializer
	boardAnalyzer
}

type Move struct {
	From      Square
	To        Square
	Promotion Piece
}

func NewGame() *Game {
	game := new(Game)
	board := make(Board)

	game.positions = []*Board{&board}
	game.initializePieces(game.Board())

	return game
}

func (g *Game) Board() *Board {
	board, exists := g.BoardAtMove(g.NumMoves())
	if !exists {
		panic("game: initialized without a board")
	}
	return board
}

func (g *Game) PreviousBoard() (*Board, bool) {
	return g.BoardAtMove(g.NumMoves() - 1)
}

func (g *Game) ComputeSquaresAttackedBySide(color PieceColor, board *Board) map[Square]bool {
	return g.boardAnalyzer.computeSquaresAttackedBySide(color, board, g)
}

func (g *Game) ComputeAttackedSquares(from Square) map[Square]bool {
	piece, exists := g.Board().GetPiece(from)
	if !exists {
		return make(map[Square]bool)
	}
	return piece.ComputeAttackedSquares(from, g)
}

func (g *Game) IsSideInCheck(color PieceColor, board *Board) bool {
	attackMap := g.ComputeSquaresAttackedBySide(color.Opponent(), board)
	return g.boardAnalyzer.isSideInCheck(color, attackMap, board, g)
}

// Moves use a 1 based index because move 0 is a valid position
func (g *Game) BoardAtMove(move int) (*Board, bool) {
	size := len(g.positions)
	if move < 0 || move >= size {
		return nil, false
	}
	return g.positions[move], true
}

func (g *Game) PieceAtMove(move int, square Square) (Piece, bool) {
	board, ok := g.BoardAtMove(move)
	if !ok {
		return nil, false
	}
	return board.GetPiece(square)
}

func (g *Game) PiecePositionAtMove(move int, piece Piece) (*Square, bool) {
	board, ok := g.BoardAtMove(move)
	if !ok {
		return nil, false
	}
	for square, p := range *board {
		if p.Id() == piece.Id() {
			return &square, true
		}
	}
	return nil, false
}

func (g *Game) NumMoves() int {
	return len(g.positions) - 1
}

func (g *Game) SideCanMove(color PieceColor) bool {
	return g.NumMoves()%2 == int(color)
}

func (g *Game) SquareHasSamePieceAtMoves(sq Square, move1 int, move2 int) bool {
	if move1 < 0 || move2 < 0 || move1 > g.NumMoves() || move2 > g.NumMoves() {
		return false
	}
	p1, _ := g.positions[move1].GetPiece(sq)
	p2, _ := g.positions[move2].GetPiece(sq)
	if p1 == nil && p2 == nil {
		return true
	}
	if p1 == nil || p2 == nil {
		return false
	}
	return p1.Id() == p2.Id()
}

func (g *Game) SquareHasChangedSinceMove(sq Square, move int) bool {
	current := g.NumMoves()
	for i := move; i < len(g.positions); i++ {
		if !g.SquareHasSamePieceAtMoves(sq, i, current) {
			return true
		}
	}
	return false
}

func (g *Game) PlanMove(move Move) (*Board, error) {
	piece, exists := g.Board().GetPiece(move.From)
	if !exists {
		return nil, fmt.Errorf("game: no piece exists at %s", move.From)
	}
	if !g.SideCanMove(piece.Color()) {
		return nil, fmt.Errorf("game: attempted to move piece out of turn")
	}

	nextPos, err := piece.PlanMoveLocally(move, g)
	if err != nil {
		return nil, fmt.Errorf("game: move failed: %v", err)
	}
	if g.IsSideInCheck(piece.Color(), nextPos) {
		return nil, fmt.Errorf("game: move failed: violates king integrity")
	}

	return nextPos, nil
}

func (g *Game) Move(move Move) error {
	board, err := g.PlanMove(move)
	if err != nil {
		return err
	}
	g.positions = append(g.positions, board)
	return nil
}
