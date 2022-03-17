package repository

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	. "login/pkg/session/entity"
)

type SessionRepository struct {
	db  *mongo.Database
	ctx context.Context
}

func (s *SessionRepository) NewSession(data Session) error {
	_, err := s.db.Collection("session").InsertOne(s.ctx, data)
	if err != nil {
		log.Println("[MONGO]", err.Error())
	}
	return err
}

func (s *SessionRepository) FindSession(uid string) (Session, error) {
	var result Session

	err := s.db.Collection("session").FindOne(s.ctx, bson.M{"_id": uid}).Decode(&result)
	if err != nil {
		log.Println("[MONGO]", err.Error())
	}

	return result, err
}

func (s *SessionRepository) DeleteSession(uid string) error {
	_, err := s.db.Collection("session").DeleteOne(s.ctx, bson.M{"_id": uid})
	if err != nil {
		log.Println("[MONGO]", err.Error())
	}

	return err
}

func NewSessionRepository(db *mongo.Database) *SessionRepository {
	return &SessionRepository{db: db}
}
