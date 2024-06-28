package main

import (
	"auth-goker/internal/db"
	"auth-goker/internal/http"
	"auth-goker/internal/repository"
	"auth-goker/internal/service"
	"github.com/joho/godotenv"
	"log"
	"sync"
)

func main() {
	_ = godotenv.Load(".env.dev")

	dbInstance := db.MustInitDb()
	userRepo := repository.UserRepo{DB: dbInstance}
	sessionRepo := repository.SessionRepo{DB: dbInstance}
	authService := service.AuthService{Ur: userRepo, Sr: sessionRepo}
	authController := http.AuthController{As: authService}

	go http.RunServer(authController)
	log.Printf("Server successfully running")

	var wg sync.WaitGroup
	wg.Add(1)

	wg.Wait()
}
