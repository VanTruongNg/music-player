package handlers

import (
	"context"
	"net/http"
	"time"

	"gateway/configs"
	authv1 "music-player/api/proto/auth/v1"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	GetUserProfile(c *gin.Context)
}
type userHandler struct {
	grpcClients *configs.GRPCClients
}

func NewUserHandler(grpcClients *configs.GRPCClients) UserHandler {
	return &userHandler{
		grpcClients: grpcClients,
	}
}

func (h *userHandler) GetUserProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.grpcClients.AuthClient.GetUserProfile(ctx, &authv1.GetUserProfileRequest{
		UserId: userID.(string),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user profile"})
		return
	}

	c.JSON(http.StatusOK, resp)
}
