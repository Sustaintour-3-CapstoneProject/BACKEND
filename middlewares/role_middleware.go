package middlewares

import (
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

// Middleware untuk memverifikasi role (admin atau user)
func RoleBasedAccess(allowedRoles []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Authorization header is required"})
			}

			// Mengambil token dari header Authorization
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

			// Pastikan klaim "role" ada dan memiliki tipe string
			role, ok := claims["role"].(string)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "role claim is missing"})
			}

			// Periksa apakah role termasuk dalam allowedRoles
			roleAllowed := false
			for _, allowedRole := range allowedRoles {
				if role == allowedRole {
					roleAllowed = true
					break
				}
			}

			if !roleAllowed {
				return c.JSON(http.StatusForbidden, map[string]string{"message": "access forbidden: insufficient role"})
			}

			return next(c) // Lanjutkan ke handler berikutnya jika role valid
		}
	}
}
