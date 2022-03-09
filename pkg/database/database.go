package database

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Username string `bson:"uname"`
	Password string `bson:"pass"`
}

type Database struct {
	db  *mongo.Database
	ctx context.Context
}

func (d *Database) DBConnect() error {
	clientOptions := options.Client()
	clientOptions.ApplyURI(os.Getenv("DB_URI"))
	client, err := mongo.NewClient(clientOptions)

	if err != nil {
		return err
	}

	err = client.Connect(d.ctx)
	if err != nil {
		return err
	}

	d.db = client.Database(os.Getenv("DB_NAME"))

	return nil
}

func (d *Database) NewUser(data User) {
	_, err := d.db.Collection("user_db").InsertOne(d.ctx, data)
	if err != nil {
		log.Println("[MONGO]", err.Error())
	}
}

func (d *Database) FindUser(user string) User {
	var result User

	err := d.db.Collection("user_db").FindOne(d.ctx, bson.M{"uname": user}).Decode(&result)
	if err != nil {
		log.Println("[MONGO]", err.Error())
	}

	return result
}
