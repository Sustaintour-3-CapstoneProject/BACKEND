package main

import (
	"backend/config"
	"backend/routes"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// Initialize Database
	config.InitDB()

	// Apply CORS middleware with custom config
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
		},
	}))

	// Register Routes
	routes.InitRoutes(e)

	
	// Start Server
	e.Logger.Fatal(e.Start(":8000"))
}
