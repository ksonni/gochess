package game

import "fmt"

var (
	pawnDelta           = make(map[PieceColor]Square)
	pawnHomeDelta       = make(map[PieceColor]Square)
	pawnAttackDeltas    = make(map[PieceColor]map[Square]bool)
	pawnEnPassantDeltas = map[Square]bool{
		{File: -1, Rank: 0}: true,
		{File: 1, Rank: 0}:  true,
	}
)

func init() {
	setPawnDeltas := func(color PieceColor, multiplier Square) {
		pawnDelta[color] = Square{File: 0, Rank: 1}.Multiplying(multiplier)
		pawnHomeDelta[color] = Square{File: 0, Rank: 2}.Multiplying(multiplier)
		pawnAttackDeltas[color] = map[Square]bool{
			Square{File: 1, Rank: 1}.Multiplying(multiplier):  true,
			Square{File: -1, Rank: 1}.Multiplying(multiplier): true,
		}
	}
	setPawnDeltas(PieceColor_White, Square{File: 1, Rank: 1})
	setPawnDeltas(PieceColor_Black, Square{File: 1, Rank: -1})
}

func NewPawn(color PieceColor) Pawn {
	return Pawn{pieceProps: newPieceProps(color)}
}

// Pawn Conforms to promotablePiece, Piece
type Pawn struct {
	pieceProps
}

func (p Pawn) Color() PieceColor {
	return p.PieceColor
}
func (p Pawn) String() string {
	if p.PieceColor == PieceColor_White {
		return "P"
	} else {
		return "p"
	}
}
func (p Pawn) Type() PieceType {
	return PieceType_Pawn
}
func (p Pawn) ComputeAttackedSquares(from Square, g *Game) map[Square]bool {
	attacked := p.computeNormalAttackedSquares(from, g)
	for sq, val := range p.computeEnPassantAttackedSquares(from, g) {
		attacked[sq] = val
	}
	return attacked
}

func (p Pawn) WithLocalMove(move Move, g *Game) (*Game, error) {
	from, to := move.From, move.To
	movement, ok := p.computePawnMovements(from, g)[to]
	if !ok {
		return nil, fmt.Errorf("pawn: illegal move")
	}
	return p.planPawnMovement(move, movement, g)
}

func (p Pawn) PlanPossibleMovesLocally(from Square, g *Game) []MovePlan {
	var plans []MovePlan
	for to, movement := range p.computePawnMovements(from, g) {
		var moves []Move
		if movement.mustPromote {
			for _, piece := range p.promotablePieces() {
				moves = append(moves, Move{From: from, To: to, Promotion: piece})
			}
		} else {
			moves = append(moves, Move{From: from, To: to})
		}
		for _, move := range moves {
			if res, err := p.planPawnMovement(move, movement, g); err == nil {
				plans = append(plans, MovePlan{Move: move, Game: res})
			}
		}
	}
	return plans
}

// Helpers

type pawnMovement struct {
	secondaryCapture *Square
	mustPromote      bool
}

func (p Pawn) planPawnMovement(move Move, movement *pawnMovement, g *Game) (*Game, error) {
	from, to, promotion := move.From, move.To, move.Promotion
	board := g.Board().Clone()
	board.jumpPiece(from, to)

	if movement.mustPromote {
		if promotion == nil || !p.canPromoteTo(promotion) {
			return nil, fmt.Errorf("pawn: either no piece or an invalid one has been provided for promotion")
		}
		board.setPiece(promotion, to)
	}
	if movement.secondaryCapture != nil {
		board.clearSquare(*movement.secondaryCapture)
	}
	newGame := g.appendingPosition(board)
	newGame.numMovesWithoutCaptureNorPawnAdvance = 0
	return newGame, nil
}

func (p Pawn) computePawnMovements(from Square, g *Game) map[Square]*pawnMovement {
	out := make(map[Square]*pawnMovement)
	for sq, movement := range p.computeNonAttackedMovements(from, g) {
		out[sq] = movement
	}
	for sq, movement := range p.computeNormalAttackedMovements(from, g) {
		out[sq] = movement
	}
	for sq, movement := range p.computeEnPassantAttackedMovements(from, g) {
		out[sq] = movement
	}
	return out
}

// Non attacking movement

func (p Pawn) computeNonAttackedMovableSquares(from Square, g *Game) map[Square]bool {
	movable := make(map[Square]bool)
	delta := pawnDelta[p.Color()]
	board := g.Board()

	m1 := from.Adding(delta)
	if _, exists := board.GetPiece(m1); !exists {
		movable[m1] = board.SquareInRange(m1)
	}
	m2 := m1.Adding(delta)
	if _, exists := board.GetPiece(m2); !exists && movable[m1] && p.homeRank(g) == from.Rank {
		movable[m2] = board.SquareInRange(m2)
	}
	return movable
}

func (p Pawn) computeNonAttackedMovements(from Square, g *Game) map[Square]*pawnMovement {
	out := make(map[Square]*pawnMovement)
	nonAttacked := p.computeNonAttackedMovableSquares(from, g)
	for to := range nonAttacked {
		out[to] = &pawnMovement{mustPromote: p.promotionRank(g) == to.Rank}
	}
	return out
}

// Normal attacking movement

func (p Pawn) attacksNormal(from Square, to Square, g *Game) bool {
	board := g.Board()
	if !board.SquareInRange(to) {
		return false
	}
	// Square must be in attacking range
	if _, ok := pawnAttackDeltas[p.Color()][to.Subtracting(from)]; !ok {
		return false
	}
	// Must be an enemy piece
	piece, exists := board.GetPiece(to)
	return !exists || piece.Color() != p.Color()
}

func (p Pawn) computeNormalAttackedSquares(from Square, g *Game) map[Square]bool {
	attackDeltas := pawnAttackDeltas[p.Color()]
	attacked := make(map[Square]bool)
	for delta := range attackDeltas {
		to := from.Adding(delta)
		if p.attacksNormal(from, to, g) {
			attacked[to] = true
		}
	}
	return attacked
}

func (p Pawn) computeNormalAttackedMovements(from Square, g *Game) map[Square]*pawnMovement {
	out := make(map[Square]*pawnMovement)
	attacked := p.computeNormalAttackedSquares(from, g)
	for to := range attacked {
		if _, exists := g.Board().GetPiece(to); exists {
			out[to] = &pawnMovement{mustPromote: p.promotionRank(g) == to.Rank}
		}
	}
	return out
}

// En-Passant movement

func (p Pawn) attacksEnPassant(from Square, to Square, g *Game) bool {
	board := g.Board()
	if !board.SquareInRange(to) {
		return false
	}
	// Square must be adjacent on the same rank
	if _, ok := pawnEnPassantDeltas[from.Subtracting(to)]; !ok {
		return false
	}
	// Square must have an enemy piece
	target, targetExists := board.GetPiece(to)
	if !targetExists || target.Color() == p.Color() {
		return false
	}
	// Enemy piece must be a pawn
	if _, isPawn := target.(Pawn); !isPawn {
		return false
	}
	// Enemy pawn must have moved 2 squares in the previous move
	homeSquare := to.Subtracting(pawnHomeDelta[target.Color()])
	if previousPiece, exists := g.Position().PieceAtPreviousMove(homeSquare); exists {
		return previousPiece.Id() == target.Id()
	}
	return false
}

func (p Pawn) computeEnPassantAttackedSquares(from Square, g *Game) map[Square]bool {
	attacked := make(map[Square]bool)
	for delta := range pawnEnPassantDeltas {
		to := from.Adding(delta)
		if p.attacksEnPassant(from, to, g) {
			attacked[to] = true
		}
	}
	return attacked
}

func (p Pawn) computeEnPassantAttackedMovements(from Square, g *Game) map[Square]*pawnMovement {
	out := make(map[Square]*pawnMovement)
	attacked := p.computeEnPassantAttackedSquares(from, g)
	for attackSq := range attacked {
		to := attackSq.Adding(pawnDelta[p.Color()])
		movement := &pawnMovement{mustPromote: p.promotionRank(g) == to.Rank}
		if attacked[attackSq] {
			movement.secondaryCapture = &attackSq
		}
		out[to] = movement
	}
	return out
}

func (p Pawn) homeRank(g *Game) int {
	board := g.Board()
	switch p.Color() {
	case PieceColor_Black:
		return board.NumRanks() - 2
	case PieceColor_White:
		return 1
	}
	return 1
}

func (p Pawn) promotionRank(g *Game) int {
	board := g.Board()
	switch p.Color() {
	case PieceColor_Black:
		return 0
	case PieceColor_White:
		return board.NumRanks() - 1
	}
	return 0
}

func (p Pawn) canPromoteTo(target Piece) bool {
	if p.Color() != target.Color() {
		return false
	}
	switch target.(type) {
	case Queen, Rook, Bishop, Knight:
		return true
	default:
		return false
	}
}

func (p Pawn) promotablePieces() []Piece {
	return []Piece{
		NewQueen(p.Color()),
		NewRook(p.Color()),
		NewKnight(p.Color()),
		NewBishop(p.Color()),
	}
}
