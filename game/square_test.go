package game

import (
	"reflect"
	"testing"

	"golang.org/x/exp/slices"
)

// Helpers

func assertSquareMapEquals(title string, got map[Square]bool, want []string, t *testing.T) {
	var out []string
	for sq := range got {
		out = append(out, sq.String())
	}
	slices.Sort(out)
	slices.Sort(want)

	if !reflect.DeepEqual(out, want) {
		t.Errorf("%s: square maps: %v, want %v", title, out, want)
	}
}

func assertSquareEmpty(title string, square string, board *Board, t *testing.T) {
	if piece, exists := board.GetPiece(sq(square)); exists {
		t.Errorf("%s: found a piece on %s when none should be present: %v", title, square, piece)
	}
}

// Shorthand for MustSquare
var sq = MustSquare
