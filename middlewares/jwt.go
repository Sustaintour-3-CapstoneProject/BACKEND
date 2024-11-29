package middlewares

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateToken(userID uint, role string) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = userID
	claims["role"] = role                              // Menambahkan role ke dalam claims
	claims["exp"] = time.Now().AddDate(0, 1, 0).Unix() // Token berlaku 1 bulan dari sekarang

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
}
