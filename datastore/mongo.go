package datastore

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDatastore contains the information of a mongo database
type MongoDatastore struct {
	DB     *mongo.Database
	Client *mongo.Client
	Logger echo.Logger
}

// NewDatastore creates a new mongo datastore
func NewDatastore(uri string, log echo.Logger) (*MongoDatastore, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	db := client.Database("apis_pf_db")

	return &MongoDatastore{
		DB:     db,
		Client: client,
		Logger: log,
	}, nil
}
