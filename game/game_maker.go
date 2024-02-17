package game

// Initializes the board to start a new game
type gameMaker struct{}

func NewGame() *Game {
	game := new(Game)
	game.board = NewBoard()
	game.castlingSquares = make(map[Square]SquareMovementStatus)
	i := gameMaker{}
	i.initializePieces(game.Board(), game.castlingSquares)
	return game
}

func (i *gameMaker) initializePieces(board *Board, castlingSquares map[Square]SquareMovementStatus) {
	whitePiecesRank := 0
	whitePawnsRank := whitePiecesRank + 1
	i.populatePawns(board, PieceColor_White, whitePawnsRank)
	i.populatePieces(board, PieceColor_White, whitePiecesRank, castlingSquares)

	blackPawnsRank := board.NumRanks() - 1 - whitePawnsRank
	blackPiecesRank := board.NumRanks() - 1 - whitePiecesRank
	i.populatePawns(board, PieceColor_Black, blackPawnsRank)
	i.populatePieces(board, PieceColor_Black, blackPiecesRank, castlingSquares)
}

func (i *gameMaker) populatePieces(board *Board, color PieceColor, rank int, castlingSquares map[Square]SquareMovementStatus) {
	pieces := []Piece{
		NewRook(color),
		NewKnight(color),
		NewBishop(color),
		NewQueen(color),
		NewKing(color),
		NewBishop(color),
		NewKnight(color),
		NewRook(color),
	}
	for file, piece := range pieces {
		t := piece.Type()
		sq := Square{File: file, Rank: rank}
		if t == PieceType_Rook || t == PieceType_King {
			castlingSquares[sq] = SquareMovementStatus_Unmoved
		}
		board.setPiece(piece, sq)
	}
}

func (i *gameMaker) populatePawns(board *Board, color PieceColor, rank int) {
	for file := 0; file < board.NumFiles(); file++ {
		board.setPiece(
			NewPawn(color),
			Square{File: file, Rank: rank},
		)
	}
}
