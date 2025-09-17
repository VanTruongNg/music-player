package handlers

import (
	"auth-service/internal/domain"
	"auth-service/internal/dto"
	"auth-service/internal/middleware"
	"auth-service/internal/services"
	tokenmanager "auth-service/internal/services/TokenManager"
	"auth-service/internal/utils"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.service.GetAllUsers(c.Request.Context())
	if err != nil {
		if derr, ok := err.(*domain.DomainError); ok {
			utils.Fail(c, derr.Status, derr.Code, derr.Message)
			return
		}
		utils.Fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	var resp []dto.UserResponse
	for _, user := range users {
		resp = append(resp, dto.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			FullName:  user.FullName,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		})
	}
	utils.Success(c, http.StatusOK, resp)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	userID := c.Param("id")

	user, err := h.service.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		if derr, ok := err.(*domain.DomainError); ok {
			utils.Fail(c, derr.Status, derr.Code, derr.Message)
			return
		}
		utils.Fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	resp := dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FullName:  user.FullName,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}
	utils.Success(c, http.StatusOK, resp)
}

func (h *UserHandler) GetMe(c *gin.Context) {
	userID, ok := middleware.MustGetUserID(c)
	if !ok {
		return
	}

	user, err := h.service.GetMe(c.Request.Context(), userID)
	if err != nil {
		if derr, ok := err.(*domain.DomainError); ok {
			utils.Fail(c, derr.Status, derr.Code, derr.Message)
			return
		}
		utils.Fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	resp := dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FullName:  user.FullName,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}
	utils.Success(c, http.StatusOK, resp)
}

func (h *UserHandler) Register(c *gin.Context) {
	var req dto.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	createdUser, err := h.service.Register(c.Request.Context(), &req)
	if err != nil {
		if derr, ok := err.(*domain.DomainError); ok {
			utils.Fail(c, derr.Status, derr.Code, derr.Message)
			return
		}
		utils.Fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	resp := dto.UserRegisterResponse{
		ID:        createdUser.ID,
		Username:  createdUser.Username,
		Email:     createdUser.Email,
		FullName:  createdUser.FullName,
		CreatedAt: createdUser.CreatedAt.Format(time.RFC3339),
	}
	utils.Success(c, http.StatusCreated, resp)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req dto.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	ctx := context.WithValue(c.Request.Context(), tokenmanager.CtxKeyIP, c.ClientIP())
	ctx = context.WithValue(ctx, tokenmanager.CtxKeyUserAgent, c.Request.UserAgent())

	user, accessToken, refreshToken, err := h.service.Login(ctx, &req)
	if err != nil {
		if derr, ok := err.(*domain.DomainError); ok {
			utils.Fail(c, derr.Status, derr.Code, derr.Message)
			return
		}
		utils.Fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	c.SetCookie(
		"xs",
		refreshToken,
		60*60*24*7, // maxAge (7 days)
		"/",        // path - send to all endpoints
		"",         // domain - current domain
		false,      // secure - false for localhost HTTP
		true,       // httpOnly - true for security
	)

	resp := dto.UserLoginResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		FullName:    user.FullName,
		CreatedAt:   user.CreatedAt.Format(time.RFC3339),
		AccessToken: accessToken,
	}
	utils.Success(c, http.StatusOK, resp)
}
