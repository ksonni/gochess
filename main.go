package main

import (
	"fmt"
	"gochess/game"

	"golang.org/x/exp/slices"
)

func main() {
	printAttacked("Rook", "e4", move{"a1", "e4"})
	printAttacked("Dark bishop", "f4", move{"c1", "f4"})
	printAttacked("Queen", "e4", move{"d1", "e4"})
	printAttacked("King", "e7", move{"e1", "e7"})
	printAttacked("Light bishop", "c4", move{"f1", "c4"})
	printAttacked("Knight", "b1", move{"b1", "b1"})
	printAttacked("Knight", "e5", move{"b1", "e5"})
	printAttacked("Knight", "c8", move{"b1", "c8"})

	printAttacked("Pawn", "a2", move{"a2", "a2"})
	printAttacked("Pawn", "h2", move{"h2", "h2"})

	// empassant allowed
	printAttacked("Pawn", "d4",
		move{"a2", "a3"}, move{"d7", "d5"},
		move{"a3", "a4"}, move{"d5", "d4"},
		move{"e2", "e4"},
	)

	// Loses empassant rights if not taken immediately
	printAttacked("Pawn", "d4",
		move{"a2", "a3"}, move{"d7", "d5"},
		move{"a3", "a4"}, move{"d5", "d4"},
		move{"e2", "e4"}, move{"h7", "h6"},
	)
}

type move = [2]string

func printAttacked(title string, testPos string, moves ...move) {
	g := game.NewGame()

	for _, move := range moves {
		fromSq := game.MustSquare(move[0])
		targetSq := game.MustSquare(move[1])
		err := g.Move(fromSq, targetSq)
		if err != nil {
			fmt.Println("WARN: Invalid")
			g.Board().JumpPiece(fromSq, targetSq)
		}
		_, exists := g.Board().GetPiece(targetSq)
		if !exists {
			panic(fmt.Sprintf("main: no piece exists at %s", fromSq))
		}
	}

	attacked := g.ComputeAttackedSquares(game.MustSquare(testPos))

	var out []string
	for sq := range attacked {
		out = append(out, sq.String())
	}
	slices.Sort(out)

	fmt.Printf("%s from %v: %v\n", title, testPos, out)
}
