package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Comment represents a comment in a project
type Comment struct {
	ID     primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Author CommentAuthor      `json:"author,omitempty" bson:"author,omitempty"`
	Date   time.Time          `json:"date,omitempty" bson:"date,omitempty"`
	Text   string             `json:"text,omitempty" bson:"text,omitempty"`
}

// CommentAuthor helps model a user inside a comment
type CommentAuthor struct {
	ID     primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name   string             `json:"name,omitempty" bson:"name,omitempty"`
	Avatar string             `json:"avatar,omitempty" bson:"avatar_url,omitempty"`
}
