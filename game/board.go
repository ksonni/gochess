package game

import "fmt"

type Square struct {
	File int
	Rank int
}

func (square Square) String() string {
	return fmt.Sprintf("%c%d", 'a'+square.File, square.Rank+1)
}

func (square Square) Adding(file int, rank int) Square {
	return Square{File: square.File + file, Rank: square.Rank + rank}
}

const boardNumFiles = 8
const boardNumRanks = 8

type Board [][]Piece

func NewBoard() Board {
	b := make([][]Piece, boardNumFiles)
	for i := range b {
		b[i] = make([]Piece, boardNumRanks)
	}
	return b
}

func (board Board) GetPiece(square Square) Piece {
	return board[square.File][square.Rank]
}

func (board Board) SetPiece(piece Piece, square Square) {
	board[square.File][square.Rank] = piece
}

func (board Board) ClearSquare(square Square) {
	board.SetPiece(nil, square)
}

func (board Board) MovePiece(start Square, end Square) {
	piece := board.GetPiece(start)
	board.ClearSquare(start)
	board.SetPiece(piece, end)
}

func (board Board) NumRanks() int {
	return boardNumRanks
}

func (board Board) NumFiles() int {
	return boardNumFiles
}

func (board Board) Clone() Board {
	copy := NewBoard()
	for i := range board {
		for j := range board[i] {
			copy[i][j] = board[i][j]
		}
	}
	return copy
}

func (board Board) HasSquare(square Square) bool {
	return square.File >= 0 && square.File < board.NumFiles() &&
		square.Rank >= 0 && square.Rank < board.NumRanks()
}
