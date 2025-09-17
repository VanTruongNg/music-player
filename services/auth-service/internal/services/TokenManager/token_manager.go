package tokenmanager

import (
	"auth-service/internal/utils/jwt"
	redisutil "auth-service/internal/utils/redis"
	"context"
	"time"
)

// ctxKey là type riêng cho context key để tránh collision
// Nên dùng chung type này ở mọi nơi inject/lấy value từ context
// See: https://golang.org/pkg/context/#WithValue
// và staticcheck SA1029
type CtxKey string

const (
	CtxKeyIP        CtxKey = "ip_address"
	CtxKeyUserAgent CtxKey = "user_agent"
)

// SessionInfo holds metadata for a refresh token session.
type SessionInfo struct {
	UserID    string    `json:"user_id"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
}

type TokenManager interface {
	GenerateTokens(ctx context.Context, userID string) (accessToken, refreshToken string, err error)
}

type tokenManager struct {
	jwtService jwt.JWTService
	redisUtil  *redisutil.RedisUtil
}

func NewTokenManager(jwtService jwt.JWTService, redisUtil *redisutil.RedisUtil) TokenManager {
	return &tokenManager{jwtService: jwtService, redisUtil: redisUtil}
}

// getStringFromContext safely gets a string value from context by key
func getStringFromContext(ctx context.Context, key CtxKey) string {
	val := ctx.Value(key)
	if s, ok := val.(string); ok {
		return s
	}
	return ""
}

// GenerateTokens generates access & refresh tokens for a user, and stores refresh session info in Redis.
func (tm *tokenManager) GenerateTokens(ctx context.Context, userID string) (string, string, error) {
	ip := getStringFromContext(ctx, CtxKeyIP)
	userAgent := getStringFromContext(ctx, CtxKeyUserAgent)

	accessToken, _, err := tm.jwtService.SignAccessToken(userID)
	if err != nil {
		return "", "", err
	}
	refreshToken, refreshJTI, err := tm.jwtService.SignRefreshToken(userID)
	if err != nil {
		return "", "", err
	}
	// Store session info in Redis for refresh token jti
	session := SessionInfo{
		UserID:    userID,
		IP:        ip,
		UserAgent: userAgent,
		CreatedAt: time.Now().UTC(),
	}
	key := "auth:session:" + refreshJTI
	refreshTTL := tm.jwtService.GetRefreshTTL()
	err = tm.redisUtil.SetJSON(ctx, key, session, refreshTTL)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}
