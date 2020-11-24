package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Contribution represents a contribution in the system
type Contribution struct {
	ID     primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	User   ContributionUser   `json:"user,omitempty" bson:"user,omitempty"`
	Amount float32            `json:"amount,omitempty" bson:"amount,omitempty"`
	Date   time.Time          `json:"date,omitempty" bson:"date,omitempty"`
}

// ContributionUser helps model a user inside a contribution
type ContributionUser struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Lastname string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
	Avatar   string             `json:"avatar,omitempty" bson:"avatar_url,omitempty"`
}
