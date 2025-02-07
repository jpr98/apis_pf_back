package app

import (
	"net/http"
	"os"

	"github.com/jpr98/apis_pf_back/datastore"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type server struct {
	router   *echo.Echo
	database *datastore.MongoDatastore
	storage  *datastore.StorageDatastore
	logger   echo.Logger
}

var appServer = server{}

// StartServer configures and intialices the web server on port 8080
func StartServer() {
	configServer()

	port := os.Getenv("PORT")
	if port == "" {
		appServer.router.Logger.Fatal("$PORT must be set")
	}

	configDatabase()
	setMiddlewares()
	setRoutes()
	appServer.router.Logger.Fatal(appServer.router.Start(":" + port))
}

func configServer() {
	appServer.router = echo.New()
	appServer.logger = appServer.router.Logger
}

func configDatabase() {
	var uri string
	password := os.Getenv("MONGO_PASSWORD")
	if password == "" {
		// Connecting to local machine mongo instance
		uri = "mongodb://localhost:27017"
	} else {
		// Connecting to Atlas mongo instance
		uri = "mongodb+srv://pf-server:" + password + "@cluster0.7ihuj.mongodb.net/apis_pf_db?retryWrites=true&w=majority"
	}

	database, err := datastore.NewDatastore(uri, appServer.logger)
	if err != nil {
		appServer.logger.Fatal(err)
	}
	appServer.database = database

	bucket := "apis-pf-bucket"
	storage, err := datastore.NewStorageDatastore(bucket, appServer.logger)
	if err != nil {
		appServer.logger.Fatal(err)
	}
	appServer.storage = storage
}

func setMiddlewares() {
	appServer.router.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}, latency=${latency_human}\n",
	}))
	appServer.router.Use(middleware.CORSWithConfig(
		middleware.CORSConfig{
			AllowOrigins:     []string{"http://localhost:3000"},
			AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodPatch},
			AllowCredentials: true,
		}))
}
