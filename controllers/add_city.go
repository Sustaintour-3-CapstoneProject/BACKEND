package controllers

import (
	"backend/config"
	"backend/models"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CreateCity(c echo.Context) error {
	// Deklarasi struct input
	type CityInput struct {
		Name string `json:"name"`
	}

	// Decode body request
	var input CityInput
	if err := json.NewDecoder(c.Request().Body).Decode(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid JSON body"})
	}

	// Validasi input
	if input.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "City name is required"})
	}

	// Cek apakah kota sudah ada
	var existingCity models.City
	if err := config.DB.Where("name = ?", input.Name).First(&existingCity).Error; err == nil {
		return c.JSON(http.StatusConflict, map[string]string{"message": "City already exists"})
	}

	// Simpan kota baru ke database
	city := models.City{Name: input.Name}
	if err := config.DB.Create(&city).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create city"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "City created successfully",
		"city":    city,
	})
}

func GetCity(c echo.Context) error {
	// Mendapatkan semua data kota dari database
	var cities []models.City
	if err := config.DB.Find(&cities).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch cities"})
	}

	// Mengembalikan data kota dalam format JSON
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "City fetched successfully",
		"cities":  cities,
	})
}
