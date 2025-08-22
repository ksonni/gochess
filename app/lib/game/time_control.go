package game

import (
	"time"
)

type TimeControl struct {
	Total     time.Duration
	Increment time.Duration
}

func (t TimeControl) Validate() bool {
	return t.Total >= 1*time.Minute && t.Total <= 24*time.Hour &&
		t.Increment >= 0 && t.Increment <= 2*time.Minute
}

var (
	TimeControl_TwoOne = TimeControl{
		Total:     2 * time.Minute,
		Increment: time.Second,
	}
	TimeControl_ThreeTwo = TimeControl{
		Total:     3 * time.Minute,
		Increment: 2 * time.Second,
	}
	TimeControl_Five = TimeControl{
		Total: 5 * time.Minute,
	}
	TimeControl_FiveTwo = TimeControl{
		Total:     5 * time.Minute,
		Increment: 2 * time.Second,
	}
	TimeControl_FiveFive = TimeControl{
		Total:     5 * time.Minute,
		Increment: 5 * time.Second,
	}
	TimeControl_FifteenTen = TimeControl{
		Total:     15 * time.Minute,
		Increment: 10 * time.Second,
	}
	TimeControl_Thirty = TimeControl{
		Total: 30 * time.Minute,
	}
	TimeControl_Hour = TimeControl{
		Total: time.Hour,
	}
)
