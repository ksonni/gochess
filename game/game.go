package game

import "fmt"

type Game struct {
	position *Position
	initializer
	boardAnalyzer
	numMoves int
}

type Move struct {
	From      Square
	To        Square
	Promotion Piece
}

func NewGame() *Game {
	game := new(Game)
	board := make(Board)

	game.position = &Position{board: &board}
	game.initializePieces(game.Board())

	return game
}

func (g *Game) Board() *Board {
	return g.position.board
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

func (g *Game) NumMoves() int {
	return g.numMoves
}

func (g *Game) SideCanMove(color PieceColor) bool {
	return g.NumMoves()%2 == int(color)
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
	g.position = g.position.Appending(board)
	g.numMoves += 1
	return nil
}

func (g *Game) Position() *Position {
	p := g.position.Clone() // Defensive copy
	return &p
}
