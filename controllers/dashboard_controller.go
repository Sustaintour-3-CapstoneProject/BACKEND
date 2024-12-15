package controllers

import (
	"backend/config"
	"backend/models"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type MonthlyUserCount struct {
	Month string
	Count int
}

// GetDashboardDataHandler godoc
// @Summary Fetch dashboard data summary
// @Description Retrieve a summary of users, destinations, video content, and destination categories for the dashboard
// @Tags Dashboard
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /dashboard/count-data [get]
func GetDashboardDataHandler(c echo.Context) error {
	var users []models.User
	var destinations []models.Destination
	var videoContents []models.VideoContent

	var natureCount, cultureCount, ecotourismCount int64

	// Fetch users, destinations, and video contents
	config.DB.Find(&users)
	config.DB.Find(&destinations)
	config.DB.Find(&videoContents)

	// Count destinations by categories
	config.DB.Model(&models.Destination{}).Where("category = ?", "Nature").Count(&natureCount)
	config.DB.Model(&models.Destination{}).Where("category = ?", "Culture").Count(&cultureCount)
	config.DB.Model(&models.Destination{}).Where("category = ?", "Ecotourism").Count(&ecotourismCount)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Dashboard data fetched successfully",
		"data": map[string]any{
			"user":         len(users),
			"destination":  len(destinations),
			"videoContent": len(videoContents),
			"destinationCategories": map[string]int64{
				"Nature":     natureCount,
				"Culture":    cultureCount,
				"Ecotourism": ecotourismCount,
			},
		},
	})
}

// GetDashboardGraphicDataHandler godoc
// @Summary Fetch monthly user registration data
// @Description Retrieve the number of user registrations grouped by month
// @Tags Dashboard
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /dashboard/graphic [get]
func GetDashboardGraphicDataHandler(c echo.Context) error {
	var userMonthlyCounts []MonthlyUserCount

	// Query for grouping user count by month with month abbreviation
	err := config.DB.Model(&models.User{}).
		Select("DATE_FORMAT(created_at, '%b') AS month, COUNT(*) AS count").
		Group("month").
		Scan(&userMonthlyCounts).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to fetch user count by month",
			"error":   err.Error(),
		})
	}

	// Predefined list of months
	months := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}

	// Initialize map with all months set to 0
	monthlyUserData := make(map[string]int)
	for _, month := range months {
		monthlyUserData[month] = 0
	}

	// Populate the map with data from the query
	for _, userCount := range userMonthlyCounts {
		monthlyUserData[userCount.Month] = userCount.Count
	}

	// Return response
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Dashboard data fetched successfully",
		"data":    monthlyUserData,
	})
}
