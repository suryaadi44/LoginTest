package repository

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	. "login/pkg/user/entity"
)

type UserRepository struct {
	db  *mongo.Database
	ctx context.Context
}

func (u *UserRepository) NewUser(data User) {
	_, err := u.db.Collection("user_db").InsertOne(u.ctx, data)
	if err != nil {
		log.Println("[MONGO]", err.Error())
	}
}

func (u *UserRepository) FindUser(user string) User {
	var result User

	err := u.db.Collection("user_db").FindOne(u.ctx, bson.M{"_id": user}).Decode(&result)
	if err != nil {
		log.Println("[MONGO]", err.Error())
	}

	return result
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{db: db}
}
