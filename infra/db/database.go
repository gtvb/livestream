package db

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	*mongo.Database
}

func NewDb() (*mongo.Database, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	mongoUri := os.Getenv("MONGODB_CONNECTION_STRING")
	mongoDatabaseName := os.Getenv("MONGODB_DATABASE_NAME")

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoUri))
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(context.TODO())

	database := client.Database(mongoDatabaseName)

	return database, nil
}
