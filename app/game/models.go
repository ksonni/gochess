package game

import (
	"github.com/google/uuid"
)

type StartGameRequest struct {
	DurationMillis  int64 `json:"duration_millis"`
	IncrementMillis int64 `json:"increment_millis"`
}

type StartGameResponse struct {
	Id uuid.UUID `json:"id"`
}
