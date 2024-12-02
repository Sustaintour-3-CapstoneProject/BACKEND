package controllers

import (
	"backend/config"
	"backend/models"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Input struct {
	Name             string       `json:"name"`
	City             string       `json:"string"`
	Address          string       `json:"address"`
	OperationalHours string       `json:"operational_hours"`
	TicketPrice      float64      `json:"ticket_price"`
	Category         string       `json:"category"`
	Facilities       string       `json:"facilities"`
	Image            []string     `json:"image"`
	Video            []VideoInput `json:"video"`
}

// video
type VideoInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url"`
}

func CreateDestination(c echo.Context) error {
	jsonBody := new(Input)
	err := json.NewDecoder(c.Request().Body).Decode(&jsonBody)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid JSON body"})
	}

	destination := new(models.Destination)
	destination.Name = jsonBody.Name
	destination.City = jsonBody.City
	destination.Address = jsonBody.Address
	destination.OperationalHours = jsonBody.OperationalHours
	destination.TicketPrice = jsonBody.TicketPrice
	destination.Category = jsonBody.Category
	destination.Facilities = jsonBody.Facilities
	// if err := c.Bind(destination); err != nil {
	// 	return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid input"})
	// }

	if err := config.DB.Create(destination).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create destination"})
	}

	for i := 0; i < len(jsonBody.Image); i++ {
		image := new(models.Image)
		image.DestinationID = destination.ID
		image.URL = jsonBody.Image[i]
		if err := config.DB.Create(image).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to add image"})
		}
	}

	//video
	for i := 0; i < len(jsonBody.Video); i++ {
		video := new(models.VideoContent)
		video.DestinationID = destination.ID
		video.URL = jsonBody.Video[i].Url
		video.Title = jsonBody.Video[i].Title
		video.Description = jsonBody.Video[i].Description
		if err := config.DB.Create(video).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to add video"})
		}
	}

	return c.JSON(http.StatusOK, destination)
}

func GetAllDestinations(c echo.Context) error {
	var destinations []models.Destination

	// Preload gambar dan video terkait dengan destinasi
	if err := config.DB.Preload("Images").Preload("VideoContents").Find(&destinations).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch destinations"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":      "Destinations fetched successfully",
		"destinations": destinations,
	})

}
