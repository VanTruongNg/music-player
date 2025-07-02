package routes

import (
	"auth-service/internal/handlers"
	"auth-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes registers user-related routes to the router group.
func RegisterUserRoutes(r *gin.RouterGroup, userHandler *handlers.UserHandler, authMiddleware *middleware.AuthMiddleware) {
	userGroup := r.Group("/auth")
	{
		// Public routes - no authentication required
		userGroup.POST("/register", userHandler.Register)
		userGroup.POST("/login", userHandler.Login)

		// Protected routes - authentication required
		protectedGroup := userGroup.Group("", authMiddleware.RequireAuth())
		{
			protectedGroup.GET("/me", userHandler.GetMe)
		}
	}
}
