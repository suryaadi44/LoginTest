package main

import (
	"context"
	"log"
	. "login/pkg/auth/controller"
	. "login/pkg/auth/repository"
	. "login/pkg/auth/service"
	. "login/pkg/database"
	. "login/pkg/middleware"
	. "login/pkg/session/repository"
	. "login/pkg/session/service"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("[Env File] Error loading .env file")
	}

	log.Println("[Start]", os.Getenv("APP_NAME"))
}

func main() {
	ctx := context.Background()

	db, err := DBConnect(ctx)
	if err != nil {
		log.Fatal("[Mongo]", err)
	}

	userRepository := NewUserRepository(db)
	sessionRepository := NewSessionRepository(db)
	userService := NewUserService(userRepository)
	sessionService := NewSessionService(sessionRepository)
	middlewareService := NewMiddlewareService(sessionService)

	controller := NewController(userService, sessionService, middlewareService)

	controller.Run()
}
