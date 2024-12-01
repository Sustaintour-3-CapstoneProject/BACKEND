package controllers

import (
	"backend/config"
	"backend/models"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type Input struct {
	Name       string `json:"name"`
	Address       string `json:"address"`
	OperationalHours string `json:"operational_hours"`
	TicketPrice float64 `json:"ticket_price"`
	Category        string `json:"category"`
	Facilities string `json:"facilities"`
	Image []string `json:"image"`

}




func CreateDestination(c echo.Context) error {
	jsonBody := new(Input)
     err := json.NewDecoder(c.Request().Body).Decode(&jsonBody)
     if err != nil {

         log.Error("empty json body")
         return nil
     }

	destination := new(models.Destination)
	destination.Name = jsonBody.Name
	destination.Address = jsonBody.Address
	destination.OperationalHours = jsonBody.OperationalHours
	destination.TicketPrice = jsonBody.TicketPrice
	destination.Category = jsonBody.Category
	destination.Facilities = jsonBody.Facilities
	
	
	if err := config.DB.Create(destination).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create destination"})
	}

	
	for i := 0; i <len(jsonBody.Image); i++ {
		image := new(models.Image)
		image.DestinationID = destination.ID
		image.URL = jsonBody.Image[i]
		if err := config.DB.Create(image).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to add image"})
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
		"message":    "Destinations fetched successfully",
		"destinations": destinations,
	})


}