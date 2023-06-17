package env

import (
	"github.com/joho/godotenv"
	"github.com/kataras/golog"
)

func Init() {
	golog.Info("Loading .env file")
	if err := godotenv.Load(); err != nil {
		golog.Fatal("No .env file can be found.")
	}
}
