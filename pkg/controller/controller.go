package controller

import (
	"github.com/gorilla/mux"
	AuthController "github.com/suryaadi44/LoginTest/internal/auth/controller"
	UserRepository "github.com/suryaadi44/LoginTest/internal/auth/repository"
	UserService "github.com/suryaadi44/LoginTest/internal/auth/service"
	Middleware "github.com/suryaadi44/LoginTest/internal/middleware"
	SessionRepository "github.com/suryaadi44/LoginTest/internal/session/repository"
	SessionService "github.com/suryaadi44/LoginTest/internal/session/service"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitializeController(router *mux.Router, db *mongo.Database) {
	userRepository := UserRepository.NewUserRepository(db)
	userService := UserService.NewUserService(userRepository)

	sessionRepository := SessionRepository.NewSessionRepository(db)
	sessionService := SessionService.NewSessionService(sessionRepository)

	middlewareService := Middleware.NewMiddlewareService(sessionService)
	controller := AuthController.NewController(router, userService, sessionService, middlewareService)
	controller.InitializeController()
}
