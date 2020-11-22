package models

import (
	"context"
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Project represents a project in the system
type Project struct {
	ID            primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Owner         primitive.ObjectID   `json:"owner,omitempty" bson:"owner,omitempty"`
	Title         string               `json:"title,omitempty" bson:"title,omitempty"`
	Description   string               `json:"description,omitempty" bson:"desc,omitempty"`
	Tags          []string             `json:"tags,omitempty" bson:"tags,omitempty"`
	Category      string               `json:"category,omitempty" bson:"category,omitempty"`
	Location      string               `json:"location,omitempty" bson:"location,omitempty"`
	Votes         []primitive.ObjectID `json:"votes,omitempty" bson:"votes,omitempty"`
	Subscriptions []primitive.ObjectID `json:"subscriptions,omitempty" bson:"subscriptions,omitempty"`
	Multimedia    []string             `json:"multimedia,omitempty" bson:"multimedia,omitempty"`
	// Comments []Comment
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
func (ps *ProjectStore) Create(p Project, ownerID string) (Project, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(ownerID)
	if err != nil {
		return Project{}, err
	}

	for index, tag := range p.Tags {
		p.Tags[index] = strings.ToLower(tag)
	}

	p.Owner = oid
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

// GetByTitle returns all projects with titles containing the given query string
func (ps *ProjectStore) GetByTitle(title string) ([]Project, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{"title", primitive.Regex{Pattern: ".*" + title + ".*", Options: ""}}}
	cursor, err := ps.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	projects := make([]Project, 0)
	for cursor.Next(ctx) {
		var project Project
		err = cursor.Decode(&project)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	cursor.Close(ctx)
	return projects, nil
}

// GetByTags returns all projects for a given set of tags
func (ps *ProjectStore) GetByTags(tags []string) ([]Project, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := ps.collection.Find(ctx, bson.D{{"tags", bson.D{{"$in", tags}}}})
	if err != nil {
		return nil, err
	}

	projects := make([]Project, 0)
	for cursor.Next(ctx) {
		var project Project
		err = cursor.Decode(&project)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	cursor.Close(ctx)
	return projects, nil
}

// GetByCategory returns all projects for a given category
func (ps *ProjectStore) GetByCategory(category string) ([]Project, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := ps.collection.Find(ctx, bson.M{"category": category})
	if err != nil {
		return nil, err
	}

	projects := make([]Project, 0)
	for cursor.Next(ctx) {
		var project Project
		err = cursor.Decode(&project)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	cursor.Close(ctx)
	return projects, nil
}

// GetByOwnerID returns projects with a given owner ID
func (ps *ProjectStore) GetByOwnerID(ownerID string) ([]Project, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(ownerID)
	if err != nil {
		return nil, err
	}

	cursor, err := ps.collection.Find(ctx, bson.M{"owner": oid})
	if err != nil {
		return nil, err
	}

	projects := make([]Project, 0)
	for cursor.Next(ctx) {
		var project Project
		err = cursor.Decode(&project)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	cursor.Close(ctx)
	return projects, nil
}

// Vote appends or removes a user to the list of votes of a project
func (ps *ProjectStore) Vote(projectID string, userID string, upvote bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pid, err := primitive.ObjectIDFromHex(projectID)
	if err != nil {
		return err
	}

	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	var updateAction string
	if upvote {
		updateAction = "$addToSet"
	} else {
		updateAction = "$pull"
	}
	update := bson.M{updateAction: bson.M{"votes": uid}}
	result, err := ps.collection.UpdateOne(ctx, bson.M{"_id": pid}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("No project found with given id")
	}

	return nil
}
