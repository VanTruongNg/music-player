package middleware

import (
	"gateway/internal/utils/jwt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	jwtVerifier *jwt.JWTVerifier
}

func NewAuthMiddleware(jwtVerifier *jwt.JWTVerifier) *AuthMiddleware {
	return &AuthMiddleware{
		jwtVerifier: jwtVerifier,
	}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
				"code":  "MISSING_AUTH_HEADER",
			})
			c.Abort()
			return
		}

		token, err := jwt.ExtractTokenFromHeader(authHeader)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format",
				"code":  "INVALID_AUTH_HEADER",
			})
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

			c.JSON(http.StatusUnauthorized, gin.H{
				"error": errorMessage,
				"code":  errorCode,
			})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_claims", claims)

		c.Next()
	})
}

func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		token, err := jwt.ExtractTokenFromHeader(authHeader)
		if err != nil {
			c.Next()
			return
		}

		claims, err := m.jwtVerifier.VerifyToken(token)
		if err != nil {
			c.Next()
			return
		}

		c.Set("user_id", claims.UserID)
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

func GetUserClaims(c *gin.Context) (*jwt.CustomClaims, bool) {
	claims, exists := c.Get("user_claims")
	if !exists {
		return nil, false
	}

	userClaims, ok := claims.(*jwt.CustomClaims)
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
