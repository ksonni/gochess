package main

import (
	"fmt"
	"gochess/game"
)

func main() {
	g := game.NewGameState()
	moves := g.PlanPossibleMovesForSide(game.PieceColor_White)
	fmt.Printf("Game initialized. There are %d moves white can play\n", len(moves))
}
