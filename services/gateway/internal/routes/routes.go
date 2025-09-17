package routes

import (
	"gateway/internal/handlers"

	"github.com/gin-gonic/gin"
)

// SetupAuthRoutes sets up authentication routes
func SetupAuthRoutes(
	router *gin.Engine,
	authHandler *handlers.AuthHandler,
	twoFAHandler *handlers.TwoFAHandler,
	userHandler *handlers.UserHandler,
) {
	api := router.Group("/api/v1")

	// Health check
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "gateway"})
	})

	// Authentication routes
	auth := api.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
		auth.POST("/logout", authHandler.Logout)
		auth.POST("/refresh", authHandler.RefreshToken)
		auth.GET("/validate", authHandler.ValidateToken)
	}

	// Two-Factor Authentication routes
	twoFA := api.Group("/2fa")
	{
		twoFA.POST("/enable", twoFAHandler.Enable2FA)
		twoFA.POST("/disable", twoFAHandler.DisableTwoFA)
		twoFA.POST("/verify", twoFAHandler.VerifyTwoFA)
	}

	// User management routes
	users := api.Group("/users")
	{
		users.GET("/:userId", userHandler.GetUserProfile)
		users.PUT("/:userId", userHandler.UpdateUserProfile)
	}
}
