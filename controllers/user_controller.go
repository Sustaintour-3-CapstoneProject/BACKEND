package controllers

import (
	"backend/config"
	"backend/helper"
	"backend/models"
	"backend/response"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

// Struct untuk validasi input login
type LoginInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
}

// LoginHandler godoc
// @Summary User login
// @Description Handle user login by verifying username and password, and return a JWT token.
// @Tags User
// @Accept json
// @Produce json
// @Param input body LoginInput true "Login credentials"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /login [post]
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

	var file string

	if user.File != "" {
		file = os.Getenv("APP_BASE") + "/" + user.File
	}

	// Response data
	data := map[string]interface{}{
		"id_user":      user.ID,
		"username":     user.Username,
		"first_name":   user.FirstName,
		"last_name":    user.LastName,
		"email":        user.Email,
		"city":         user.City,
		"role":         user.Role,
		"token":        token,
		"file":         file,
		"phone_number": user.PhoneNumber,
		"gender":       user.Gender,
	}

	response := helper.APIResponse("Login successful", http.StatusOK, "success", data)
	return c.JSON(http.StatusOK, response)
}

// LogoutHandler godoc
// @Summary Log out a user
// @Description Process the logout request and return a success message
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /logout [get]
func LogoutHandler(c echo.Context) error {

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Berhasil Logout",
	})
}

// Struct untuk validasi input registrasi
type RegisterInput struct {
	Username  string `json:"username" validate:"required,alphanum"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	City      string `json:"city" validate:"required"`
	Password  string `json:"password" validate:"required,min=6"`
	Role      string `json:"role,omitempty" validate:"omitempty,oneof=admin user"`
}

// RegisterHandler godoc
// @Summary User registration
// @Description Handle user registration by validating input and creating a new user in the database.
// @Tags User
// @Accept json
// @Produce json
// @Param input body RegisterInput true "Registration details"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /register [post]
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

	hashedPassword, err := helper.HashPassword(input.Password)
	if err != nil {
		response := helper.APIResponse("Failed to hash password", http.StatusInternalServerError, "error", nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	var role string
	if input.Role == "" {
		role = "user"
	} else {
		role = input.Role
	}

	// Create user object
	user := models.User{
		Username:    input.Username,
		FirstName:   input.FirstName,
		LastName:    input.LastName,
		Email:       input.Email,
		City:        input.City,
		Password:    hashedPassword,
		Role:        role,
		File:        "",
		PhoneNumber: "",
		Gender:      "",
	}

	// Save to the database
	if result := config.DB.Create(&user); result.Error != nil {
		response := helper.APIResponse("Failed to register", http.StatusInternalServerError, "error", nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	// Generate JWT token
	token, err := helper.GenerateJWT(user.ID, user.Username, user.Role)
	if err != nil {
		response := helper.APIResponse("Failed to generate token", http.StatusInternalServerError, "error", nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	// Response data
	data := map[string]interface{}{
		"id_user":      user.ID,
		"username":     user.Username,
		"first_name":   user.FirstName,
		"last_name":    user.LastName,
		"email":        user.Email,
		"city":         user.City,
		"role":         user.Role,
		"file":         user.File,
		"token":        token,
		"phone_number": user.PhoneNumber,
		"gender":       user.Gender,
	}

	response := helper.APIResponse("Registration successful", http.StatusOK, "success", data)
	return c.JSON(http.StatusOK, response)
}

type UserCategoryInput struct {
	UserID   int      `json:"userID"`
	Category []string `json:"category"`
}

// CreateUserCategoryHandler godoc
// @Summary Create or update user categories
// @Description Assign categories to a user by updating their profile.
// @Tags User
// @Accept json
// @Produce json
// @Param input body UserCategoryInput true "User categories"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /user/category [post]
func CreateUserCategoryHandler(c echo.Context) error {
	var input UserCategoryInput

	// Bind input
	if err := c.Bind(&input); err != nil {
		response := helper.APIResponse("Invalid request", http.StatusBadRequest, "error", nil)
		return c.JSON(http.StatusBadRequest, response)
	}

	var user models.User
	result := config.DB.First(&user, "id = ?", input.UserID)
	if result.Error != nil || user.ID == 0 {
		response := helper.APIResponse("User not found", http.StatusUnauthorized, "error", nil)
		return c.JSON(http.StatusUnauthorized, response)
	}

	categoryString := strings.ToUpper(strings.Join(input.Category, ","))

	user.Category = categoryString

	if result := config.DB.Save(&user); result.Error != nil {
		response := helper.APIResponse("Failed to update user", http.StatusInternalServerError, "error", nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	response := helper.APIResponse("Create User Category success", http.StatusOK, "success", user)
	return c.JSON(http.StatusOK, response)
}

// GetAllUserHandler godoc
// @Summary Get all users
// @Description Retrieve a list of users, optionally filtered by name.
// @Tags User
// @Accept json
// @Produce json
// @Param name query string false "Search by first or last name"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /user/{id} [get]
func GetAllUserHandler(c echo.Context) error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var users []models.User

	queryName := c.QueryParam("name")

	query := config.DB

	if queryName != "" {
		query = query.Where("first_name LIKE ? OR last_name LIKE ?", "%"+queryName+"%", "%"+queryName+"%")
	}

	query.Find(&users)

	var responses []response.UserResponse

	for i := 0; i < len(users); i++ {
		var file string

		if users[i].File != "" {
			file = os.Getenv("APP_BASE") + "/" + users[i].File
		} else {
			file = "https://static-00.iconduck.com/assets.00/profile-default-icon-2048x2045-u3j7s5nj.png"
		}

		var response = response.UserResponse{
			Username:    users[i].Username,
			FirstName:   users[i].FirstName,
			LastName:    users[i].LastName,
			Email:       users[i].Email,
			City:        users[i].City,
			Role:        users[i].Role,
			Category:    users[i].Category,
			File:        file,
			PhoneNumber: users[i].PhoneNumber,
			Gender:      users[i].Gender,
		}

		responses = append(responses, response)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Users fetched successfully",
		"data":    responses,
	})
}

// GetDetailUserHandler godoc
// @Summary Get user details
// @Description Retrieve detailed information for a specific user by ID.
// @Tags User
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /user/{id} [get]
func GetDetailUserHandler(c echo.Context) error {
	var user models.User

	id := c.Param("id")

	config.DB.Where("id = ?", id).Find(&user)

	var file string

	if user.File != "" {
		file = os.Getenv("APP_BASE") + "/" + user.File
	} else {
		file = "https://static-00.iconduck.com/assets.00/profile-default-icon-2048x2045-u3j7s5nj.png"
	}

	var response = response.UserResponse{
		Username:    user.Username,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		City:        user.City,
		Role:        user.Role,
		Category:    user.Category,
		File:        file,
		PhoneNumber: user.PhoneNumber,
		Gender:      user.Gender,
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "User fetched successfully",
		"data":    response,
	})
}

// EditUserHandler godoc
// @Summary Edit user profile
// @Description Update user information such as username, email, password, and categories.
// @Tags User
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "User ID"
// @Param username formData string false "Username"
// @Param first_name formData string false "First Name"
// @Param last_name formData string false "Last Name"
// @Param email formData string false "Email"
// @Param city formData string false "City"
// @Param password formData string false "Password"
// @Param role formData string false "Role"
// @Param phone_number formData string false "Phone Number"
// @Param gender formData string false "Gender"
// @Param file formData file false "Profile Image"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users/{id} [put]
func EditUserHandler(c echo.Context) error {
	// Get user ID from the request (e.g., from URL parameter or JWT claims)
	id := c.Param("id")

	// Fetch user from the database
	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	// Update fields based on the form data
	username := c.FormValue("username")
	firstName := c.FormValue("first_name")
	lastName := c.FormValue("last_name")
	email := c.FormValue("email")
	city := c.FormValue("city")
	password := c.FormValue("password")
	role := c.FormValue("role")
	phoneNumber := c.FormValue("phone_number")
	gender := c.FormValue("gender")

	// Handle empty required fields
	if username == "" || firstName == "" || lastName == "" || email == "" || city == "" || phoneNumber == "" || gender == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "All fields except password are required"})
	}

	var existUserByUsername models.User
	config.DB.Where("username = ? AND id != ?", username, id).Find(&existUserByUsername)

	if existUserByUsername.ID != 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Username already used"})
	}

	var existUserByEmail models.User
	config.DB.Where("email = ? AND id != ?", email, id).Find(&existUserByEmail)

	if existUserByEmail.ID != 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Email already used"})
	}

	// Update optional fields
	user.Username = username
	user.FirstName = firstName
	user.LastName = lastName
	user.Email = email
	user.City = city
	user.PhoneNumber = phoneNumber
	user.Gender = gender

	// Handle role (optional, default to existing)
	if role != "" {
		user.Role = role
	}

	// Handle file upload (optional)
	file, err := c.FormFile("file")
	if err == nil {
		src, err := file.Open()
		if err != nil {
			response := helper.APIResponse("Failed to process file", http.StatusInternalServerError, "error", nil)
			return c.JSON(http.StatusInternalServerError, response)
		}
		defer src.Close()

		// Save uploaded file
		filePath := "assets/" + file.Filename
		dst, err := os.Create(filePath)
		if err != nil {
			response := helper.APIResponse("Failed to save file", http.StatusInternalServerError, "error", nil)
			return c.JSON(http.StatusInternalServerError, response)
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			response := helper.APIResponse("Failed to save file", http.StatusInternalServerError, "error", nil)
			return c.JSON(http.StatusInternalServerError, response)
		}
		user.File = filePath
	}

	// Hash password (if provided)
	if password != "" {
		hashedPassword, err := helper.HashPassword(password)
		if err != nil {
			response := helper.APIResponse("Failed to hash password", http.StatusInternalServerError, "error", nil)
			return c.JSON(http.StatusInternalServerError, response)
		}
		user.Password = hashedPassword
	}

	// Save updated user to the database
	if result := config.DB.Save(&user); result.Error != nil {
		response := helper.APIResponse("Failed to update user", http.StatusInternalServerError, "error", nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	// Generate a new token
	token, err := helper.GenerateJWT(user.ID, user.Username, user.Role)
	if err != nil {
		response := helper.APIResponse("Failed to generate token", http.StatusInternalServerError, "error", nil)
		return c.JSON(http.StatusInternalServerError, response)
	}

	// Prepare response data
	data := map[string]interface{}{
		"id_user":      user.ID,
		"username":     user.Username,
		"first_name":   user.FirstName,
		"last_name":    user.LastName,
		"email":        user.Email,
		"city":         user.City,
		"role":         user.Role,
		"file":         user.File,
		"token":        token,
		"phone_number": user.PhoneNumber,
		"gender":       user.Gender,
	}

	response := helper.APIResponse("User updated successfully", http.StatusOK, "success", data)
	return c.JSON(http.StatusOK, response)
}

func DeleteUser(c echo.Context) error {
	id := c.Param("id")
	userID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid user ID"})
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found"})
	}

	tx := config.DB.Begin()

	if err := tx.Delete(&user).Error; err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to delete user"})
	}

	tx.Commit()

	return c.JSON(http.StatusOK, map[string]string{"message": "User successfully deleted"})
}
