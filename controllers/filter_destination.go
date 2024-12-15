package controllers

import (
	"backend/config"
	"backend/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

// FilterDestinations godoc
// @Summary Filter destinations based on city and category
// @Description Retrieve destinations filtered by city and/or category
// @Tags Destinations
// @Accept json
// @Produce json
// @Param city query string false "City to filter destinations by"
// @Param category query array false "Categories to filter destinations by"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /destinations [get]
func FilterDestinations(c echo.Context) error {
	// Mengambil query parameter 'city' untuk filter berdasarkan kota
	city := c.QueryParam("city")

	// Mengambil query parameter 'category' untuk filter berdasarkan kategori
	categories := c.QueryParams()["category"] // Mendapatkan daftar kategori yang dipilih pengguna, bisa lebih dari satu

	var destinations []models.Destination
	query := config.DB.Preload("Images").Preload("VideoContents")

	// Filter berdasarkan kota
	if city != "" {
		query = query.Joins("JOIN cities ON destinations.city_id = cities.id").
			Where("cities.name = ?", city)
	}

	// Filter berdasarkan kategori
	if len(categories) > 0 {
		query = query.Where("category IN ?", categories)
	}

	// Ambil data destinasi yang sudah difilter
	if err := query.Find(&destinations).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch destinations"})
	}

	// Kembalikan hasil dalam bentuk JSON
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":      "Destinations filtered successfully",
		"destinations": destinations,
	})
}
