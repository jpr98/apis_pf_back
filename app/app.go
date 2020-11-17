package app

import (
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
	appServer.router.Logger.Fatal(appServer.router.Start(":8080"))
}

func configServer() {
	appServer.router = echo.New()
	appServer.logger = appServer.router.Logger
	database, err := datastore.NewDatastore(appServer.logger)
	if err != nil {
		appServer.logger.Fatal(err)
	}
	appServer.database = database
}

func setMiddlewares() {
	appServer.router.Use(middleware.Logger())
}
