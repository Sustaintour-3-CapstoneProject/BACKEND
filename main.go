package main

import (
	"backend/config"
	"backend/routes"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// Initialize Database
	config.InitDB()

	// Register Routes
	routes.RegisterRoutes(e)

	// Start Server
	e.Logger.Fatal(e.Start(":8000"))
}
