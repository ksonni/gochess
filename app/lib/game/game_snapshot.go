package game

import "time"

// Serialiable copy of game.Game
type GameSnapshot struct {
	Moves         []Move               `json:"moves"`
	Result        *ResultData          `json:"result,omitempty"`
	SnapshotTime  int64                `json:"snapshot_time"`
	RemainingTime map[PieceColor]int64 `json:"remaining_time"`
}

func (g *Game) Snapshot() GameSnapshot {
	result, _ := g.Result()
	remaining := make(map[PieceColor]int64)
	for color, clock := range g.clocks {
		remaining[color] = clock.RemainingTime().Milliseconds()
	}
	return GameSnapshot{
		Moves:         g.moves,
		Result:        result,
		SnapshotTime:  time.Now().UnixMilli(),
		RemainingTime: remaining,
	}
}
