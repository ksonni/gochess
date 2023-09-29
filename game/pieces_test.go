package game

import (
	"testing"
)

func TestPiecesAttackedSquares(t *testing.T) {
	runAttackedSquaresTests(map[string]targetedSquaresTestInput{
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

func TestPlanPossibleMoves(t *testing.T) {
	runPossibleMovesTests(map[string]targetedSquaresTestInput{
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
			[]string{"d5", "e5", "f5"},
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

type targetedSquaresTestInput struct {
	moves      []testMove
	testSquare string
	targeted   []string
	jump       bool
}

type targetedSquaresTestFn = func(title string, g *Game,
	testSquare string, want []string, t *testing.T)

func runTargetedSquaresTests(tests map[string]targetedSquaresTestInput,
	testFn targetedSquaresTestFn, t *testing.T) {
	for title, test := range tests {
		g := playGame(title, test.moves, test.jump, t)
		testFn(title, g, test.testSquare, test.targeted, t)
	}
}

func runAttackedSquaresTests(tests map[string]targetedSquaresTestInput, t *testing.T) {
	runTargetedSquaresTests(tests, testAttackedSquares, t)
}

func testAttackedSquares(title string, g *Game, testSquare string, want []string, t *testing.T) {
	attackedResult := g.ComputeAttackedSquares(sq(testSquare))
	assertSquareMapEquals(title, attackedResult, want, t)
}

func runPossibleMovesTests(tests map[string]targetedSquaresTestInput, t *testing.T) {
	runTargetedSquaresTests(tests, testPossibleMoves, t)
}

func testPossibleMoves(title string, g *Game, testSquare string, want []string, t *testing.T) {
	possibleMoves := g.PlanPossibleMoves(sq(testSquare))
	sqMap := make(map[Square]bool)
	for _, move := range possibleMoves {
		sqMap[move.To] = true
	}
	assertSquareMapEquals(title, sqMap, want, t)
}
