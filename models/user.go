package models

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// User model represents a user on the system
type User struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name,omitempty" bson:"name,omitempty"`
	Lastname  string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
	Email     string             `json:"email,omitempty" bson:"email,omitempty"`
	Password  string             `json:"password,omitempty" bson:"password,omitempty"`
	Status    string             `json:"status,omitempty" bson:"status,omitempty"`
	Avatar    string             `json:"avatar_url,omitempty" bson:"avatar,omitempty"`
	Bio       string             `json:"bio,omitempty" bson:"bio,omitempty"`
	Location  string             `json:"location,omitempty" bson:"location,omitempty"`
	Birthdate string             `json:"birthdate,omitempty" bson:"birthdate,omitempty"`

	Projects            []Project `json:"projects,omitempty" bson:"-"`
	ContributedProjects []Project `json:"contributed_projects,omitempty" bson:"-"`
	VotedProjects       []Project `json:"voted_projects,omitempty" bson:"-"`
}

// UserStore contains all the CRUD operations for the User model
type UserStore struct {
	collection *mongo.Collection
}

// NewUserStore creates a user store with a mongo database
func NewUserStore(database *mongo.Database) *UserStore {
	return &UserStore{database.Collection("users")}
}

// Create stores a new user in the users collection
func (us *UserStore) Create(u User) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	u.Password, err = generatePassword(u.Password)
	if err != nil {
		return User{}, err
	}

	u.Status = "active"

	result, err := us.collection.InsertOne(ctx, u)
	if err != nil {
		return User{}, err
	}

	generatedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return User{}, errors.New("Invalid generated id on user")
	}
	u.ID = generatedID

	return u, nil
}

// ValidEmail checks if an email is already taken
func (us *UserStore) ValidEmail(email string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user User
	err := us.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	return err != nil
}

// GetByID gets a user with a given id from the database
func (us *UserStore) GetByID(id string) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user User
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return User{}, err
	}

	err = us.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// GetByEmail retrieves a user by a given email
func (us *UserStore) GetByEmail(email string) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user User
	err := us.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func generatePassword(plainTextPassword string) (string, error) {
	bytePassword, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytePassword), nil
}
