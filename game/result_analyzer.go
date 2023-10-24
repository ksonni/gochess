package game

import "fmt"

type resultAnalyzer struct{}

type Result int

const (
	GameResult_Active Result = iota
	GameResult_Draw
	GameResult_Checkmate
)

func (r Result) String() string {
	switch r {
	case GameResult_Active:
		return "Active"
	case GameResult_Draw:
		return "Draw"
	case GameResult_Checkmate:
		return "Checkmate"
	default:
		return fmt.Sprintf("%d", int(r))
	}
}

type DrawReason int

const (
	DrawReason_Stalemate DrawReason = iota
	DrawReason_3FoldRepetition
	DrawReason_InusfficientMaterial
	DrawReason_50Moves
)

func (r DrawReason) String() string {
	switch r {
	case DrawReason_Stalemate:
		return "Stalemate"
	case DrawReason_3FoldRepetition:
		return "3-fold repetition"
	case DrawReason_InusfficientMaterial:
		return "Insufficient material"
	case DrawReason_50Moves:
		return "50 moves without capture or pawn advance"
	default:
		return fmt.Sprintf("%d", int(r))
	}
}

type ResultData struct {
	Result     Result
	DrawReason *DrawReason
}

func (a *resultAnalyzer) isKingInCheck(color PieceColor, attackMap map[Square]bool, g *Game) bool {
	square, ok := g.Board().GetKingSquare(color)
	if !ok {
		return false
	}
	return attackMap[*square]
}

func (a *resultAnalyzer) computeResult(g *Game) ResultData {
	side := g.MovingSide()
	result := GameResult_Active
	drawReason := DrawReason_Stalemate

	if possibleMoves := g.PlanPossibleMovesForSide(side); len(possibleMoves) == 0 {
		if g.IsSideInCheck(side) {
			result = GameResult_Checkmate
		} else {
			result = GameResult_Draw
			drawReason = DrawReason_Stalemate
		}
	} else if reason, isDraw := a.testForDraw(g, side); isDraw {
		result = GameResult_Draw
		drawReason = *reason
	}

	out := ResultData{Result: result}
	if result == GameResult_Draw {
		out.DrawReason = &drawReason
	}
	return out
}

func (a *resultAnalyzer) testForDraw(g *Game, color PieceColor) (*DrawReason, bool) {
	drawReason := DrawReason_InusfficientMaterial
	if a.hasInsufficientMaterial(g, color) {
		drawReason = DrawReason_InusfficientMaterial
	} else if a.hasReached3FoldRepetition(g, color) {
		drawReason = DrawReason_3FoldRepetition
	} else if a.qualifiesFor50MoveRule(g, color) {
		drawReason = DrawReason_50Moves
	} else {
		return nil, false
	}
	return &drawReason, true
}

func (a *resultAnalyzer) hasInsufficientMaterial(g *Game, color PieceColor) bool {
	counts := g.CountPieces()
	black, white := counts[PieceColor_Black], counts[PieceColor_White]
	nBlack, nWhite := len(black), len(white)
	// Sanity check
	if black[PieceType_King] != 1 || white[PieceType_King] != 1 {
		return false
	}
	// King vs king
	if nBlack == 1 && nWhite == 1 {
		return true
	}
	// King vs King & (Knight/Bishop)
	if nBlack+nWhite == 3 {
		return black[PieceType_Knight] == 1 || black[PieceType_Bishop] == 1 ||
			white[PieceType_Knight] == 1 || white[PieceType_Bishop] == 1
	}
	// King vs King with both sides having an opposite color bishop
	if nBlack == 2 && nWhite == 2 {
		var bishopSqs []Square
		for square, piece := range g.Board().pieces {
			if piece.Type() == PieceType_Bishop {
				bishopSqs = append(bishopSqs, square)
			}
		}
		return len(bishopSqs) == 2 && bishopSqs[0].Color() != bishopSqs[1].Color()
	}
	return false
}

/**
* Position is considered the same, if all of the following are true (cf. https://www.fide.com/FIDE/handbook/LawsOfChess.pdf clause 9.2):
*
* 1. the same player has the move
* 2. pieces of the same kind and colour occupy the same squares
* 3. the possible moves of all the pieces of both players are the same, including:
* 3.1. castling rights
* 3.2. en passant captures
*
* Note: When a king or a rook is forced to move, it will lose its castling rights, if any, only after it is moved. (This should be obvious, but apparently it is not.)
*
* This means that our hash has the following inputs:
* - side to move (e.g. "White")
* - piece positions
* - castling rights for White and for Black ("White: yes/no, Black: yes/no")
* - en passant square (e.g. "a3", or "-" if none)
*
* The most straightforward way to implement this seems to have a string representation of the board, and then hash that.
 */
func (a *resultAnalyzer) hasReached3FoldRepetition(g *Game, color PieceColor) bool {
	// TODO/implement
	return false
}

/**
 * Last 50 consecutive moves have been made by each player without:
 * 1. pawn move
 * 2. capture
 *
 * (cf. https://www.fide.com/FIDE/handbook/LawsOfChess.pdf clause 9.3)
 *
 * We will add a counter to the game state, and increment it by 1 after each move. When it reaches 100, the game is a draw.
 * Whenever a pawn is moved or a capture is made, the counter is reset to 0.
 */
func (a *resultAnalyzer) qualifiesFor50MoveRule(g *Game, color PieceColor) bool {
	return g.numMovesWithoutCaptureNorPawnAdvance >= 100
}
