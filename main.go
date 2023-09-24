package main

import (
	"fmt"
	"gochess/game"
)

func main() {
	g := game.NewGame()
	moves := g.PlanPossibleMovesForSide(game.PieceColor_Black)
	fmt.Printf("Game initialized. There are %d moves white can play\n", len(moves))
}
