package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/oklog/ulid/v2"
)

type JWTService interface {
	SignAccessToken(userID string) (string, string, error)
	SignRefreshToken(userID string) (string, string, error)
	GetRefreshTTL() time.Duration // Get refresh token TTL from config
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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(j.cfg.AccessSecret))
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
	var secret string
	if isRefresh {
		secret = j.cfg.RefreshSecret
	} else {
		secret = j.cfg.AccessSecret
	}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
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

func (j *jwtService) GetRefreshTTL() time.Duration {
	return j.cfg.RefreshTTL
}
