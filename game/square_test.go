package game

import (
	"github.com/samber/lo"
	"golang.org/x/exp/slices"
	"reflect"
	"testing"
)

// Helpers

func assertMovePlanSquares(title string, got []MovePlan, want []string, t *testing.T) {
	var out []string
	for _, move := range got {
		out = append(out, move.To.String())
	}
	assertSquaresEqual(title, out, want, t)
}

func assertSquareMapEquals(title string, got map[Square]bool, want []string, t *testing.T) {
	var out []string
	for sq := range got {
		out = append(out, sq.String())
	}
	assertSquaresEqual(title, out, want, t)
}

func assertSquaresEqual(title string, got []string, want []string, t *testing.T) {
	got = lo.Uniq(got)
	slices.Sort(got)

	want = lo.Uniq(want)
	slices.Sort(want)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("%s: square maps: %v, want %v", title, got, want)
	}
}

func assertSquareEmpty(title string, square string, board *Board, t *testing.T) {
	if piece, exists := board.GetPiece(sq(square)); exists {
		t.Errorf("%s: found a piece on %s when none should be present: %v", title, square, piece)
	}
}

// Shorthand for MustSquare
var sq = MustSquare
