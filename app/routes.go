package app

import (
	"github.com/jpr98/apis_pf_back/controllers"
	"github.com/jpr98/apis_pf_back/models"
	"github.com/labstack/echo/v4/middleware"
)

func setRoutes() {
	setUserRoutes()
	setProjectRoutes()
}

func setUserRoutes() {
	userStore := models.NewUserStore(appServer.database.DB)
	usersController := controllers.NewUsersController(*userStore)

	appServer.router.POST("/signup", usersController.Create)
	appServer.router.GET("/login", usersController.Login)

	u := appServer.router.Group("/users")
	u.Use(middleware.JWT([]byte("secret")))
	u.GET("/:id", usersController.GetByID)
	u.PATCH("/:id", usersController.Update)
}

func setProjectRoutes() {
	projectStore := models.NewProjectStore(appServer.database.DB)
	projectsController := controllers.NewProjectsController(*projectStore)

	appServer.router.GET("projects/:id", projectsController.GetByID)
	appServer.router.GET("/projects/search", projectsController.SearchProject)
	appServer.router.GET("/projects/owned", projectsController.GetByOwner)

	p := appServer.router.Group("/projects")
	p.Use(middleware.JWT([]byte("secret")))
	p.POST("/new", projectsController.Create)
	p.POST("/:id/vote", projectsController.VoteForProject)
}
