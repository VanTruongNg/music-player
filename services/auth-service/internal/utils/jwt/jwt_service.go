package jwt

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/oklog/ulid/v2"
)

type JWTService interface {
	SignAccessToken(userID string) (string, string, error)
	SignRefreshToken(userID string) (string, string, error)
	VerifyToken(tokenStr string, isRefresh bool) (*CustomClaims, error)
	ExtractTokenFromHeader(authHeader string) (string, error)
	GetRefreshTTL() time.Duration
	GetAccessKID() string
}

type jwtService struct {
	cfg *JWTConfig
}

func NewJWTService(cfg *JWTConfig) *jwtService {
	return &jwtService{
		cfg: cfg,
	}
}

func (j *jwtService) SignAccessToken(userID string) (string, string, error) {
	jti := ulid.Make().String()
	now := time.Now().UTC()
	exp := now.Add(j.cfg.AccessTTL)

	claims := &CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
			ID:        jti,
			Subject:   userID,
		},
	}

	// Use EdDSA (Ed25519) for access token
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)

	// Set the key ID in the header
	token.Header["kid"] = j.cfg.AccessKID

	signed, err := token.SignedString(j.cfg.AccessPrivateKey)
	return signed, jti, err
}

func (j *jwtService) SignRefreshToken(userID string) (string, string, error) {
	jti := ulid.Make().String()
	now := time.Now().UTC()
	exp := now.Add(j.cfg.RefreshTTL)

	claims := &CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
			ID:        jti,
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(j.cfg.RefreshSecret))
	return signed, jti, err
}

func (j *jwtService) VerifyToken(tokenStr string, isRefresh bool) (*CustomClaims, error) {
	claims := &CustomClaims{}

	var keyFunc jwt.Keyfunc
	if isRefresh {
		keyFunc = func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrUnexpectedSigningMethod
			}
			return []byte(j.cfg.RefreshSecret), nil
		}
	} else {
		keyFunc = func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
				return nil, ErrUnexpectedSigningMethod
			}

			kidInterface, exists := token.Header["kid"]
			if !exists {
				return nil, ErrTokenInvalid
			}

			kid, ok := kidInterface.(string)
			if !ok {
				return nil, ErrTokenInvalid
			}

			publicKey, err := j.cfg.GetPublicKeyFromJWKS(kid)
			if err != nil {
				return nil, ErrTokenInvalid
			}

			return publicKey, nil
		}
	}

	token, err := jwt.ParseWithClaims(tokenStr, claims, keyFunc)

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

func (j *jwtService) GetRefreshTTL() time.Duration {
	return j.cfg.RefreshTTL
}

func (j *jwtService) GetAccessKID() string {
	return j.cfg.AccessKID
}

func (j *jwtService) ExtractTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", ErrTokenInvalid
	}

	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return "", ErrTokenInvalid
	}

	// Extract token part
	token := strings.TrimPrefix(authHeader, bearerPrefix)
	if token == "" {
		return "", ErrTokenInvalid
	}

	return token, nil
}
