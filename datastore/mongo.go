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
func NewDatastore(password string, log echo.Logger) (*MongoDatastore, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://pf-server:" + password + "@cluster0.7ihuj.mongodb.net/apis_pf_db?retryWrites=true&w=majority"))
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
