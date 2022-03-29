package main

import (
	"context"
	"log"
	Controller "login/pkg/auth/controller"
	UserRepository "login/pkg/auth/repository"
	UserService "login/pkg/auth/service"
	Database "login/pkg/database"
	Middleware "login/pkg/middleware"
	SessionRepository "login/pkg/session/repository"
	SessionService "login/pkg/session/service"
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

	db, err := Database.DBConnect(ctx)
	if err != nil {
		log.Fatal("[Mongo]", err)
	}

	userRepository := UserRepository.NewUserRepository(db)
	sessionRepository := SessionRepository.NewSessionRepository(db)
	userService := UserService.NewUserService(userRepository)
	sessionService := SessionService.NewSessionService(sessionRepository)
	middlewareService := Middleware.NewMiddlewareService(sessionService)

	controller := Controller.NewController(userService, sessionService, middlewareService)

	controller.Run()
}
