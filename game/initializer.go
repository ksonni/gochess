package game

// Initializes the board to start a new game
type initializer struct{}

func (i *initializer) initializePieces(board Board) {
	whitePiecesRank := 0
	whitePawnsRank := whitePiecesRank + 1
	i.populatePawns(board, PieceColor_White, whitePawnsRank)
	i.populatePieces(board, PieceColor_White, whitePiecesRank)

	blackPawnsRank := board.NumRanks() - 1 - whitePawnsRank
	blackPiecesRank := board.NumRanks() - 1 - whitePiecesRank
	i.populatePawns(board, PieceColor_Black, blackPawnsRank)
	i.populatePieces(board, PieceColor_Black, blackPiecesRank)
}

func (i *initializer) populatePieces(board Board, color PieceColor, rank int) {
	pieces := []Piece{
		Rook{PieceColor: color},
		Knight{PieceColor: color},
		Bishop{PieceColor: color},
		Queen{PieceColor: color},
		King{PieceColor: color},
		Bishop{PieceColor: color},
		Knight{PieceColor: color},
		Rook{PieceColor: color},
	}
	for file, piece := range pieces {
		board.SetPiece(piece, Square{File: file, Rank: rank})
	}
}

func (i *initializer) populatePawns(board Board, color PieceColor, rank int) {
	for file := 0; file < board.NumFiles(); file++ {
		board.SetPiece(
			Pawn{PieceColor: color},
			Square{File: file, Rank: rank},
		)
	}
}
