package jwt

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	SignAccessToken(userID, sid string, av uint64) (string, time.Time, error)
	SignRefreshToken(userID string, jti string) (string, time.Time, error)
	VerifyAccessToken(tokenStr string, isRefresh bool) (*AccessClaims, error)
	ExtractTokenFromHeader(authHeader string) (string, error)
	GetAccessTTL() time.Duration
	GetRefreshTTL() time.Duration
	GetJWKS() (*JWKS, error)
}

type jwtService struct {
	cfg *JWTConfig
}

func NewJWTService(cfg *JWTConfig) JWTService {
	return &jwtService{
		cfg: cfg,
	}
}

func (j *jwtService) SignAccessToken(userID, sid string, av uint64) (string, time.Time, error) {
	now := time.Now().UTC()
	exp := now.Add(j.cfg.AccessTTL)

	claims := &AccessClaims{
		SID: sid,
		AV:  av,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	token.Header["kid"] = j.cfg.AccessKID

	signed, err := token.SignedString(j.cfg.AccessPrivateKey)
	return signed, exp, err
}

func (j *jwtService) SignRefreshToken(userID string, jti string) (string, time.Time, error) {
	now := time.Now().UTC()
	exp := now.Add(j.cfg.RefreshTTL)

	claims := &RefreshClaims{
		UserID: userID,
		JTI:    jti,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(j.cfg.RefreshSecret))
	return signed, exp, err
}

func (j *jwtService) VerifyAccessToken(tokenStr string, isRefresh bool) (*AccessClaims, error) {
	claims := &AccessClaims{}

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

func (j *jwtService) GetAccessTTL() time.Duration {
	return j.cfg.AccessTTL
}

func (j *jwtService) GetRefreshTTL() time.Duration {
	return j.cfg.RefreshTTL
}

func (j *jwtService) ExtractTokenFromHeader(authHeader string) (string, error) {
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

func (j *jwtService) GetJWKS() (*JWKS, error) {
	if j.cfg.JWKSFile == "" {
		return nil, fmt.Errorf("JWKS file not configured")
	}

	jwksData, err := os.ReadFile(j.cfg.JWKSFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read JWKS file: %w", err)
	}

	var jwks JWKS
	if err := json.Unmarshal(jwksData, &jwks); err != nil {
		return nil, fmt.Errorf("failed to parse JWKS: %w", err)
	}

	return &jwks, nil
}
