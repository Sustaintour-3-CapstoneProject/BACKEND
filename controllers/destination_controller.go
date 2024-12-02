package controllers

import (
	"backend/config"
	"backend/models"
	"encoding/json"
	"fmt"
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
	// Decode JSON body
	jsonBody := new(Input)
	if err := json.NewDecoder(c.Request().Body).Decode(jsonBody); err != nil {
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
		CityID:           city.ID, // Gunakan ID dari City yang ditemukan
		Position:         jsonBody.Position,
		Address:          jsonBody.Address,
		OperationalHours: jsonBody.OperationalHours,
		TicketPrice:      jsonBody.TicketPrice,
		Category:         jsonBody.Category,
		Facilities:       jsonBody.Facilities,
	}

	// Simpan destinasi ke database
	if err := config.DB.Create(&destination).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create destination"})
	}

	// Simpan data gambar
	if len(jsonBody.Image) > 0 {
		for _, img := range jsonBody.Image {
			image := models.Image{
				DestinationID: destination.ID,
				URL:           img,
			}
			if err := config.DB.Create(&image).Error; err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to add image"})
			}
		}
	} else {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "No images provided"})
	}

	// Simpan data video
	if len(jsonBody.Video) > 0 {
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
	} else {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "No videos provided"})
	}

	// Muat ulang destinasi dengan properti City dan relasi lainnya (Images, VideoContents)
	if err := config.DB.Preload("City").Preload("Images").Preload("VideoContents").First(&destination, destination.ID).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch destination with related data"})
	}

	// Kembalikan respons dengan destinasi lengkap, termasuk gambar dan video
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

	// Cari CityID berdasarkan nama kota
	var city models.City
	if err := config.DB.Where("name = ?", jsonBody.City).First(&city).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "City not found"})
	}

	// Perbarui data destinasi
	destination.Name = jsonBody.Name
	destination.CityID = city.ID // Gunakan CityID yang benar
	destination.Address = jsonBody.Address
	destination.OperationalHours = jsonBody.OperationalHours
	destination.TicketPrice = jsonBody.TicketPrice
	destination.Category = jsonBody.Category
	destination.Facilities = jsonBody.Facilities

	// Simpan perubahan ke database
	if err := config.DB.Save(&destination).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Gagal memperbarui destinasi"})
	}

	// Kembalikan respons berhasil
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

	// Preload semua data terkait
	err := config.DB.
		Preload("City").
		Preload("Images").
		Preload("VideoContents").
		Find(&destinations).Error

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch destinations"})
	}

	// Debug log untuk melihat hasil preload
	for _, dest := range destinations {
		fmt.Printf("Destination ID: %d, VideoContents: %v\n", dest.ID, dest.VideoContents)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":      "Destinations fetched successfully",
		"destinations": destinations,
	})
}
