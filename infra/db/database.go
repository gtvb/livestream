package db

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	*mongo.Database
}

// Inicia uma nova conexão com o banco Mongo, já
// selecionando a Database contida na variável de
// ambiente MONGODB_DATABASE_NAME.
func NewDb() (*Database, error) {
	mongoUri := os.Getenv("MONGODB_CONNECTION_STRING")
	mongoDatabaseName := os.Getenv("MONGODB_DATABASE_NAME")

	client, err := mongo.Connect(context.TODO(),
		options.Client().
			ApplyURI(mongoUri).SetAuth(options.Credential{
			Username: os.Getenv("MONGODB_USERNAME"),
			Password: os.Getenv("MONGODB_PASSWORD"),
		}),
	)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = client.Ping(context.TODO(), options.Client().ReadPreference)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	database := client.Database(mongoDatabaseName)

	return &Database{database}, nil
}
