package config

import (
	"backend/models"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	// Load the appropriate .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connectDatabase()
}

func TestInitDB() {
	// Load the test-specific .env file
	err := godotenv.Load(".env.test")
	if err != nil {
		log.Fatal("Error loading .env.test file")
	}

	connectDatabase()
}

func connectDatabase() {
	// Get database credentials from environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Build DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	// Connect to the database
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// AutoMigrate models
	DB.AutoMigrate(
		&models.User{},
		&models.Destination{},
		&models.VideoContent{},
		&models.Image{},
		&models.City{},
		&models.VideoContentView{},
		&models.Route{},
		&models.RouteDestination{},
	)
}
