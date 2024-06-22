package main

import (
	"auth-goker/internal/db"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env.dev")

	_ = db.MustInitDb()
}
