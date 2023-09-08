package game

import (
	"fmt"
	"strings"
)

type Square struct {
	File int
	Rank int
}

func (square Square) String() string {
	return fmt.Sprintf("%c%d", 'a'+square.File, square.Rank+1)
}

func MustSquare(notation string) Square {
	square, err := ParseSquare(notation)
	if err != nil {
		panic(fmt.Sprintf("game: failed to parse invalid square %s", notation))
	}
	return *square
}

func ParseSquare(notation string) (*Square, error) {
	var file rune
	var rank int
	_, err := fmt.Sscanf(strings.ToLower(notation), "%c%d", &file, &rank)
	if err != nil {
		return nil, err
	}
	return &Square{File: int(file - 'a'), Rank: rank - 1}, nil
}

func (square Square) Adding(delta Square) Square {
	return Square{
		File: square.File + delta.File,
		Rank: square.Rank + delta.Rank,
	}
}

func (square Square) Subtracting(delta Square) Square {
	return square.Adding(delta.Multiplying(Square{File: -1, Rank: -1}))
}

func (square Square) Multiplying(delta Square) Square {
	return Square{
		File: square.File * delta.File,
		Rank: square.Rank * delta.Rank,
	}
}
