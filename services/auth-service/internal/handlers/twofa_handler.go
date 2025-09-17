package handlers

import (
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
		utils.Fail(c, http.StatusBadRequest, "SETUP_2FA_FAILED", err.Error())
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
		utils.Fail(c, http.StatusBadRequest, "ENABLE_2FA_FAILED", err.Error())
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
		utils.Fail(c, http.StatusBadRequest, "VERIFY_2FA_FAILED", err.Error())
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"message": "2FA verified"})
}

func (h *TwoFAHandler) Disable2FA(c *gin.Context) {
	userID := c.Param("id")
	if err := h.service.Disable2FA(c.Request.Context(), userID); err != nil {
		utils.Fail(c, http.StatusBadRequest, "DISABLE_2FA_FAILED", err.Error())
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"message": "2FA disabled"})
}
