package database

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBConnect(ctx context.Context) (*mongo.Database, error) {
	clientOptions := options.Client()
	clientOptions.ApplyURI(os.Getenv("DB_URI"))
	client, err := mongo.NewClient(clientOptions)

	if err != nil {
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return client.Database(os.Getenv("DB_NAME")), nil
}
