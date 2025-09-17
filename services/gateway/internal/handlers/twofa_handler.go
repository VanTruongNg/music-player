package handlers

import (
	"context"
	"net/http"
	"time"

	"gateway/configs"
	authv1 "music-player/api/proto/auth/v1"

	"github.com/gin-gonic/gin"
)

// TwoFAHandler handles 2FA-related HTTP requests
type TwoFAHandler struct {
	grpcClients *configs.GRPCClients
}

// NewTwoFAHandler creates a new TwoFAHandler
func NewTwoFAHandler(grpcClients *configs.GRPCClients) *TwoFAHandler {
	return &TwoFAHandler{
		grpcClients: grpcClients,
	}
}

// Enable2FA enables two-factor authentication for a user
func (h *TwoFAHandler) Enable2FA(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "User not authenticated",
		})
		return
	}

	userID, ok := userIDInterface.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Invalid user ID format",
		})
		return
	}

	// Create context with timeout for gRPC call
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	// Call auth service via gRPC
	resp, err := h.grpcClients.AuthClient.EnableTwoFA(ctx, &authv1.EnableTwoFARequest{
		UserId: userID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Authentication service unavailable",
			"error":   err.Error(),
		})
		return
	}

	if !resp.Success {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": resp.Message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"message":   resp.Message,
		"qrCodeUrl": resp.QrCodeUrl,
		"secretKey": resp.SecretKey,
	})
}

// DisableTwoFA disables two-factor authentication
func (h *TwoFAHandler) DisableTwoFA(c *gin.Context) {
	var req struct {
		UserID string `json:"userId" binding:"required"`
		Code   string `json:"code" binding:"required,len=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}

	// Create context with timeout for gRPC call
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	// Call auth service via gRPC
	resp, err := h.grpcClients.AuthClient.DisableTwoFA(ctx, &authv1.DisableTwoFARequest{
		UserId:    req.UserID,
		TwoFaCode: req.Code,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Authentication service unavailable",
			"error":   err.Error(),
		})
		return
	}

	if !resp.Success {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": resp.Message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": resp.Message,
	})
}

// VerifyTwoFA verifies a 2FA code (for general verification)
func (h *TwoFAHandler) VerifyTwoFA(c *gin.Context) {
	var req struct {
		UserID string `json:"userId" binding:"required"`
		Code   string `json:"code" binding:"required,len=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}

	// Create context with timeout for gRPC call
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// Call auth service via gRPC
	resp, err := h.grpcClients.AuthClient.VerifyTwoFA(ctx, &authv1.VerifyTwoFARequest{
		UserId:    req.UserID,
		TwoFaCode: req.Code,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Authentication service unavailable",
			"error":   err.Error(),
		})
		return
	}

	if !resp.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": resp.Message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": resp.Message,
		"valid":   resp.Valid,
	})
}
