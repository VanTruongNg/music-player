package routes

import (
	"gateway/internal/handlers"
	"gateway/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupAuthRoutes sets up authentication routes with JWT middleware
func SetupAuthRoutes(
	router *gin.Engine,
	authHandler *handlers.AuthHandler,
	twoFAHandler handlers.TwoFAHandler,
	userHandler handlers.UserHandler,
	authMiddleware *middleware.AuthMiddleware,
) {
	api := router.Group("/api/v1")

	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "gateway"})
	})

	auth := api.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)

		// Protected auth routes
		authProtected := auth.Group("")
		authProtected.Use(authMiddleware.RequireAuth())
		{
			authProtected.POST("/logout", authHandler.Logout)
			authProtected.POST("/refresh", authHandler.RefreshToken)
			authProtected.GET("/validate", authHandler.ValidateToken)
		}
	}

	// Two-Factor Authentication routes (all protected)
	twoFA := api.Group("/2fa")
	twoFA.Use(authMiddleware.RequireAuth())
	{
		twoFA.POST("/setup", twoFAHandler.Setup2FA)
		twoFA.POST("/enable", twoFAHandler.Enable2FA)
	}

	// User management routes (all protected)
	users := api.Group("/users")
	users.Use(authMiddleware.RequireAuth())
	{
		users.GET("", userHandler.GetUserProfile)
	}
}
