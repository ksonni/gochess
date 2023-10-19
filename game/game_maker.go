package game

// Initializes the board to start a new game
type gameMaker struct{}

func NewGame() *Game {
	game := new(Game)
	game.position = &Position{board: NewBoard()}
	game.numMoves = 0
	game.numMovesWithoutCaptureNorPawnAdvance = 0
	i := gameMaker{}
	i.initializePieces(game.Board())
	return game
}

func (i *gameMaker) initializePieces(board *Board) {
	whitePiecesRank := 0
	whitePawnsRank := whitePiecesRank + 1
	i.populatePawns(board, PieceColor_White, whitePawnsRank)
	i.populatePieces(board, PieceColor_White, whitePiecesRank)

	blackPawnsRank := board.NumRanks() - 1 - whitePawnsRank
	blackPiecesRank := board.NumRanks() - 1 - whitePiecesRank
	i.populatePawns(board, PieceColor_Black, blackPawnsRank)
	i.populatePieces(board, PieceColor_Black, blackPiecesRank)
}

func (i *gameMaker) populatePieces(board *Board, color PieceColor, rank int) {
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
		board.setPiece(piece, Square{File: file, Rank: rank})
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
