package jwt

import "errors"

// Predefined errors for JWT operations in gateway
var (
	ErrTokenExpired            = errors.New("jwt: token expired")
	ErrTokenInvalid            = errors.New("jwt: invalid token")
	ErrUnexpectedSigningMethod = errors.New("jwt: unexpected signing method")
	ErrMissingKID              = errors.New("jwt: missing kid in token header")
	ErrInvalidKID              = errors.New("jwt: invalid kid in token header")
	ErrJWKSFetchFailed         = errors.New("jwt: failed to fetch JWKS")
	ErrKeyNotFound             = errors.New("jwt: key not found in JWKS")
	ErrInvalidKeyFormat        = errors.New("jwt: invalid key format")
)
