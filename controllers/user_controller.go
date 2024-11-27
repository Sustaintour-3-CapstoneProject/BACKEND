package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"backend/config"
	"backend/middlewares"
	"backend/models"
)

func Register(c echo.Context) error {
	user := new(models.User)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid input"})
	}

	// Validasi input
	if user.Username == "" || user.FirstName == "" || user.LastName == "" || user.PhoneNumber == "" || user.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "semua kolom wajib diisi"})
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	if err := config.DB.Create(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "could not register user"})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "user registered successfully",
		"user": map[string]interface{}{
			"username":     user.Username,
			"first_name":   user.FirstName,
			"last_name":    user.LastName,
			"phone_number": user.PhoneNumber,
		},
	})
}

func Login(c echo.Context) error {
	userInput := new(models.User)
	if err := c.Bind(userInput); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid input"})
	}

	// Validasi input
	if userInput.Username == "" || userInput.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "username and password are required"})
	}

	// Cari user berdasarkan username
	var user models.User
	if err := config.DB.Where("username = ?", userInput.Username).First(&user).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "user tidak ditemukan"})
	}

	// Verifikasi password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "password tidak ditemukan"})
	}

	// Generate JWT
	token, err := middlewares.GenerateToken(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "could not generate token"})
	}

	// Response sukses
	return c.JSON(http.StatusOK, map[string]string{
		"message": "login successful",
		"token":   token,
	})
}
