package game

import "fmt"

type Game struct {
	position *Position
	numMoves int
	boardAnalyzer
}

type Move struct {
	From      Square
	To        Square
	Promotion Piece
}

func (g *Game) Board() *Board {
	return g.position.board
}

func (g *Game) Move(move Move) error {
	result, err := g.WithMove(move)
	if err != nil {
		return err
	}
	g.numMoves = result.numMoves
	g.position = result.position
	return nil
}

func (g *Game) WithMove(move Move) (*Game, error) {
	piece, exists := g.Board().GetPiece(move.From)
	if !exists {
		return nil, fmt.Errorf("game: no piece exists at %s", move.From)
	}
	if !g.SideCanMove(piece.Color()) {
		return nil, fmt.Errorf("game: attempted to move piece out of turn")
	}

	nextPos, err := piece.WithLocalMove(move, g)
	if err != nil {
		return nil, fmt.Errorf("game: move failed: %v", err)
	}
	if nextPos.IsSideInCheck(piece.Color()) {
		return nil, fmt.Errorf("game: move failed: violates king integrity")
	}

	return nextPos, nil
}

func (g *Game) ComputeSquaresAttackedBySide(color PieceColor) map[Square]bool {
	return g.boardAnalyzer.computeSquaresAttackedBySide(color, g)
}

func (g *Game) ComputeAttackedSquares(from Square) map[Square]bool {
	piece, exists := g.Board().GetPiece(from)
	if !exists {
		return make(map[Square]bool)
	}
	return piece.ComputeAttackedSquares(from, g)
}

func (g *Game) IsSideInCheck(color PieceColor) bool {
	attackMap := g.ComputeSquaresAttackedBySide(color.Opponent())
	return g.boardAnalyzer.isSideInCheck(color, attackMap, g.Board(), g)
}

func (g *Game) NumMoves() int {
	return g.numMoves
}

func (g *Game) SideCanMove(color PieceColor) bool {
	return g.NumMoves()%2 == int(color)
}

func (g *Game) Position() *Position {
	p := g.position.Clone() // Defensive copy
	return &p
}

// Helpers

func (g *Game) appendingPosition(board *Board) *Game {
	return &Game{
		position: g.position.Appending(board),
		numMoves: g.numMoves + 1,
	}
}

func (g *Game) withPosition(board *Board) *Game {
	return &Game{
		position: g.position.Setting(board),
		numMoves: g.numMoves + 1,
	}
}
