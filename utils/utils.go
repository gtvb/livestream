package utils

import (
	"context"

	"github.com/gtvb/livestream/infra/db"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	UserCollectionTest       = "users_test"
	LiveStreamCollectionTest = "livestreams_test"
)

type TestContainer struct {
	databaseName string
	ctx          context.Context

	MongoContainer *mongodb.MongoDBContainer
	Database       *db.Database
}

// Cria um novo container que roda a imagem do Mongo
func NewTestContainer(databaseName string) (*TestContainer, error) {
	ctx := context.Background()

	mongoContainer, err := mongodb.Run(ctx, "mongo")
	if err != nil {
		return nil, err
	}

	return &TestContainer{ctx: ctx, databaseName: databaseName, MongoContainer: mongoContainer}, nil
}

// Esse método cria uma nova database dentro do container Mongo e
// instancia o nosso wrapper da Database, para que possamos
// utilizá-lo na criação dos repositórios
func (tc *TestContainer) SetupDatabaseWrapper() error {
	endpoint, err := tc.MongoContainer.ConnectionString(tc.ctx)
	if err != nil {
		return err
	}

	mongoClient, err := mongo.Connect(tc.ctx, options.Client().ApplyURI(endpoint))
	if err != nil {
		return err
	}

	db := &db.Database{mongoClient.Database(tc.databaseName)}
	tc.Database = db

	return nil
}

// Esse método finaliza o container
func (tc *TestContainer) Terminate() error {
	if err := tc.MongoContainer.Terminate(tc.ctx); err != nil {
		return err
	}

	return nil
}
