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
	"time"

	"github.com/labstack/echo/v4"
)

func CreateRoute(c echo.Context) error {
	// Decode JSON body
	jsonBody := new(request.CreateRouteInput)
	if err := json.NewDecoder(c.Request().Body).Decode(jsonBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid JSON body"})
	}

	var originCity models.City
	if err := config.DB.Where("name = ?", jsonBody.OriginCityName).First(&originCity).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Origin City not found"})
	}

	var destinationCity models.City
	if err := config.DB.Where("name = ?", jsonBody.DestinationCityName).First(&destinationCity).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Destination City not found"})
	}

	route := models.Route{
		UserID:              jsonBody.UserID,
		OriginCityName:      originCity.Name,
		DestinationCityName: destinationCity.Name,
		Distance:            jsonBody.Distance,
		Time:                jsonBody.Time,
		Cost:                jsonBody.Cost,
	}

	if err := config.DB.Create(&route).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create route"})
	}

	for i := 0; i < len(jsonBody.Destinations); i++ {
		routeDestination := models.RouteDestination{
			RouteID:       route.ID,
			DestinationID: jsonBody.Destinations[i],
			CreatedAt:     time.Now(),
		}

		config.DB.Create(&routeDestination)
	}

	return c.JSON(http.StatusOK, route)
}

func calculateDistance(originCity models.City, destinationCity models.City) float64 {
	lat1, _ := strconv.ParseFloat(originCity.Lat, 64)
	lon1, _ := strconv.ParseFloat(originCity.Long, 64)
	lat2, _ := strconv.ParseFloat(destinationCity.Lat, 64)
	lon2, _ := strconv.ParseFloat(destinationCity.Long, 64)

	return helper.Haversine(lat1, lon1, lat2, lon2)
}

func GetRouteByUser(c echo.Context) error {
	var routes []models.Route

	userID := c.QueryParam("user_id")

	var user models.User
	result := config.DB.First(&user, "id = ?", userID)
	if result.Error != nil || user.ID == 0 {
		response := helper.APIResponse("User not found", http.StatusUnauthorized, "error", nil)
		return c.JSON(http.StatusUnauthorized, response)
	}

	err := config.DB.Where("user_id = ?", userID).Find(&routes).Error
	if err != nil {

		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	var responses []response.RouteResponse

	for i := 0; i < len(routes); i++ {
		var routeDestinations []models.RouteDestination
		config.DB.Where("route_id = ?", routes[i].ID).Find(&routeDestinations)

		var destinations []models.Destination

		for j := 0; j < len(routeDestinations); j++ {
			var destination models.Destination
			config.DB.Where("id = ?", routeDestinations[j].DestinationID).Find(&destination)

			println(destination.Name)
			destinations = append(destinations, destination)
		}

		var response = response.RouteResponse{
			ID:                  routes[i].ID,
			UserID:              routes[i].UserID,
			OriginCityName:      routes[i].OriginCityName,
			DestinationCityName: routes[i].DestinationCityName,
			Distance:            routes[i].Distance,
			Time:                routes[i].Time,
			Cost:                routes[i].Cost,
			CreatedAt:           routes[i].CreatedAt,
			Destinations:        destinations,
		}

		responses = append(responses, response)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Routes fetched successfully",
		"data":    responses,
	})
}

func DeleteRoute(c echo.Context) error {
	id := c.Param("id")
	routeID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid route ID"})
	}

	var route models.Route
	if err := config.DB.First(&route, routeID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Route not found"})
	}

	tx := config.DB.Begin()

	if err := tx.Where("route_id = ?", route.ID).Delete(&route.Destinations).Error; err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to delete related routes destinations"})
	}

	if err := tx.Delete(&route).Error; err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to delete route"})
	}

	tx.Commit()

	return c.JSON(http.StatusOK, map[string]string{"message": "Route and related data successfully deleted"})
}

func GetDestinationsByRoute(c echo.Context) error {
	originCityName := c.QueryParam("origin")
	destinationCityName := c.QueryParam("destination")

	var originCity models.City
	if err := config.DB.Where("name = ?", originCityName).First(&originCity).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Origin City not found"})
	}

	var destinationCity models.City
	if err := config.DB.Where("name = ?", destinationCityName).First(&destinationCity).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Destination City not found"})
	}

	distance := calculateDistance(originCity, destinationCity)

	var IDs []uint
	IDs = append(IDs, originCity.ID)
	IDs = append(IDs, destinationCity.ID)

	var destinations []models.Destination
	query := config.DB.Where("city_id IN ?", IDs)

	err := query.Find(&destinations).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch destinations"})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"distance":     distance,
		"destinations": destinations,
	})
}
