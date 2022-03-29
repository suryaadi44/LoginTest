package main

import (
	"context"
	"log"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/suryaadi44/LoginTest/pkg/controller"
	Database "github.com/suryaadi44/LoginTest/pkg/database"
	Server "github.com/suryaadi44/LoginTest/pkg/server"
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

	router := mux.NewRouter()

	controller.InitializeController(router, db)

	server := Server.NewServer(os.Getenv("PORT"), router)
	server.Run()
}
