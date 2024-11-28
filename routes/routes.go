package routes

import (
	"github.com/labstack/echo/v4"

	"backend/controllers"
)

func RegisterRoutes(e *echo.Echo) {
	e.POST("/register", controllers.Register)
	e.POST("/login", controllers.Login)
	e.POST("/logout", controllers.Logout)

}
