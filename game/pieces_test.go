package game

import (
	"testing"
)

func TestPiecesAttackedSquares(t *testing.T) {
	runAttackedSquaresTests(map[string]testAttackedSquaresInput{
		"Rook": {
			[]testMove{{"a1", "e4"}},
			"e4",
			[]string{
				"a4", "b4", "c4", "d4", "f4", "g4",
				"h4", "e3", "e5", "e6", "e7",
			},
			true,
		},
		"Dark square bishop": {
			[]testMove{{"c1", "f4"}},
			"f4",
			[]string{"e3", "g3", "e5", "g5", "d6", "h6", "c7"},
			true,
		},
		"Light square bishop": {
			[]testMove{{"f1", "c4"}},
			"c4",
			[]string{"a6", "b3", "b5", "d3", "d5", "e6", "f7"},
			true,
		},
		"King": {
			[]testMove{{"e1", "e6"}},
			"e6",
			[]string{"d5", "d6", "d7", "e5", "e7", "f5", "f6", "f7"},
			true,
		},
		"Queen": {
			[]testMove{{"d1", "e4"}},
			"e4",
			[]string{
				"a4", "b4", "b7", "c4", "c6", "d3", "d4", "d5", "e3", "e5",
				"e6", "e7", "f3", "f4", "f5", "g4", "g6", "h4", "h7",
			},
			true,
		},
		"Knight towards left": {
			[]testMove{},
			"b1",
			[]string{"a3", "c3"},
			true,
		},
		"Knight center": {
			[]testMove{{"b1", "e5"}},
			"e5",
			[]string{"c4", "c6", "d3", "d7", "f3", "f7", "g4", "g6"},
			true,
		},
		"Knight top rim": {
			[]testMove{{"b1", "c8"}},
			"c8",
			[]string{"a7", "b6", "d6", "e7"},
			true,
		},
	}, t)
}

// Helpers

type testAttackedSquaresInput struct {
	moves      []testMove
	testSquare string
	attacked   []string
	jump       bool
}

func runAttackedSquaresTests(tests map[string]testAttackedSquaresInput, t *testing.T) {
	for title, test := range tests {
		g := playGame(title, test.moves, test.jump, t)
		testAttackedSquares(title, g, test.testSquare, test.attacked, t)
	}
}

func testAttackedSquares(title string, g *Game, testSquare string, want []string, t *testing.T) {
	attackedResult := g.ComputeAttackedSquares(sq(testSquare))
	assertAttackedSquaresEqual(title, attackedResult, want, t)
}
