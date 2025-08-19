package jwt

import (
	"encoding/json"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(payload interface{}, secret []byte) (string, error) {
	data, err := structToClaims(payload)
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	return token.SignedString(secret)
}

func ParseToken(tokenStr string, secret []byte, out interface{}) error {
	parser := jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	token, err := parser.ParseWithClaims(tokenStr, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return fmt.Errorf("failed to decode token: %v", err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("failed to decode claims")
	}
	return claimsToStruct(claims, out)
}

// Helpers

func structToClaims(v interface{}) (jwt.MapClaims, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return jwt.MapClaims(m), nil
}

func claimsToStruct(claims jwt.MapClaims, out interface{}) error {
	b, err := json.Marshal(claims)
	if err != nil {
		return fmt.Errorf("marshal claims: %w", err)
	}
	if err := json.Unmarshal(b, out); err != nil {
		return fmt.Errorf("unmarshal into struct: %w", err)
	}
	return nil
}
