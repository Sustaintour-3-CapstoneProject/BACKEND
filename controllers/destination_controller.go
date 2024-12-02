package controllers

import (
	"backend/config"
	"backend/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Input struct {
	Name             string       `json:"name"`
	City             string       `json:"city"`
	Position         float64      `json:"position"`
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
	if err := json.NewDecoder(c.Request().Body).Decode(&jsonBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid JSON body"})
	}

	// Cari City berdasarkan nama
	var city models.City
	if err := config.DB.Where("name = ?", jsonBody.City).First(&city).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "City not found"})
	}

	// Buat destinasi baru
	destination := models.Destination{
		Name:             jsonBody.Name,
		CityID:           city.ID,
		Position:         jsonBody.Position,
		Address:          jsonBody.Address,
		OperationalHours: jsonBody.OperationalHours,
		TicketPrice:      jsonBody.TicketPrice,
		Category:         jsonBody.Category,
		Facilities:       jsonBody.Facilities,
	}

	if err := config.DB.Create(&destination).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create destination"})
	}

	// Simpan data gambar
	for _, img := range jsonBody.Image {
		image := models.Image{
			DestinationID: destination.ID,
			URL:           img,
		}
		if err := config.DB.Create(&image).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to add image"})
		}
	}

	// Simpan data video
	for _, vid := range jsonBody.Video {
		video := models.VideoContent{
			DestinationID: destination.ID,
			Title:         vid.Title,
			Description:   vid.Description,
			URL:           vid.Url,
		}
		if err := config.DB.Create(&video).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to add video"})
		}
	}

	return c.JSON(http.StatusOK, destination)
}

func UpdateDestination(c echo.Context) error {
	// Ambil ID destinasi dari parameter URL
	id := c.Param("id")
	destinationID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "ID destinasi tidak valid"})
	}

	// Parse body request ke struct Input
	jsonBody := new(Input)
	if err := json.NewDecoder(c.Request().Body).Decode(jsonBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid JSON body"})
	}

	// Cari destinasi berdasarkan ID
	var destination models.Destination
	if err := config.DB.First(&destination, destinationID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Destinasi tidak ditemukan"})
	}

	// Perbarui data destinasi
	destination.Name = jsonBody.Name
	destination.City = jsonBody.City
	destination.Address = jsonBody.Address
	destination.OperationalHours = jsonBody.OperationalHours
	destination.TicketPrice = jsonBody.TicketPrice
	destination.Category = jsonBody.Category
	destination.Facilities = jsonBody.Facilities

	// Simpan perubahan ke database
	if err := config.DB.Save(&destination).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Gagal memperbarui destinasi"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Destinasi berhasil diperbarui"})
}

// Fungsi untuk menghapus destinasi berdasarkan ID
func DeleteDestination(c echo.Context) error {
	// Ambil ID destinasi dari parameter URL
	id := c.Param("id")
	destinationID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "ID destinasi tidak valid"})
	}

	// Cari destinasi berdasarkan ID
	var destination models.Destination
	if err := config.DB.First(&destination, destinationID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Destinasi tidak ditemukan"})
	}

	// Hapus destinasi
	if err := config.DB.Delete(&destination).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Gagal menghapus destinasi"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Destinasi berhasil dihapus"})
}

func GetAllDestinations(c echo.Context) error {
	var destinations []models.Destination

	// Preload semua data terkait: City, Images, dan VideoContents
	if err := config.DB.
		Preload("City").          // Preload relasi City
		Preload("Images").        // Preload relasi Images
		Preload("VideoContents"). // Preload relasi VideoContents
		Find(&destinations).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch destinations"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":      "Destinations fetched successfully",
		"destinations": destinations,
	})
}
