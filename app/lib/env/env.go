package env

import (
	"fmt"
	"os"
)

func MustEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("Env variable %s not set!", key))
	}
	return value
}
