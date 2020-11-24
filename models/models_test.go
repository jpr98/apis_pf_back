package models

import (
	"os"
	"testing"
	"time"

	"github.com/jpr98/apis_pf_back/datastore"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetCommentsAuthors(t *testing.T) {
	var logger echo.Logger
	password := os.Getenv("MONGO_PASSWORD")
	uri := "mongodb+srv://pf-server:" + password + "@cluster0.7ihuj.mongodb.net/apis_pf_db?retryWrites=true&w=majority"

	datastore, _ := datastore.NewDatastore(uri, logger)
	ps := NewProjectStore(datastore.DB)
	project := Project{}

	ps.getCommentsAuthors(&project)
	if len(project.Comments) != 0 {
		t.Error("Comments length should be 0")
	}
}

func TestGetContirbutionsUsers(t *testing.T) {
	var logger echo.Logger
	password := os.Getenv("MONGO_PASSWORD")
	uri := "mongodb+srv://pf-server:" + password + "@cluster0.7ihuj.mongodb.net/apis_pf_db?retryWrites=true&w=majority"

	datastore, _ := datastore.NewDatastore(uri, logger)
	ps := NewProjectStore(datastore.DB)
	project := Project{}

	ps.getContributionsUsers(&project)
	if len(project.Comments) != 0 {
		t.Error("Contributions length should be 0")
	}
}

func TestGeneratePassword(t *testing.T) {
	password := "test"
	hashedPassword, err := generatePassword(password)
	if err != nil {
		t.Error("Password generation failed")
	}

	if len(hashedPassword) <= len(password) {
		t.Error("Password not hashed")
	}

	if hashedPassword == password {
		t.Error("Password not hashed")
	}
}

func TestComment(t *testing.T) {
	author := CommentAuthor{primitive.NewObjectID(), "Name", "avatar.com/url"}
	comment := Comment{primitive.NewObjectID(), author, time.Now(), "Content"}

	if comment.Author != author {
		t.Error("Author should be set in comment")
	}

	if comment.Text != "Content" {
		t.Error("Text should be set in comment")
	}
}

func TestContribute(t *testing.T) {
	author := CommentAuthor{primitive.NewObjectID(), "Name", "avatar.com/url"}
	comment := Comment{primitive.NewObjectID(), author, time.Now(), "Content"}

	if comment.Author != author {
		t.Error("Author should be set in comment")
	}

	if comment.Text != "Content" {
		t.Error("Text should be set in comment")
	}
}
