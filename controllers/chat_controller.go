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

// ChatHandler godoc
// @Summary Send a chat message to Gemini API
// @Description Process a chat message and get a response from Gemini API
// @Tags Chat
// @Accept json
// @Produce json
// @Param input body Input true "Chat Message"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /chat [post]
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
