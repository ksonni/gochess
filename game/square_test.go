package game

import (
	"reflect"
	"testing"

	"golang.org/x/exp/slices"
)

// Helpers

func assertAttackedSquaresEqual(title string, got map[Square]bool, want []string, t *testing.T) {
	var out []string
	for sq := range got {
		out = append(out, sq.String())
	}
	slices.Sort(out)
	slices.Sort(want)

	if !reflect.DeepEqual(out, want) {
		t.Errorf("%s: attacked squares: %v, want %v", title, out, want)
	}
}

// Shorthand for MustSquare
var sq = MustSquare
