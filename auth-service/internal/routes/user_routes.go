package routes

import (
	"auth-service/internal/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes registers user-related routes to the router group.
func RegisterUserRoutes(r *gin.RouterGroup, userHandler *handlers.UserHandler) {
	userGroup := r.Group("/users")
	{
		userGroup.POST("/register", userHandler.Register)
		userGroup.POST("/login", userHandler.Login)
	}
}
