package routes

import (
	"auth-service/internal/handlers"
	"auth-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterTwoFARoutes registers 2FA endpoints under user resource.
func RegisterTwoFARoutes(rg *gin.RouterGroup, handler *handlers.TwoFAHandler, authMiddleware *middleware.AuthMiddleware) {
	user := rg.Group("/auth/:id/2fa", authMiddleware.RequireAuth())
	{
		user.POST("/setup", handler.Setup2FA)
		user.POST("/enable", handler.Enable2FA)
		user.POST("/verify", handler.Verify2FA)
		user.POST("/disable", handler.Disable2FA)
	}
}
