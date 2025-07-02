package middleware

import (
	"auth-service/internal/utils"
	customjwt "auth-service/internal/utils/jwt"
	redisutil "auth-service/internal/utils/redis"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	ContextKeyUserID   = "user_id"
	ContextKeyUserJTI  = "user_jti"
	ContextKeyUserData = "user_data"
)

type AuthMiddleware struct {
	jwtService customjwt.JWTService
	redisUtil  *redisutil.RedisUtil
}

func NewAuthMiddleware(jwtService customjwt.JWTService, redisUtil *redisutil.RedisUtil) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
		redisUtil:  redisUtil,
	}
}

func (mw *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		token, err := mw.jwtService.ExtractTokenFromHeader(authHeader)
		if err != nil {
			utils.Fail(c, http.StatusUnauthorized, "MISSING_TOKEN", "Authorization token is required")
			c.Abort()
			return
		}

		claims, err := mw.jwtService.VerifyToken(token, false)
		if err != nil {
			var code, message string
			switch err {
			case customjwt.ErrTokenExpired:
				code = "TOKEN_EXPIRED"
				message = "Access token has expired"
			case customjwt.ErrTokenInvalid:
				code = "TOKEN_INVALID"
				message = "Invalid access token"
			default:
				code = "TOKEN_ERROR"
				message = "Token verification failed"
			}
			utils.Fail(c, http.StatusUnauthorized, code, message)
			c.Abort()
			return
		}
		if claims == nil {
			utils.Fail(c, http.StatusUnauthorized, "TOKEN_INVALID", "Invalid access token")
			c.Abort()
			return
		}

		if !mw.validateSession(c) {
			c.Abort()
			return
		}

		mw.setUserContext(c, claims)
		c.Next()
	}
}

func (mw *AuthMiddleware) setUserContext(c *gin.Context, claims *customjwt.CustomClaims) {
	c.Set(ContextKeyUserID, claims.UserID)
	c.Set(ContextKeyUserJTI, claims.ID)
	c.Set(ContextKeyUserData, claims)
}

func (mw *AuthMiddleware) validateSession(c *gin.Context) bool {
	refreshToken, err := GetRefreshTokenFromCookie(c)
	if err != nil {
		utils.Fail(c, http.StatusUnauthorized, "SESSION_INVALID", "Session expired - please login again")
		return false
	}

	refreshClaims, err := mw.jwtService.VerifyToken(refreshToken, true)
	if err != nil {
		utils.Fail(c, http.StatusUnauthorized, "SESSION_INVALID", "Invalid session - please login again")
		return false
	}

	redisKey := "auth:session:" + refreshClaims.ID
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// Try to get session data from Redis as JSON
	var sessionData map[string]interface{}
	err = mw.redisUtil.GetJSON(ctx, redisKey, &sessionData)
	if err != nil {
		utils.Fail(c, http.StatusUnauthorized, "SESSION_REVOKED", "Session has been revoked - please login again")
		return false
	}

	return true
}
