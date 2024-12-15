package main

import (
	"backend/config"
	_ "backend/docs"
	"backend/routes"
	"log"
	"os"

	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// @title           TRIPWISE API
// @version         1.0
// @description     API Documentation TripWise Capstone Project - Team 3
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      www.tripwise.my.id

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

func main() {
	e := echo.New()

	// Initialize Database
	config.InitDB()

	os.Mkdir("assets", 0777)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

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
