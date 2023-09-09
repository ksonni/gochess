package game

import "fmt"

type Game struct {
	positions []*Board
	initializer
	boardAnalyzer
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

func (g *Game) PlanMove(from Square, to Square) (*Board, error) {
	return g.planMove(from, to, nil)
}

func (g *Game) PlanMoveWithPromotionLocally(from Square, to Square, promotion Piece) (*Board, error) {
	return g.planMove(from, to, promotion)
}

func (g *Game) Move(from Square, to Square) error {
	return g.move(from, to, nil)
}

func (g *Game) MoveWithPromotion(from Square, to Square, promotion Piece) error {
	return g.move(from, to, promotion)
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

// Helpers

func (g *Game) planMove(from Square, to Square, promotion Piece) (*Board, error) {
	piece, exists := g.Board().GetPiece(from)
	if !exists {
		return nil, fmt.Errorf("game: no piece exists at %s", from)
	}
	if !g.SideCanMove(piece.Color()) {
		return nil, fmt.Errorf("game: attempted to move piece out of turn")
	}

	var nextPos *Board
	var err error

	if promotion == nil {
		nextPos, err = piece.PlanMoveLocally(from, to, g)
	} else if promoPiece, ok := piece.(PromotablePiece); ok {
		nextPos, err = promoPiece.PlanMoveWithPromotionLocally(from, to, promotion, g)
	} else {
		return nil, fmt.Errorf("game: piece at %s doesn't support promotion", from)
	}
	if err != nil {
		return nil, fmt.Errorf("game: move failed: %v", err)
	}
	if g.IsSideInCheck(piece.Color(), nextPos) {
		return nil, fmt.Errorf("game: move failed: violates king integrity")
	}

	return nextPos, nil
}

func (g *Game) move(from Square, to Square, promotion Piece) error {
	board, err := g.planMove(from, to, promotion)
	if err != nil {
		return err
	}
	g.positions = append(g.positions, board)
	return nil
}
