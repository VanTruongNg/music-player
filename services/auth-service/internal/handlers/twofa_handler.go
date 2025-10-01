package handlers

import (
	"auth-service/internal/domain"
	"auth-service/internal/services"
	"auth-service/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TwoFAHandler struct {
	service services.TwoFAService
}

func NewTwoFAHandler(service services.TwoFAService) *TwoFAHandler {
	return &TwoFAHandler{service: service}
}

func (h *TwoFAHandler) Setup2FA(c *gin.Context) {
	userID := c.Param("id")
	result, err := h.service.Setup2FA(c.Request.Context(), userID)
	if err != nil {
		if domainErr, ok := err.(*domain.DomainError); ok {
			utils.Fail(c, domainErr.Status, domainErr.Code, domainErr.Message)
		} else {
			utils.Fail(c, http.StatusInternalServerError, "SETUP_2FA_FAILED", "Internal server error")
		}
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"secret": result.Secret, "otp_url": result.OTPURL})
}

func (h *TwoFAHandler) Enable2FA(c *gin.Context) {
	userID := c.Param("id")
	var req struct {
		Code string `json:"code"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Code == "" {
		utils.Fail(c, http.StatusBadRequest, "INVALID_CODE", "Invalid code")
		return
	}
	if err := h.service.Enable2FA(c.Request.Context(), userID, req.Code); err != nil {
		if derr, ok := err.(*domain.DomainError); ok {
			utils.Fail(c, derr.Status, derr.Code, derr.Message)
			return
		}
		utils.Fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"message": "2FA enabled"})
}

func (h *TwoFAHandler) Verify2FA(c *gin.Context) {
	userID := c.Param("id")
	var req struct {
		Code string `json:"code"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Code == "" {
		utils.Fail(c, http.StatusBadRequest, "INVALID_CODE", "Invalid code")
		return
	}
	if err := h.service.Verify2FA(c.Request.Context(), userID, req.Code); err != nil {
		if domainErr, ok := err.(*domain.DomainError); ok {
			utils.Fail(c, domainErr.Status, domainErr.Code, domainErr.Message)
		} else {
			utils.Fail(c, http.StatusInternalServerError, "VERIFY_2FA_FAILED", "Internal server error")
		}
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"message": "2FA verified"})
}

func (h *TwoFAHandler) Disable2FA(c *gin.Context) {
	userID := c.Param("id")
	if err := h.service.Disable2FA(c.Request.Context(), userID); err != nil {
		if domainErr, ok := err.(*domain.DomainError); ok {
			utils.Fail(c, domainErr.Status, domainErr.Code, domainErr.Message)
		} else {
			utils.Fail(c, http.StatusInternalServerError, "DISABLE_2FA_FAILED", "Internal server error")
		}
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"message": "2FA disabled"})
}
