package middleware

import (
	"context"
	"gateway/internal/utils"
	"gateway/internal/utils/jwt"
	"net/http"
	"time"

	redisutil "gateway/internal/utils/redis"

	"github.com/gin-gonic/gin"
)

const (
	ContextKeyUserID   = "user_id"
	ContextKeyUserData = "user_claims"
	ContextUserSID     = "user_sid"
)

type AuthMiddleware struct {
	jwtVerifier jwt.JWTVerifier
	redisUtil   *redisutil.RedisUtil
}

func NewAuthMiddleware(jwtVerifier jwt.JWTVerifier, redisUtil *redisutil.RedisUtil) *AuthMiddleware {
	return &AuthMiddleware{
		jwtVerifier: jwtVerifier,
		redisUtil:   redisUtil,
	}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		token, err := m.jwtVerifier.ExtractTokenFromHeader(authHeader)
		if err != nil {
			utils.Fail(c, http.StatusUnauthorized, "MISSING_TOKEN", "Authorization token is required")
			c.Abort()
			return
		}

		claims, err := m.jwtVerifier.VerifyToken(token)
		if err != nil {
			var errorCode string
			var errorMessage string

			switch err {
			case jwt.ErrTokenExpired:
				errorCode = "TOKEN_EXPIRED"
				errorMessage = "Token has expired"
			case jwt.ErrTokenInvalid:
				errorCode = "TOKEN_INVALID"
				errorMessage = "Token is invalid"
			case jwt.ErrMissingKID:
				errorCode = "MISSING_KID"
				errorMessage = "Token missing key ID"
			case jwt.ErrUnexpectedSigningMethod:
				errorCode = "INVALID_SIGNING_METHOD"
				errorMessage = "Invalid token signing method"
			default:
				errorCode = "TOKEN_VERIFICATION_FAILED"
				errorMessage = "Token verification failed"
			}

			utils.Fail(c, http.StatusUnauthorized, errorCode, errorMessage)
			c.Abort()
			return
		}

		if claims == nil {
			utils.Fail(c, http.StatusUnauthorized, "INVALID_CLAIMS", "Token claims are invalid")
			c.Abort()
			return
		}
		if m.validateSession(c, claims) {
			utils.Fail(c, http.StatusUnauthorized, "SESSION_INVALID", "User session is invalid or has expired")
			c.Abort()
			return
		}

		m.setUserContext(c, claims)

		c.Next()
	})
}

func (m *AuthMiddleware) setUserContext(c *gin.Context, claims *jwt.AccessClaims) {
	c.Set(ContextKeyUserID, claims.Subject)
	c.Set(ContextKeyUserData, claims)
	c.Set(ContextUserSID, claims.SID)
}

func (m *AuthMiddleware) validateSession(c *gin.Context, claims *jwt.AccessClaims) bool {
	redisKey := "auth:session:" + claims.SID
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var sessionData map[string]interface{}
	err := m.redisUtil.GetJSON(ctx, redisKey, &sessionData)
	return err != nil
}

func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		token, err := m.jwtVerifier.ExtractTokenFromHeader(authHeader)
		if err != nil {
			c.Next()
			return
		}

		claims, err := m.jwtVerifier.VerifyToken(token)
		if err != nil {
			c.Next()
			return
		}

		c.Set("user_id", claims.Subject)
		c.Set("user_claims", claims)
		c.Next()
	})
}

func GetUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", false
	}

	userIDStr, ok := userID.(string)
	return userIDStr, ok
}

func GetUserClaims(c *gin.Context) (*jwt.AccessClaims, bool) {
	claims, exists := c.Get("user_claims")
	if !exists {
		return nil, false
	}

	userClaims, ok := claims.(*jwt.AccessClaims)
	return userClaims, ok
}

func CORSMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})
}
