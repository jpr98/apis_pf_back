package app

import (
	"github.com/jpr98/apis_pf_back/controllers"
	"github.com/jpr98/apis_pf_back/models"
	"github.com/labstack/echo/v4/middleware"
)

func setRoutes() {
	setUserRoutes()
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
