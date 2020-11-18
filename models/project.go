package models

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Project represents a project in the system
type Project struct {
	ID          primitive.ObjectID
	Owner       primitive.ObjectID
	Title       string
	Description string
	Tags        []string
	Category    string
	Location    string
	Votes       []primitive.ObjectID
	// Comments []Comment
	Subscriptions []primitive.ObjectID
	Multimedia    []string
}

// ProjectStore contains all the CRUD operations of Project
type ProjectStore struct {
	collection *mongo.Collection
}

// NewProjectStore creates a project store with a mongo database
func NewProjectStore(database *mongo.Database) *ProjectStore {
	return &ProjectStore{database.Collection("projects")}
}

// Create receives a project object and tries to insert it to the project store
func (ps *ProjectStore) Create(p Project) (Project, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := ps.collection.InsertOne(ctx, p)
	if err != nil {
		return Project{}, err
	}

	generatedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return Project{}, errors.New("Invalid generated id on project")
	}
	p.ID = generatedID

	return p, nil
}

// GetByID finds a project with a given id
func (ps *ProjectStore) GetByID(id string) (Project, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var project Project
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Project{}, err
	}

	err = ps.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&project)
	if err != nil {
		return Project{}, err
	}

	return project, nil
}

// func (ps *ProjectStore) GetByTitle(title string) ([]Project, error) {

// }

// func (ps *ProjectStore) GetByTags(tags []string) ([]Project, error) {

// }

// func (ps *ProjectStore) GetByCategories(categories []string) ([]Project, error) {

// }
