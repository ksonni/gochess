package auth

import "github.com/google/uuid"

type UserClaims struct {
	Id uuid.UUID `json:"id"`
}

func NewBasicClaims() UserClaims {
	return UserClaims{
		Id: uuid.New(),
	}
}
