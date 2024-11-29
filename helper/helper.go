package helper

import (
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Struct untuk API Response
type Response struct {
	Meta meta        `json:"meta"`
	Data interface{} `json:"data"`
}

type meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

// APIResponse untuk membungkus response API
func APIResponse(message string, code int, status string, data interface{}) Response {
	meta := meta{
		Message: message,
		Code:    code,
		Status:  status,
	}
	jsonresponse := Response{
		Meta: meta,
		Data: data,
	}
	return jsonresponse
}

// FormatValidationError memformat error validasi menjadi array string
func FormatValidationError(err error) []string {
	var errors []string
	for _, e := range err.(validator.ValidationErrors) {
		errors = append(errors, e.Error())
	}
	return errors
}

// ValidateInput memvalidasi struct input menggunakan library validator
func ValidateInput(input interface{}) error {
	validate := validator.New()
	return validate.Struct(input)
}

// HashPassword mengenkripsi password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash mencocokkan password dengan hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Struct untuk JWT Claims
type JWTClaims struct {
	Username string `json:"username"`
	UserID   uint   `json:"user_id"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateJWT membuat token JWT
func GenerateJWT(userID uint, username, role string) (string, error) {
	claims := &JWTClaims{
		Username: username,
		UserID:   userID,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
}
