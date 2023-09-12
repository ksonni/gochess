package game

// Initializes the board to start a new game
type initializer struct{}

func (i *initializer) initializePieces(board *Board) {
	whitePiecesRank := 0
	whitePawnsRank := whitePiecesRank + 1
	i.populatePawns(board, PieceColor_White, whitePawnsRank)
	i.populatePieces(board, PieceColor_White, whitePiecesRank)

	blackPawnsRank := board.NumRanks() - 1 - whitePawnsRank
	blackPiecesRank := board.NumRanks() - 1 - whitePiecesRank
	i.populatePawns(board, PieceColor_Black, blackPawnsRank)
	i.populatePieces(board, PieceColor_Black, blackPiecesRank)
}

func (i *initializer) populatePieces(board *Board, color PieceColor, rank int) {
	pieces := []Piece{
		Rook{PieceProps: NewPieceProps(color)},
		Knight{PieceProps: NewPieceProps(color)},
		Bishop{PieceProps: NewPieceProps(color)},
		Queen{PieceProps: NewPieceProps(color)},
		King{PieceProps: NewPieceProps(color)},
		Bishop{PieceProps: NewPieceProps(color)},
		Knight{PieceProps: NewPieceProps(color)},
		Rook{PieceProps: NewPieceProps(color)},
	}
	for file, piece := range pieces {
		board.SetPiece(piece, Square{File: file, Rank: rank})
	}
}

func (i *initializer) populatePawns(board *Board, color PieceColor, rank int) {
	for file := 0; file < board.NumFiles(); file++ {
		board.SetPiece(
			Pawn{PieceProps: NewPieceProps(color)},
			Square{File: file, Rank: rank},
		)
	}
}
