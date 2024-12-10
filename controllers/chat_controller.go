package controllers

import (
	"backend/helper"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Input struct {
	Message string `json:"message"`
}

func ChatHandler(c echo.Context) error {
	var input Input
	if err := json.NewDecoder(c.Request().Body).Decode(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid JSON body"})
	}

	response, err := helper.CallGeminiAPI(input.Message)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Chat successfully sent!",
		"data":    response,
	})
}
