package routes

import (
	"backend/controllers"
	"backend/middlewares"

	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo) {
	e.GET("/dashboard/count-data", controllers.GetDashboardDataHandler)
	e.GET("/dashboard/graphic", controllers.GetDashboardGraphicDataHandler)

	userGroup := e.Group("/user", middlewares.AuthorizedAccess)
	userGroup.GET("", controllers.GetAllUserHandler)
	userGroup.POST("/category", controllers.CreateUserCategoryHandler)
	userGroup.GET("/:id", controllers.GetDetailUserHandler)
	userGroup.PUT("/change-password/:id", controllers.ChangePasswordHandler)
	userGroup.PUT("/:id", controllers.EditUserHandler)
	userGroup.DELETE("/:id", controllers.DeleteUser)

	e.POST("/register", controllers.RegisterHandler)
	e.POST("/login", controllers.LoginHandler)
	e.GET("/logout", controllers.LogoutHandler)

	destinationGroup := e.Group("/destination", middlewares.AuthorizedAccess)
	destinationGroup.GET("", controllers.GetAllDestinations)
	destinationGroup.GET("/personalized", controllers.GetPersonalizedDestinationByUser)
	destinationGroup.GET("/:id", controllers.GetDetailDestination)

	destinationVideoContentGroup := e.Group("/video-content")
	destinationVideoContentGroup.GET("", controllers.GetAllVideoContents)
	destinationVideoContentGroup.GET("/most", controllers.GetMostViewedVideoContent)

	e.POST("/city", controllers.CreateCity)
	e.GET("/city", controllers.GetCity)

	e.POST("/chat", controllers.ChatHandler)

	destinationGroup.POST("", controllers.CreateDestination, middlewares.AdminOnly)
	destinationGroup.POST("/assets", controllers.CreateDestinationAssetsHandler, middlewares.AdminOnly)
	destinationGroup.PUT("/assets", controllers.UpdateDestinationAssetsHandler)
	destinationGroup.PUT("/:id", controllers.UpdateDestination, middlewares.AdminOnly)
	destinationGroup.DELETE("/:id", controllers.DeleteDestination, middlewares.AdminOnly)

	routeGroup := e.Group("/route", middlewares.AuthorizedAccess)
	routeGroup.POST("", controllers.CreateRoute)
	routeGroup.GET("", controllers.GetRouteByUser)
	routeGroup.GET("/destination", controllers.GetDestinationsByRoute)
	routeGroup.DELETE("/:id", controllers.DeleteRoute)
}
