package env

import (
	"github.com/kataras/golog"
	"os"
)

func EnsureEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		golog.Fatal("Configuration ", key, " is not found, or is empty, despite being required.")
	}
	return value
}
