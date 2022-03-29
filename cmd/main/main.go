package main

import (
	"context"
	"log"
	"os"

	Controller "github.com/suryaadi44/LoginTest/pkg/auth/controller"
	UserRepository "github.com/suryaadi44/LoginTest/pkg/auth/repository"
	UserService "github.com/suryaadi44/LoginTest/pkg/auth/service"
	Database "github.com/suryaadi44/LoginTest/pkg/database"
	Middleware "github.com/suryaadi44/LoginTest/pkg/middleware"
	SessionRepository "github.com/suryaadi44/LoginTest/pkg/session/repository"
	SessionService "github.com/suryaadi44/LoginTest/pkg/session/service"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("[Env File] Error loading .env file")
	}

	appName, present := os.LookupEnv("APP_NAME")

	if !present {
		log.Fatal("[Env Var] Env variable not configure correctly")
	}

	log.Println("[Start]", appName)
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
