package jwt

import (
	"crypto/ed25519"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// JWTVerifier handles JWT token verification using JWKS
type JWTVerifier struct {
	jwksClient *JWKSClient
}

// NewJWTVerifier creates a new JWT verifier
func NewJWTVerifier(jwksClient *JWKSClient) *JWTVerifier {
	return &JWTVerifier{
		jwksClient: jwksClient,
	}
}

func (v *JWTVerifier) VerifyToken(tokenStr string) (*CustomClaims, error) {
	claims := &CustomClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, ErrUnexpectedSigningMethod
		}

		kidInterface, exists := token.Header["kid"]
		if !exists {
			return nil, ErrMissingKID
		}

		kid, ok := kidInterface.(string)
		if !ok {
			return nil, ErrInvalidKID
		}

		publicKey, err := v.getPublicKeyFromJWKS(kid)
		if err != nil {
			return nil, fmt.Errorf("failed to get public key: %w", err)
		}

		return publicKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrTokenInvalid
	}

	if !token.Valid {
		return nil, ErrTokenInvalid
	}

	return claims, nil
}

func (v *JWTVerifier) getPublicKeyFromJWKS(kid string) (ed25519.PublicKey, error) {
	jwk, err := v.jwksClient.GetKeyByKID(kid)
	if err != nil {
		return nil, err
	}

	if jwk.Kty != "OKP" || jwk.Crv != "Ed25519" {
		return nil, fmt.Errorf("%w: type=%s, curve=%s", ErrInvalidKeyFormat, jwk.Kty, jwk.Crv)
	}

	pubKeyBytes, err := decodeBase64URL(jwk.X)
	if err != nil {
		return nil, fmt.Errorf("failed to decode public key: %w", err)
	}

	if len(pubKeyBytes) != ed25519.PublicKeySize {
		return nil, fmt.Errorf("%w: invalid Ed25519 public key size: %d", ErrInvalidKeyFormat, len(pubKeyBytes))
	}

	return ed25519.PublicKey(pubKeyBytes), nil
}

func ExtractTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", ErrTokenInvalid
	}

	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return "", ErrTokenInvalid
	}

	token := strings.TrimPrefix(authHeader, bearerPrefix)
	if token == "" {
		return "", ErrTokenInvalid
	}

	return token, nil
}

func decodeBase64URL(s string) ([]byte, error) {
	s = strings.TrimRight(s, "=")
	switch len(s) % 4 {
	case 2:
		s += "=="
	case 3:
		s += "="
	}

	s = strings.Replace(s, "-", "+", -1)
	s = strings.Replace(s, "_", "/", -1)

	return base64.StdEncoding.DecodeString(s)
}
