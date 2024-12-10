package controllers

import (
	"backend/config"
	"backend/helper"
	"backend/models"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// Struct untuk validasi input login
type LoginInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
}

// Struct untuk validasi input registrasi
type RegisterInput struct {
	Username  string   `json:"username" validate:"required,alphanum"`
	FirstName string   `json:"first_name" validate:"required"`
	LastName  string   `json:"last_name" validate:"required"`
	Email     string   `json:"email" validate:"required,email"`
	City      string   `json:"city" validate:"required"`
	Password  string   `json:"password" validate:"required,min=6"`
	Role      string   `json:"role,omitempty" validate:"omitempty,oneof=admin user"`
	Category  []string `json:"category"`
}

// LoginHandler menangani proses login
func LoginHandler(c echo.Context) error {
	var input LoginInput

	// Bind input
	if err := c.Bind(&input); err != nil {
		response := helper.APIResponse("Invalid request", http.StatusBadRequest, "error", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	// Validasi input
	if err := helper.ValidateInput(&input); err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponse("Validation error", http.StatusBadRequest, "error", errors)
		return c.JSON(http.StatusBadRequest, response)
	}

	// Cari user berdasarkan username
	var user models.User
	result := config.DB.First(&user, "username = ?", input.Username)
	if result.Error != nil || user.ID == 0 {
		response := helper.APIResponse("Username not found", http.StatusUnauthorized, "error", nil)
		return c.JSON(http.StatusUnauthorized, response)
	}

	// Cek password
	if !helper.CheckPasswordHash(input.Password, user.Password) {
		response := helper.APIResponse("Incorrect password", http.StatusUnauthorized, "error", nil)
		return c.JSON(http.StatusUnauthorized, response)
	}

	// Generate token JWT
	token, err := helper.GenerateJWT(user.ID, user.Username, user.Role)
	if err != nil {
		response := helper.APIResponse("Failed to generate token", http.StatusInternalServerError, "error", nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	// Response data
	data := map[string]interface{}{
		"id_user":    user.ID,
		"username":   user.Username,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":      user.Email,
		"city":       user.City,
		"role":       user.Role,
		"token":      token,
	}

	response := helper.APIResponse("Login successful", http.StatusOK, "success", data)
	return c.JSON(http.StatusOK, response)
}

func LogoutHandler(c echo.Context) error {
	cookie := &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	}
	c.SetCookie(cookie)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Berhasil Logout",
	})
}

// RegisterHandler menangani proses registrasi
func RegisterHandler(c echo.Context) error {
	var input RegisterInput

	// Bind input
	if err := c.Bind(&input); err != nil {
		response := helper.APIResponse("Invalid request", http.StatusBadRequest, "error", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	// Validasi input
	if err := helper.ValidateInput(&input); err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.APIResponse("Validation error", http.StatusBadRequest, "error", errors)
		return c.JSON(http.StatusBadRequest, response)
	}

	// Hash password
	hashedPassword, err := helper.HashPassword(input.Password)
	if err != nil {
		response := helper.APIResponse("Failed to hash password", http.StatusInternalServerError, "error", nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	// Default role
	if input.Role == "" {
		input.Role = "user"
	}

	categoryString := strings.ToUpper(strings.Join(input.Category, ","))

	// Membuat objek user
	user := models.User{
		Username:  input.Username,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		City:      input.City,
		Password:  hashedPassword,
		Role:      input.Role,
		Category:  categoryString,
	}

	// Simpan ke database
	if result := config.DB.Create(&user); result.Error != nil {
		response := helper.APIResponse("Failed to register", http.StatusInternalServerError, "error", nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	// Generate token JWT
	token, err := helper.GenerateJWT(user.ID, user.Username, user.Role)
	if err != nil {
		response := helper.APIResponse("Failed to generate token", http.StatusInternalServerError, "error", nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	categories := strings.Split(user.Category, ",")
	for i, category := range categories {
		categories[i] = strings.ToUpper(strings.TrimSpace(category))
	}

	// Response data
	data := map[string]interface{}{
		"id_user":    user.ID,
		"username":   user.Username,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":      user.Email,
		"city":       user.City,
		"role":       user.Role,
		"token":      token,
		"category":   categories,
	}

	response := helper.APIResponse("Registration successful", http.StatusOK, "success", data)
	return c.JSON(http.StatusOK, response)
}
