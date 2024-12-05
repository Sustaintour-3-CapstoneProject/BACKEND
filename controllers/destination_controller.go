package controllers

import (
	"backend/config"
	"backend/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Input struct {
	Name             string       `json:"name"`
	City             string       `json:"city"`
	Position         float64      `json:"position"`
	Address          string       `json:"address"`
	OperationalHours string       `json:"operational_hours"`
	TicketPrice      float64      `json:"ticket_price"`
	Category         string       `json:"category"`
	Description      string       `json:"description"`
	Facilities       string       `json:"facilities"`
	Image            []string     `json:"image"`
	Video            []VideoInput `json:"video_contents"`
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
		Description:      jsonBody.Description,
	}

	// Simpan destinasi ke database
	if err := config.DB.Create(&destination).Error; err != nil {
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
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to add image"})
		}
	}

	// Muat ulang destinasi dengan properti City
	if err := config.DB.Preload("City").Preload("Images").Preload("VideoContents").First(&destination, destination.ID).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch destination with related data"})
	}

	// Kembalikan respons dengan properti City yang lengkap
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

	queryName := c.QueryParam("name")
	queryCityID := c.QueryParam("city")

	query := config.DB.
		Preload("City").
		Preload("Images").
		Preload("VideoContents")

	if queryName != "" {
		query = query.Where("name LIKE ?", "%"+queryName+"%")
	}

	if queryCityID != "" {
		query = query.Where("city_id = ?", queryCityID)
	}

	err := query.Find(&destinations).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch destinations"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":      "Destinations fetched successfully",
		"destinations": destinations,
	})
}

func GetDetailDestination(c echo.Context) error {
	var destination models.Destination

	id := c.Param("id")

	err := config.DB.
		Preload("City").
		Preload("Images").
		Preload("VideoContents").
		First(&destination, "id = ?", id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"message": "Destination not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch destination details"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":     "Destination details fetched successfully",
		"destination": destination,
	})
}

func GetMostViewedVideoContent(c echo.Context) error {
	var results []struct {
		ID          uint   `json:"id"`
		Title       string `json:"title"`
		URL         string `json:"url"`
		Description string `json:"description"`
		ViewCount   int64  `json:"view_count"`
	}

	err := config.DB.Table("video_contents").
		Select("video_contents.id, video_contents.title, video_contents.url, video_contents.description, COUNT(video_content_views.id) as view_count").
		Joins("LEFT JOIN video_content_views ON video_contents.id = video_content_views.video_content_id").
		Group("video_contents.id").
		Order("view_count DESC").
		Scan(&results).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch data"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Videos with view count fetched successfully",
		"data":    results,
	})
}
