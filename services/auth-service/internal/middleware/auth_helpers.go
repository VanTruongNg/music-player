package middleware

import (
	"auth-service/internal/utils"
	customjwt "auth-service/internal/utils/jwt"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUserIDFromContext(c *gin.Context) (string, error) {
	userID, exists := c.Get(ContextKeyUserID)
	if !exists {
		return "", errors.New("user ID not found in context")
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return "", errors.New("user ID is not a valid string")
	}

	return userIDStr, nil
}

func GetUserJTIFromContext(c *gin.Context) (string, error) {
	userJTI, exists := c.Get(ContextKeyUserJTI)
	if !exists {
		return "", errors.New("user JTI not found in context")
	}

	userJTIStr, ok := userJTI.(string)
	if !ok {
		return "", errors.New("user JTI is not a valid string")
	}

	return userJTIStr, nil
}

func GetUserDataFromContext(c *gin.Context) (*customjwt.AccessClaims, error) {
	userData, exists := c.Get(ContextKeyUserData)
	if !exists {
		return nil, errors.New("user data not found in context")
	}

	claims, ok := userData.(*customjwt.AccessClaims)
	if !ok {
		return nil, errors.New("user data is not valid claims")
	}

	return claims, nil
}

func GetCurrentUser(c *gin.Context) (userID string, claims *customjwt.AccessClaims, err error) {
	userID, err = GetUserIDFromContext(c)
	if err != nil {
		return "", nil, err
	}

	claims, err = GetUserDataFromContext(c)
	if err != nil {
		return "", nil, err
	}

	return userID, claims, nil
}

func MustGetUserID(c *gin.Context) (string, bool) {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		utils.Fail(c, http.StatusUnauthorized, "UNAUTHORIZED", "User ID not found in context")
		return "", false
	}
	return userID, true
}

func MustGetCurrentUser(c *gin.Context) (string, *customjwt.AccessClaims, bool) {
	userID, claims, err := GetCurrentUser(c)
	if err != nil {
		utils.Fail(c, http.StatusUnauthorized, "UNAUTHORIZED", "User data not found in context")
		return "", nil, false
	}
	return userID, claims, true
}

func GetRefreshTokenFromCookie(c *gin.Context) (string, error) {
	refreshToken, err := c.Cookie("xs")
	if err != nil {
		return "", errors.New("refresh token cookie not found")
	}

	if refreshToken == "" {
		return "", errors.New("refresh token cookie is empty")
	}

	return refreshToken, nil
}

func MustGetRefreshTokenFromCookie(c *gin.Context) (string, bool) {
	refreshToken, err := GetRefreshTokenFromCookie(c)
	if err != nil {
		utils.Fail(c, http.StatusUnauthorized, "MISSING_REFRESH_TOKEN", "Refresh token cookie not found")
		return "", false
	}
	return refreshToken, true
}
