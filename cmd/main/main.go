package main

import (
	"context"
	"log"
	. "login/pkg/database"
	. "login/pkg/user/controller"
	. "login/pkg/user/repository"
	. "login/pkg/user/service"
	"os"

	"github.com/joho/godotenv"
)

var db = UserRepository{}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
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
	userService := NewUserService(userRepository)
	controller := NewController(userService)

	controller.Run()
}
