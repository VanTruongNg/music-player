package handlers

import (
	"context"
	"net/http"
	"time"

	"gateway/configs"
	authv1 "music-player/api/proto/auth/v1"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	grpcClients *configs.GRPCClients
}

func NewAuthHandler(grpcClients *configs.GRPCClients) *AuthHandler {
	return &AuthHandler{
		grpcClients: grpcClients,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email     string `json:"email" binding:"required"`
		Password  string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	resp, err := h.grpcClients.AuthClient.Login(ctx, &authv1.LoginRequest{
		Email:     req.Email,
		Password:  req.Password,
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
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": resp.Message,
		})
		return
	}

	c.SetCookie(
		"refresh_token",
		resp.RefreshToken,
		int(resp.ExpiresIn)*24*7,
		"/",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"message":     resp.Message,
		"accessToken": resp.AccessToken,
		"expiresIn":   resp.ExpiresIn,
		"user": gin.H{
			"id":           resp.User.Id,
			"username":     resp.User.Username,
			"email":        resp.User.Email,
			"fullName":     resp.User.FullName,
			"twoFaEnabled": resp.User.TwoFaEnabled,
			"createdAt":    resp.User.CreatedAt,
		},
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required,min=3,max=32"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6,max=64"`
		FullName string `json:"fullName,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	resp, err := h.grpcClients.AuthClient.Register(ctx, &authv1.RegisterRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		FullName: req.FullName,
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

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": resp.Message,
		"userId":  resp.UserId,
	})
}

func (h *AuthHandler) ValidateToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Authorization header required",
		})
		return
	}

	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid authorization header format",
		})
		return
	}
	token := authHeader[len(bearerPrefix):]

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.grpcClients.AuthClient.ValidateToken(ctx, &authv1.ValidateTokenRequest{
		AccessToken: token,
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
		"user": gin.H{
			"id":       resp.User.Id,
			"username": resp.User.Username,
			"email":    resp.User.Email,
			"fullName": resp.User.FullName,
		},
		"expiresAt": resp.ExpiresAt,
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Refresh token not found",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.grpcClients.AuthClient.RefreshToken(ctx, &authv1.RefreshTokenRequest{
		RefreshToken: refreshToken,
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
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": resp.Message,
		})
		return
	}

	c.SetCookie(
		"refresh_token",
		resp.RefreshToken,
		int(resp.ExpiresIn)*24*7,
		"/",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"message":     resp.Message,
		"accessToken": resp.AccessToken,
		"expiresIn":   resp.ExpiresIn,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	refreshToken, _ := c.Cookie("refresh_token")

	var accessToken string
	if authHeader != "" {
		const bearerPrefix = "Bearer "
		if len(authHeader) >= len(bearerPrefix) && authHeader[:len(bearerPrefix)] == bearerPrefix {
			accessToken = authHeader[len(bearerPrefix):]
		}
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	resp, err := h.grpcClients.AuthClient.Logout(ctx, &authv1.LogoutRequest{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Authentication service unavailable",
			"error":   err.Error(),
		})
		return
	}

	c.SetCookie(
		"refresh_token",
		"",
		-1,
		"/",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": resp.Message,
	})
}
