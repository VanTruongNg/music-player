package routes

import (
	"auth-service/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterJWKSRoutes(r *gin.RouterGroup, jwksHandler *handlers.JWKSHandler) {
	r.GET("/.well-known/jwks.json", jwksHandler.GetJWKS)
}
