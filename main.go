package main

import (
	"fmt"
	"gochess/game"
)

func main() {
	fmt.Printf("Game initialized at move %d\n", game.NewGame().NumMoves())
}
