package routes

import (
	"auth-service/internal/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterTwoFARoutes registers 2FA endpoints under user resource.
func RegisterTwoFARoutes(rg *gin.RouterGroup, handler *handlers.TwoFAHandler) {
	user := rg.Group("/users/:id/2fa")
	{
		user.POST("/setup", handler.Setup2FA)
		user.POST("/enable", handler.Enable2FA)
		user.POST("/verify", handler.Verify2FA)
		user.POST("/disable", handler.Disable2FA)
	}
}
