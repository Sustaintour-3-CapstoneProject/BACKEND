package main

import (
	"backend/config"
	"backend/models"
	"backend/routes"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// Initialize Database
	config.InitDB()

	// Run Migrations
	config.DB.AutoMigrate(&models.User{})

	// Register Routes
	routes.RegisterRoutes(e)

	// Start Server
	e.Logger.Fatal(e.Start(":8000"))
}
