package game

import (
	"crypto/md5"
	"fmt"
)

type Game struct {
	position                             *Position
	numMoves                             int
	numMovesWithoutCaptureNorPawnAdvance int
	history                              map[string]int
	threeFoldDrawReached                 bool
	resultAnalyzer
}

type Move struct {
	From      Square
	To        Square
	Promotion Piece
}

type MovePlan struct {
	Move
	Game *Game
}

func (g *Game) Board() *Board {
	return g.position.board
}

func (g *Game) Move(move Move) error {
	if _, exists := g.history[g.hash3Fold()]; !exists {
		// This is the first move, and we need to initialize history
		// XXX Ideally the state of setting up the board would be a separate game state, and there would be Game.Start(), after which the game would always be in a valid state, but for now we will just do this
		g.history = make(map[string]int)
		g.history[g.hash3Fold()] = 1
	}
	result, err := g.WithMove(move)
	if err != nil {
		return err
	}
	g.numMoves = result.numMoves
	g.position = result.position
	g.numMovesWithoutCaptureNorPawnAdvance = result.numMovesWithoutCaptureNorPawnAdvance
	g.history = result.history
	g.threeFoldDrawReached = result.threeFoldDrawReached
	return nil
}

func (g *Game) WithMove(move Move) (*Game, error) {
	piece, exists := g.Board().GetPiece(move.From)
	if !exists {
		return nil, fmt.Errorf("game: no piece exists at %s", move.From)
	}
	if g.MovingSide() != piece.Color() {
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
	attacked := make(map[Square]bool)
	for square, piece := range g.Board().pieces {
		if piece.Color() != color {
			continue
		}
		pieceAttacked := piece.ComputeAttackedSquares(square, g)
		for pieceSquare, val := range pieceAttacked {
			attacked[pieceSquare] = val
		}
	}
	return attacked
}

func (g *Game) ComputeAttackedSquares(from Square) map[Square]bool {
	piece, exists := g.Board().GetPiece(from)
	if !exists {
		return make(map[Square]bool)
	}
	return piece.ComputeAttackedSquares(from, g)
}

func (g *Game) PlanPossibleMovesForSide(color PieceColor) []MovePlan {
	var out []MovePlan
	for sq, piece := range g.Board().pieces {
		if piece.Color() != color {
			continue
		}
		out = append(out, g.PlanPossibleMoves(sq)...)
	}
	return out
}

func (g *Game) PlanPossibleMoves(from Square) []MovePlan {
	piece, exists := g.Board().GetPiece(from)
	if !exists {
		return []MovePlan{}
	}
	moves := piece.PlanPossibleMovesLocally(from, g)
	var out []MovePlan
	color := piece.Color()
	for _, move := range moves {
		if move.Game.IsSideInCheck(color) {
			continue
		}
		out = append(out, move)
	}
	return out
}

func (g *Game) IsSideInCheck(color PieceColor) bool {
	attackMap := g.ComputeSquaresAttackedBySide(color.Opponent())
	return g.resultAnalyzer.isKingInCheck(color, attackMap, g)
}

func (g *Game) NumMoves() int {
	return g.numMoves
}

func (g *Game) MovingSide() PieceColor {
	return PieceColor(g.NumMoves() % 2)
}

func (g *Game) Position() *Position {
	p := g.position.Clone() // Defensive copy
	return &p
}

func (g *Game) CountPieces() map[PieceColor]map[PieceType]int {
	out := map[PieceColor]map[PieceType]int{
		PieceColor_White: {},
		PieceColor_Black: {},
	}
	for _, piece := range g.Board().pieces {
		out[piece.Color()][piece.Type()] += 1
	}
	return out
}

func (g *Game) String() string {
	s := g.string3Fold()
	s += fmt.Sprintln("moves:", g.numMoves)
	return s
}

// Helpers

/**
 * String representation of the game state for the purpose of the 3-fold-move draw.
 *
 * This is used in computing the 3-fold-move draw, i.e. it must capture the FIDE rules for 3-fold-move draw, and only those rules. See Game.String() for a more complete representation of the game state.
 */
func (g *Game) string3Fold() string {
	s := ""
	if g.position != nil && g.position.board != nil {
		s += g.position.board.String()
	}
	if g.MovingSide() == PieceColor_White {
		s += "White to move\n"
	} else {
		s += "Black to move\n"
	}
	whiteCanCastle, blackCanCastle := g.castlingRights()
	s += fmt.Sprintln("castling rights: White:", whiteCanCastle, "Black:", blackCanCastle)
	enPassantSquare, canEnPassant := g.enPassantSquare()
	if canEnPassant {
		s += fmt.Sprintln("en-passant square:", enPassantSquare)
	} else {
		s += fmt.Sprintln("en-passant square: -")
	}

	return s
}

/**
 * Hash representation of the game state.
 *
 * Used in computing the 3-fold-move draw, i.e. two positions will have a different hash if and only if they are different according to the 3-fold-move-draw rule
 */
func (g *Game) hash3Fold() string {
	data := g.string3Fold()
	hash := md5.Sum([]byte(data))
	return fmt.Sprintf("%x", hash)
}

func (g *Game) enPassantSquare() (Square, bool) {
	// find all pawns in the game
	for square, piece := range g.position.board.pieces {
		if pawn, ok := piece.(Pawn); ok && pawn.Type() == PieceType_Pawn {
			attackedSquares := pawn.computeEnPassantAttackedSquares(square, g)
			for x := range attackedSquares {
				if attackedSquares[x] {
					return x, true
				}
			}

		}
	}
	return Square{}, false
}

func (g *Game) castlingRights() (bool, bool) {
	whiteCanCastle := !g.position.SquareHasEverChanged(Square{4, 0}) && (!g.position.SquareHasEverChanged(Square{0, 0}) || !g.position.SquareHasEverChanged(Square{7, 0}))
	blackCanCastle := !g.position.SquareHasEverChanged(Square{4, 7}) && (!g.position.SquareHasEverChanged(Square{0, 7}) || !g.position.SquareHasEverChanged(Square{7, 7}))
	return whiteCanCastle, blackCanCastle
}

func (g *Game) appendingPosition(board *Board) *Game {
	appending := true
	return g._withPosition(board, appending)
}

func (g *Game) withPosition(board *Board) *Game {
	appending := false
	return g._withPosition(board, appending)
}

func (g *Game) _withPosition(board *Board, appending bool) *Game {
	newHistory := make(map[string]int)
	for k, v := range g.history {
		newHistory[k] = v
	}
	var myBoard *Position
	if appending {
		myBoard = g.position.Appending(board)
	} else {
		myBoard = g.position.Setting(board)
	}
	newGameState := &Game{
		position:                             myBoard,
		numMoves:                             g.numMoves + 1,
		numMovesWithoutCaptureNorPawnAdvance: g.numMovesWithoutCaptureNorPawnAdvance + 1,
		history:                              newHistory,
		threeFoldDrawReached:                 g.threeFoldDrawReached,
	}
	newHash := newGameState.hash3Fold()
	if _, exists := newGameState.history[newHash]; !exists {
		newGameState.history[newHash] = 0
	}
	newGameState.history[newHash] += 1

	if newGameState.history[newHash] > 2 {
		newGameState.threeFoldDrawReached = true
	}

	return newGameState
}
