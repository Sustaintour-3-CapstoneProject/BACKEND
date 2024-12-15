package main

import (
	"backend/config"
	"backend/routes"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// Initialize Database
	config.InitDB()

	os.Mkdir("assets", 0777)

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

	e.Static("/assets", "./assets")

	// Register Routes
	routes.InitRoutes(e)

	// Start Server
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	e.Logger.Fatal(e.Start(":" + os.Getenv("APP_PORT")))
}
