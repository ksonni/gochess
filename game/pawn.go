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

// Conforms to promotablePiece, Piece
type Pawn struct {
	pieceProps
}

func (p Pawn) Color() PieceColor {
	return p.PieceColor
}
func (p Pawn) String() string {
	return "p"
}

func (p Pawn) canMove(from Square, to Square, g *Game) bool {
	_, ok := p.computePawnMovement(from, to, g)
	return ok
}

func (p Pawn) move(from Square, to Square, g *Game) (*Board, error) {
	return p.moveAndPromote(from, to, nil, g)
}

func (p Pawn) moveAndPromote(from Square, to Square, promotion Piece, g *Game) (*Board, error) {
	movement, ok := p.computePawnMovement(from, to, g)
	if !ok {
		return nil, fmt.Errorf("pawn: illegal move")
	}

	board := g.Board().Clone()
	board.JumpPiece(from, to)

	if movement.mustPromote {
		if promotion != nil && p.canPromoteTo(promotion) {
			board.SetPiece(promotion, to)
		}
		return nil, fmt.Errorf("pawn: promotion required, normal movement not possible")
	}
	if movement.secondaryCapture != nil {
		board.ClearSquare(*movement.secondaryCapture)
	}
	return &board, nil
}

// Helpers

type pawnMovement struct {
	secondaryCapture *Square
	mustPromote      bool
}

func (p Pawn) computePawnMovement(from Square, to Square, g *Game) (*pawnMovement, bool) {
	defaultMovement := pawnMovement{mustPromote: p.promotionRank(g) == to.Rank}
	movement := new(pawnMovement)
	board := g.Board()

	if nonAttacking := p.computeNonAttackingMovableSquares(from, g); nonAttacking[to] {
		movement = &defaultMovement
	} else if _, exists := board.GetPiece(to); exists {
		if p.computeNormalAttackedSquares(from, g)[to] {
			movement = &defaultMovement
		}
	} else {
		attacked := p.computeEnPassantAttackedSquares(from, g)
		target := to.Subtracting(pawnDelta[p.Color()])
		if attacked[target] {
			defaultMovement.secondaryCapture = &target
			movement = &defaultMovement
		}
	}
	return movement, movement != nil
}

func (p Pawn) computeNonAttackingMovableSquares(from Square, g *Game) map[Square]bool {
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

func (p Pawn) computeAttackedSquares(from Square, g *Game) map[Square]bool {
	attacked := p.computeNormalAttackedSquares(from, g)
	for sq, val := range p.computeEnPassantAttackedSquares(from, g) {
		attacked[sq] = val
	}
	return attacked
}

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
	if previousPiece, exists := g.PieceAtMove(g.NumMoves()-1, homeSquare); exists {
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
