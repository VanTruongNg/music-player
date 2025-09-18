package jwt

import "github.com/golang-jwt/jwt/v5"

// CustomClaims represents the JWT claims structure
// This should match the claims structure from auth-service
type CustomClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}
