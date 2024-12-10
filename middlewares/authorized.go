package middlewares

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func AuthorizedAccess(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Unauthorized Access"})
		}

		tokenString := strings.TrimPrefix(authHeader, bearerPrefix)

		secretKey := os.Getenv("JWT_SECRET_KEY")
		if secretKey == "" {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Unauthorized Access"})
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})
		if err != nil {
			var validationErr *jwt.ValidationError
			if errors.As(err, &validationErr) && validationErr.Errors&jwt.ValidationErrorExpired != 0 {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Unauthorized Access"})
			}
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Unauthorized Access"})
		}

		_, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Unauthorized Access"})
		}

		return next(c)
	}
}
