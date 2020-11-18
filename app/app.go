package app

import (
	"os"

	"github.com/jpr98/apis_pf_back/datastore"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type server struct {
	router   *echo.Echo
	database *datastore.MongoDatastore
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

	setMiddlewares()
	setRoutes()
	appServer.router.Logger.Fatal(appServer.router.Start(":" + port))
}

func configServer() {
	appServer.router = echo.New()
	appServer.logger = appServer.router.Logger

	password := os.Getenv("MONGO_PASSWORD")
	appServer.logger.Infof("Got PASSWORD from env: ", password)
	database, err := datastore.NewDatastore(password, appServer.logger)
	if err != nil {
		appServer.logger.Fatal(err)
	}
	appServer.database = database
}

func setMiddlewares() {
	appServer.router.Use(middleware.Logger())
}
