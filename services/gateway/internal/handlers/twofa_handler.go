package handlers

import (
	"context"
	"gateway/configs"
	"gateway/internal/utils"
	authv1 "music-player/api/proto/auth/v1"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type TwoFAHandler interface {
	Setup2FA(c *gin.Context)
	Enable2FA(c *gin.Context)
}

// TwoFAHandler handles 2FA-related HTTP requests
type twoFAHandler struct {
	grpcClients *configs.GRPCClients
}

// NewTwoFAHandler creates a new TwoFAHandler
func NewTwoFAHandler(grpcClients *configs.GRPCClients) TwoFAHandler {
	return &twoFAHandler{
		grpcClients: grpcClients,
	}
}

func (h *twoFAHandler) Setup2FA(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Fail(c, 401, "UNAUTHORIZED", "User ID not found in context")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.grpcClients.AuthClient.SetupTwoFA(ctx, &authv1.SetupTwoFARequest{
		UserId: userID.(string),
	})

	if err != nil {
		utils.Fail(c, http.StatusInternalServerError, "SETUP_2FA_FAILED", err.Error())
		return
	}

	if !resp.Success {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": resp.Message,
		})
		return
	}

	utils.Success(c, 200, resp)
}

func (h *twoFAHandler) Enable2FA(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		utils.Fail(c, http.StatusBadGateway, "UNAUTHORIZED", "User ID not found in context")
		return
	}

	var req struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request payload")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.grpcClients.AuthClient.EnableTwoFA(ctx, &authv1.EnableTwoFARequest{
		UserId: userId.(string),
		Code:   req.Code,
	})
	if err != nil {
		utils.Fail(c, http.StatusInternalServerError, "ENABLE_2FA_FAILED", err.Error())
		return
	}

	if !resp.Success {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": resp.Message,
		})
		return
	}

	utils.Success(c, http.StatusOK, resp)
}
