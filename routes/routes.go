package routes

import (
	"backend/controllers"
	"backend/middlewares"

	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo) {
	// Rute khusus admin dengan middleware AdminOnly
	adminGroup := e.Group("/admin")
	adminGroup.Use(middlewares.AdminOnly)

	// Rute autentikasi
	e.POST("/register", controllers.RegisterHandler)
	e.POST("/login", controllers.LoginHandler)
	e.GET("/logout", controllers.LogoutHandler)

	// Route untuk destinasi
	e.POST("/destination", controllers.CreateDestination)
	e.GET("/destination", controllers.GetAllDestinations)
	e.PUT("/destination/:id", controllers.UpdateDestination) // Endpoint untuk mengubah destinasi
	e.DELETE("/destination/:id", controllers.DeleteDestination)

	// route untuk personalized recommendation
	e.GET("/destinations", controllers.FilterDestinations)

	// route untuk menambahkan kota
	e.POST("/city", controllers.CreateCity)
	e.GET("/city", controllers.GetCity)

}
