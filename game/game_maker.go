package game

func NewGameState() *GameState {
	game := new(GameState)
	game.board = NewBoard()
	game.castlingSquares = make(map[Square]SquareMovementStatus)
	initializePieces(game.Board(), game.castlingSquares)
	return game
}

func initializePieces(board *Board, castlingSquares map[Square]SquareMovementStatus) {
	whitePiecesRank := 0
	whitePawnsRank := whitePiecesRank + 1
	populatePawns(board, PieceColor_White, whitePawnsRank)
	populatePieces(board, PieceColor_White, whitePiecesRank, castlingSquares)

	blackPawnsRank := board.NumRanks() - 1 - whitePawnsRank
	blackPiecesRank := board.NumRanks() - 1 - whitePiecesRank
	populatePawns(board, PieceColor_Black, blackPawnsRank)
	populatePieces(board, PieceColor_Black, blackPiecesRank, castlingSquares)
}

func populatePieces(board *Board, color PieceColor, rank int, castlingSquares map[Square]SquareMovementStatus) {
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

func populatePawns(board *Board, color PieceColor, rank int) {
	for file := 0; file < board.NumFiles(); file++ {
		board.setPiece(
			NewPawn(color),
			Square{File: file, Rank: rank},
		)
	}
}
