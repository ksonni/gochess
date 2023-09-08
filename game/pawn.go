package game

var (
	pawnDelta           = make(map[PieceColor]Square)
	pawnHomeDelta       = make(map[PieceColor]Square)
	pawnAttackDeltas    = make(map[PieceColor]map[Square]bool)
	pawnEmPassantDeltas = map[Square]bool{
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

// Conforms to promotablePiece
type Pawn struct {
	pieceProps
}

func (p Pawn) Color() PieceColor {
	return p.PieceColor
}
func (p Pawn) String() string {
	return "p"
}

// TODO: implement properly
func (p Pawn) move(from Square, to Square, g *Game) (*Board, error) {
	b := g.Board().Clone()
	b.JumpPiece(from, to)
	return &b, nil
}

// TODO: implement properly
func (p Pawn) canMove(from Square, to Square, g *Game) bool {
	// panic("pawn: not implemented")
	return true
}

func (p Pawn) computeAttackedSquares(from Square, g *Game) map[Square]bool {
	attacked := p.computeNormalAttackedSquares(from, g)
	for sq, val := range p.computeEmPassantAttackedSquares(from, g) {
		attacked[sq] = val
	}
	return attacked
}

// TODO
func (p Pawn) moveAndPromote(from Square, to Square, promotion Piece, g *Game) (*Board, error) {
	panic("pawn: promotion not implemented")
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

func (p Pawn) attacksEmPassant(from Square, to Square, g *Game) bool {
	board := g.Board()
	if !board.SquareInRange(to) {
		return false
	}
	// Square must be horizontally adjacent
	if _, ok := pawnEmPassantDeltas[from.Subtracting(to)]; !ok {
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

func (p Pawn) computeEmPassantAttackedSquares(from Square, g *Game) map[Square]bool {
	attacked := make(map[Square]bool)
	for delta := range pawnEmPassantDeltas {
		to := from.Adding(delta)
		if p.attacksEmPassant(from, to, g) {
			attacked[to] = true
		}
	}
	return attacked
}
