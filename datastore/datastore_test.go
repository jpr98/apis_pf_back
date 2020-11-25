package datastore

import (
	"os"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestNewDatastore(t *testing.T) {
	uri := ""
	var logger echo.Logger

	_, err := NewDatastore(uri, logger)
	if err == nil {
		t.Error("Error should not be nil if datastore can't connect")
	}

	password := os.Getenv("MONGO_PASSWORD")
	uri = "mongodb+srv://pf-server:" + password + "@cluster0.7ihuj.mongodb.net/apis_pf_db?retryWrites=true&w=majority"
}

func TestNewStorageDatastore(t *testing.T) {
	bucket := ""
	var logger echo.Logger

	_, err := NewStorageDatastore(bucket, logger)
	if err == nil {
		t.Error("Error should not be nil if datastore can't connect")
	}
}
