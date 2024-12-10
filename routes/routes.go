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
	destinationGroup := e.Group("/destination")

	destinationGroup.POST("", controllers.CreateDestination)
	destinationGroup.GET("", controllers.GetAllDestinations)
	destinationGroup.GET("/personalized", controllers.GetPersonalizedDestinationByUser)
	destinationGroup.GET("/:id", controllers.GetDetailDestination)
	destinationGroup.DELETE("/:id", controllers.DeleteDestination)
	destinationGroup.PUT("/:id", controllers.UpdateDestination)

	destinationVideoContentGroup := e.Group("/video-content")
	destinationVideoContentGroup.GET("/most", controllers.GetMostViewedVideoContent)

	// route untuk personalized recommendation
	e.GET("/destinations", controllers.FilterDestinations)

	// route untuk menambahkan kota
	e.POST("/city", controllers.CreateCity)
	e.GET("/city", controllers.GetCity)

	routeGroup := e.Group("/route", middlewares.AuthorizedAccess)
	routeGroup.POST("", controllers.CreateRoute, middlewares.AdminOnly)
	routeGroup.GET("", controllers.GetRouteByUser, middlewares.AdminOnly)
	routeGroup.DELETE("/:id", controllers.DeleteRoute, middlewares.AdminOnly)

}
