package game

import "fmt"

type Square struct {
	File int
	Rank int
}

func (square Square) String() string {
	return fmt.Sprintf("%c%d", 'a'+square.File, square.Rank+1)
}

const boardNumFiles = 8
const boardNumRanks = 8

type Board [boardNumFiles][boardNumRanks]Piece

func (board *Board) GetPiece(square Square) Piece {
	return board[square.File][square.Rank]
}

func (board *Board) SetPiece(piece Piece, square Square) {
	board[square.File][square.Rank] = piece
}

func (board *Board) ClearSquare(square Square) {
	board.SetPiece(nil, square)
}

func (board *Board) MovePiece(start Square, end Square) {
	piece := board.GetPiece(start)
	board.ClearSquare(start)
	board.SetPiece(piece, end)
}

func (board *Board) NumRanks() int {
	return boardNumRanks
}

func (board *Board) NumFiles() int {
	return boardNumFiles
}
