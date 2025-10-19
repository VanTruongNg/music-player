package jwt

import "errors"

// Predefined errors for JWT operations
var (
	ErrTokenExpired            = errors.New("jwt: token expired")
	ErrTokenInvalid            = errors.New("jwt: invalid token")
	ErrUnexpectedSigningMethod = errors.New("jwt: unexpected signing method")
	ErrInvalidJWTConfig        = errors.New("jwt: invalid JWT config in environment")
	ErrSessionNotFound         = errors.New("jwt: session not found")
	ErrSessionRevoked          = errors.New("jwt: session revoked")
)
