package routes

import (
	"backend/controllers"
	"backend/middlewares"

	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo) {
	// e.POST("/register", controllers.Register)
	// e.POST("/login", controllers.Login)
	// e.POST("/logout", controllers.Logout)

	// // Route untuk admin hanya bisa diakses oleh admin
	// e.POST("/admin-only", controllers.AdminDashboard, middlewares.AdminOnly)

	// // Route untuk user biasa
	// e.GET("/user-profile", controllers.UserProfile)

	// // Route untuk admin (hanya bisa diakses oleh admin)
	// adminGroup := e.Group("/admin")
	// adminGroup.Use(middlewares.RoleBasedAccess([]string{"admin"})) // Hanya admin yang boleh akses
	// adminGroup.GET("/dashboard", controllers.AdminDashboard)

	// // Route untuk user (hanya bisa diakses oleh user biasa)
	// userGroup := e.Group("/user")
	// userGroup.Use(middlewares.RoleBasedAccess([]string{"user", "admin"})) // User dan admin boleh akses
	// userGroup.GET("/profile", controllers.UserProfile)

	// Rute khusus admin dengan middleware AdminOnly
	adminGroup := e.Group("/api/v1/admin")
	adminGroup.Use(middlewares.AdminOnly)

	// Rute autentikasi
	e.POST("/api/v1/register", controllers.RegisterHandler)
	e.POST("/api/v1/login", controllers.LoginHandler)

}
