package handlers

import (
	"context"
	"net/http"
	"time"

	"gateway/configs"
	authv1 "music-player/api/proto/auth/v1"

	"github.com/gin-gonic/gin"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	grpcClients *configs.GRPCClients
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(grpcClients *configs.GRPCClients) *UserHandler {
	return &UserHandler{
		grpcClients: grpcClients,
	}
}

// GetUserProfile gets user profile information
func (h *UserHandler) GetUserProfile(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "User ID is required",
		})
		return
	}

	// Create context with timeout for gRPC call
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// Call auth service via gRPC
	resp, err := h.grpcClients.AuthClient.GetUserProfile(ctx, &authv1.GetUserProfileRequest{
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
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": resp.Message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": resp.Message,
		"user": gin.H{
			"id":           resp.User.Id,
			"username":     resp.User.Username,
			"email":        resp.User.Email,
			"fullName":     resp.User.FullName,
			"twoFaEnabled": resp.User.TwoFaEnabled,
			"createdAt":    resp.User.CreatedAt,
			"updatedAt":    resp.User.UpdatedAt,
		},
	})
}

// UpdateUserProfile updates user profile information
func (h *UserHandler) UpdateUserProfile(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "User ID is required",
		})
		return
	}

	var req struct {
		FullName string `json:"fullName,omitempty"`
		Email    string `json:"email,omitempty"`
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
	resp, err := h.grpcClients.AuthClient.UpdateUserProfile(ctx, &authv1.UpdateUserProfileRequest{
		UserId:   userID,
		FullName: req.FullName,
		Email:    req.Email,
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
		"user": gin.H{
			"id":           resp.User.Id,
			"username":     resp.User.Username,
			"email":        resp.User.Email,
			"fullName":     resp.User.FullName,
			"twoFaEnabled": resp.User.TwoFaEnabled,
			"createdAt":    resp.User.CreatedAt,
			"updatedAt":    resp.User.UpdatedAt,
		},
	})
}
