package handlers

import (
	"auth-service/internal/services"
	"auth-service/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TwoFAHandler handles HTTP requests for 2FA features.
type TwoFAHandler struct {
	service *services.TwoFAService
}

// NewTwoFAHandler creates a new TwoFAHandler.
func NewTwoFAHandler(service *services.TwoFAService) *TwoFAHandler {
	return &TwoFAHandler{service: service}
}

// Setup2FA generates a new 2FA secret and OTP URL for the user.
func (h *TwoFAHandler) Setup2FA(c *gin.Context) {
	userID := c.Param("id")
	result, err := h.service.Setup2FA(c.Request.Context(), userID)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "SETUP_2FA_FAILED", err.Error())
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"secret": result.Secret, "otp_url": result.OTPURL})
}

// Enable2FA enables 2FA for the user after verifying the provided code.
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

// Verify2FA verifies the 2FA code for the user.
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

// Disable2FA disables 2FA for the user.
func (h *TwoFAHandler) Disable2FA(c *gin.Context) {
	userID := c.Param("id")
	if err := h.service.Disable2FA(c.Request.Context(), userID); err != nil {
		utils.Fail(c, http.StatusBadRequest, "DISABLE_2FA_FAILED", err.Error())
		return
	}
	utils.Success(c, http.StatusOK, gin.H{"message": "2FA disabled"})
}
