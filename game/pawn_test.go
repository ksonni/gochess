package game

import "testing"

func TestPawnAttackedSquares(t *testing.T) {
	runAttackedSquaresTests(map[string]testAttackedSquaresInput{
		"Pawn left flank": {
			[]testMove{},
			"a2",
			[]string{"b3"},
			false,
		},
		"Pawn right flank": {
			[]testMove{},
			"h2",
			[]string{"g3"},
			false,
		},
		"Pawn en-passant": {
			[]testMove{
				{"a2", "a3"}, {"d7", "d5"},
				{"a3", "a4"}, {"d5", "d4"},
				{"e2", "e4"},
			},
			"d4",
			[]string{"c3", "e3", "e4"},
			false,
		},
		"Pawn en-passant right lost if not used immediately": {
			[]testMove{
				{"a2", "a3"}, {"d7", "d5"},
				{"a3", "a4"}, {"d5", "d4"},
				{"e2", "e4"}, {"h7", "h6"},
			},
			"d4",
			[]string{"c3", "e3"},
			false,
		},
	}, t)
}
