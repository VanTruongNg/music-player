package handlers

import (
	"auth-service/internal/domain"
	"auth-service/internal/dto"
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

	resp := dto.UserLoginResponse{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		FullName:     user.FullName,
		CreatedAt:    user.CreatedAt.Format(time.RFC3339),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	utils.Success(c, http.StatusOK, resp)
}
