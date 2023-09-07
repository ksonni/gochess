package main

import (
	"fmt"
	"gochess/game"

	"golang.org/x/exp/slices"
)

func main() {
	printAttacked("Rook", "a1", "e4")
	printAttacked("Dark bishop", "c1", "f4")
	printAttacked("Queen", "d1", "e4")
	printAttacked("King", "e1", "e7")
	printAttacked("Light bishop", "f1", "c4")
	printAttacked("Knight", "b1", "b1")
	printAttacked("Knight", "b1", "e5")
	printAttacked("Knight", "b1", "c8")

	// TODO
	// printAttacked("Pawn", "d2", "d4")
}

func printAttacked(title string, piecePos string, pieceTargetPos string) {
	g := game.NewGame()

	sourceSq := game.MustSquare(piecePos)
	targetSq := game.MustSquare(pieceTargetPos)

	g.Board.JumpPiece(sourceSq, targetSq)

	piece := g.Board.GetPiece(targetSq)

	attacked := piece.ComputeAttackedSquares(targetSq, g)

	var out []string
	for sq, _ := range attacked {
		out = append(out, sq.String())
	}
	slices.Sort(out)

	fmt.Printf("%s from %v: %v\n", title, targetSq, out)
}
