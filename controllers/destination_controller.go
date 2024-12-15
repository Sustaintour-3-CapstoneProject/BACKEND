package controllers

import (
	"backend/config"
	"backend/helper"
	"backend/models"
	"backend/request"
	"backend/response"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type CreateDestinationAssetsInput struct {
	DestinationID int                  `json:"destinationID"`
	Images        []string             `json:"images"`
	VideoContents []request.VideoInput `json:"video_contents"`
}

func CreateDestination(c echo.Context) error {
	// Decode JSON body
	jsonBody := new(request.CreateDestinationInput)
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
	jsonBody := new(request.CreateDestinationInput)
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

	config.DB.Where("destination_id = ?", destination.ID).Delete(&destination.Images)
	config.DB.Where("destination_id = ?", destination.ID).Delete(&destination.VideoContents)

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

	// Kembalikan respons berhasil
	return c.JSON(http.StatusOK, map[string]string{"message": "Destinasi berhasil diperbarui"})
}

// Fungsi untuk menghapus destinasi berdasarkan ID
func DeleteDestination(c echo.Context) error {
	// Retrieve the destination ID from the URL parameter
	id := c.Param("id")
	destinationID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid destination ID"})
	}

	// Find the destination by ID, including its related entities
	var destination models.Destination
	if err := config.DB.Preload("Images").Preload("VideoContents").First(&destination, destinationID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Destination not found"})
	}

	// Start a transaction to ensure atomicity
	tx := config.DB.Begin()

	// Delete related Images
	if err := tx.Where("destination_id = ?", destination.ID).Delete(&destination.Images).Error; err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to delete related images"})
	}

	// Delete related VideoContents
	if err := tx.Where("destination_id = ?", destination.ID).Delete(&destination.VideoContents).Error; err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to delete related video contents"})
	}

	// Delete the destination
	if err := tx.Delete(&destination).Error; err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to delete destination"})
	}

	// Commit the transaction
	tx.Commit()

	return c.JSON(http.StatusOK, map[string]string{"message": "Destination and related data successfully deleted"})
}

func GetAllDestinations(c echo.Context) error {
	var destinations []models.Destination
	var destinationResponses []response.DestinationResponse

	queryName := c.QueryParam("name")
	queryCityName := c.QueryParam("city")
	querySort := c.QueryParam("sort")
	queryCategory := c.QueryParam("category")

	query := config.DB.
		Preload("City").
		Preload("Images").
		Preload("VideoContents")

	if queryName != "" {
		query = query.Where("name LIKE ?", "%"+queryName+"%")
	}

	if queryCityName != "" {
		var city models.City
		if err := config.DB.Where("name = ?", queryCityName).First(&city).Error; err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Origin City not found"})
		}
		query = query.Where("city_id = ?", city.ID)
	}

	if queryCategory != "" {
		query = query.Where("category = ?", queryCategory)
	}

	if querySort != "" {
		if querySort == "oldest" {
			query = query.Order("created_at ASC")
		} else {
			query = query.Order("created_at DESC")
		}
	}

	err := query.Find(&destinations).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch destinations"})
	}

	for _, dest := range destinations {
		var facilitiesArray []string
		if dest.Facilities != "" {
			facilitiesArray = make([]string, 0)
			for _, facility := range strings.Split(dest.Facilities, ",") {
				facilitiesArray = append(facilitiesArray, strings.TrimSpace(facility))
			}
		}

		destinationResponse := response.DestinationResponse{
			ID:               dest.ID,
			Name:             dest.Name,
			City:             response.City{ID: dest.City.ID, Name: dest.City.Name},
			Position:         dest.Position,
			Address:          dest.Address,
			OperationalHours: dest.OperationalHours,
			TicketPrice:      dest.TicketPrice,
			Category:         dest.Category,
			Description:      dest.Description,
			Facilities:       facilitiesArray,
			CreatedAt:        dest.CreatedAt,
			Images:           convertImagesToResponse(dest.Images),
			VideoContents:    convertVideosToResponse(dest.VideoContents),
		}

		destinationResponses = append(destinationResponses, destinationResponse)
	}

	// Return the response with the destinations
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":      "Destinations fetched successfully",
		"destinations": destinationResponses,
	})
}

func GetDetailDestination(c echo.Context) error {
	var destination models.Destination
	var destinationResponse response.DestinationResponse

	id := c.Param("id")

	// Fetch the destination details with related data
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

	// Convert Facilities field to an array of strings and trim spaces
	var facilitiesArray []string
	if destination.Facilities != "" {
		// Split the Facilities string and trim spaces for each element
		facilitiesArray = make([]string, 0)
		for _, facility := range strings.Split(destination.Facilities, ",") {
			facilitiesArray = append(facilitiesArray, strings.TrimSpace(facility))
		}
	}

	// Populate the response struct with the destination details
	destinationResponse = response.DestinationResponse{
		ID:               destination.ID,
		Name:             destination.Name,
		City:             response.City{ID: destination.City.ID, Name: destination.City.Name},
		Position:         destination.Position,
		Address:          destination.Address,
		OperationalHours: destination.OperationalHours,
		TicketPrice:      destination.TicketPrice,
		Category:         destination.Category,
		Description:      destination.Description,
		Facilities:       facilitiesArray, // Facilities as array of strings
		CreatedAt:        destination.CreatedAt,
		Images:           convertImagesToResponse(destination.Images),
		VideoContents:    convertVideosToResponse(destination.VideoContents),
	}

	// Return the response
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":     "Destination details fetched successfully",
		"destination": destinationResponse,
	})
}

func GetMostViewedVideoContent(c echo.Context) error {
	var results []struct {
		ID          uint   `json:"id"`
		Name        string `json:"name"`
		Address     string `json:"address"`
		Description string `json:"description"`
		ViewCount   int64  `json:"view_count"`
	}

	// Query to join Destination with VideoContentView and calculate view counts
	err := config.DB.Table("destinations").
		Select("destinations.id, destinations.name, destinations.address, destinations.description, COUNT(video_content_views.id) as view_count").
		Joins("LEFT JOIN video_content_views ON destinations.id = video_content_views.destination_id").
		Group("destinations.id").
		Order("view_count DESC").
		Scan(&results).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch data"})
	}

	// Convert the result to a format similar to DestinationResponse
	var responseResults []map[string]interface{}
	for _, result := range results {
		var videos []models.VideoContent
		_ = config.DB.
			Find(&videos, "destination_id = ?", result.ID)

		responseResults = append(responseResults, map[string]interface{}{
			"id":          result.ID,
			"name":        result.Name,
			"address":     result.Address,
			"description": result.Description,
			"view_count":  result.ViewCount,
			"videos":      videos,
		})
	}

	// Return the result as JSON
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Destinations with view count fetched successfully",
		"data":    responseResults,
	})
}

func convertImagesToResponse(images []models.Image) []response.Image {
	var imageResponses []response.Image
	for _, img := range images {
		imageResponses = append(imageResponses, response.Image{
			DestinationID: img.DestinationID,
			URL:           img.URL,
		})
	}
	return imageResponses
}

func convertVideosToResponse(videos []models.VideoContent) []response.VideoContent {
	var videoResponses []response.VideoContent
	for _, video := range videos {
		videoResponses = append(videoResponses, response.VideoContent{
			DestinationID: video.DestinationID,
			Title:         video.Title,
			URL:           video.URL,
			Description:   video.Description,
		})
	}
	return videoResponses
}

func GetPersonalizedDestinationByUser(c echo.Context) error {
	var destinations []models.Destination
	var destinationResponses []response.DestinationResponse

	userID := c.QueryParam("user_id")

	var user models.User
	result := config.DB.First(&user, "id = ?", userID)
	if result.Error != nil || user.ID == 0 {
		response := helper.APIResponse("User not found", http.StatusUnauthorized, "error", nil)
		return c.JSON(http.StatusUnauthorized, response)
	}

	query := config.DB.
		Preload("City").
		Preload("Images").
		Preload("VideoContents")

	categories := strings.Split(user.Category, ",")
	for i := range categories {
		categories[i] = strings.ToUpper(strings.TrimSpace(categories[i]))
	}
	query = query.Where("category IN ?", categories)

	err := query.Find(&destinations).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch destinations"})
	}

	for _, dest := range destinations {
		var facilitiesArray []string
		if dest.Facilities != "" {
			facilitiesArray = make([]string, 0)
			for _, facility := range strings.Split(dest.Facilities, ",") {
				facilitiesArray = append(facilitiesArray, strings.TrimSpace(facility))
			}
		}

		destinationResponse := response.DestinationResponse{
			ID:               dest.ID,
			Name:             dest.Name,
			City:             response.City{ID: dest.City.ID, Name: dest.City.Name},
			Position:         dest.Position,
			Address:          dest.Address,
			OperationalHours: dest.OperationalHours,
			TicketPrice:      dest.TicketPrice,
			Category:         dest.Category,
			Description:      dest.Description,
			Facilities:       facilitiesArray,
			CreatedAt:        dest.CreatedAt,
			Images:           convertImagesToResponse(dest.Images),
			VideoContents:    convertVideosToResponse(dest.VideoContents),
		}

		destinationResponses = append(destinationResponses, destinationResponse)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Destinations fetched successfully",
		"data":    destinationResponses,
	})
}

func CreateDestinationAssetsHandler(c echo.Context) error {
	var input CreateDestinationAssetsInput

	// Bind input
	if err := c.Bind(&input); err != nil {
		response := helper.APIResponse("Invalid request", http.StatusBadRequest, "error", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	var destination models.Destination
	result := config.DB.First(&destination, "id = ?", input.DestinationID)
	if result.Error != nil || destination.ID == 0 {
		response := helper.APIResponse("Destination not found", http.StatusUnauthorized, "error", nil)
		return c.JSON(http.StatusUnauthorized, response)
	}

	for i := 0; i < len(input.Images); i++ {
		image := new(models.Image)
		image.DestinationID = destination.ID
		image.URL = input.Images[i]
		if err := config.DB.Create(image).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to add image"})
		}
	}

	//video
	for i := 0; i < len(input.VideoContents); i++ {
		video := new(models.VideoContent)
		video.DestinationID = destination.ID
		video.URL = input.VideoContents[i].Url
		video.Title = input.VideoContents[i].Title
		video.Description = input.VideoContents[i].Description
		if err := config.DB.Create(video).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to add image"})
		}
	}

	response := helper.APIResponse("Create Destination Assets success", http.StatusOK, "success", destination)
	return c.JSON(http.StatusOK, response)
}

func GetAllVideoContents(c echo.Context) error {
	var videos []models.VideoContent

	err := config.DB.Find(&videos).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch videos"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Video Contents fetched successfully",
		"videos":  convertVideosToResponse(videos),
	})
}
