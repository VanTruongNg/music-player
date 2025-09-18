package handlers

import (
	"auth-service/internal/utils/jwt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type JWKSHandler struct {
	jwtService jwt.JWTService
}

func NewJWKSHandler(jwtService jwt.JWTService) *JWKSHandler {
	return &JWKSHandler{
		jwtService: jwtService,
	}
}

func (h *JWKSHandler) GetJWKS(c *gin.Context) {
	jwks, err := h.jwtService.GetJWKS()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to retrieve JWKS",
		})
		return
	}

	c.Header("Content-Type", "application/json")
	c.Header("Cache-Control", "public, max-age=3600")
	c.JSON(http.StatusOK, jwks)
}
