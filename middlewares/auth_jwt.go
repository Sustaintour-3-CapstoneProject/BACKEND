package middlewares

import (
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

// Middleware untuk memverifikasi role admin
func AdminOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Ambil token dari header Authorization
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Authorization header is required"})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})
		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "invalid or expired token"})
		}

		// Verifikasi klaim token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "invalid token claims"})
		}

		// Periksa role di klaim token
		role, ok := claims["role"].(string)
		if !ok || role != "admin" {
			return c.JSON(http.StatusForbidden, map[string]string{"message": "Access forbidden: only admins are allowed"})
		}

		return next(c) // Lanjutkan ke handler jika validasi role sukses
	}
}
