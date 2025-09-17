package tokenmanager

import (
	"auth-service/internal/utils/jwt"
	redisutil "auth-service/internal/utils/redis"
	"context"
	"time"
)

type CtxKey string

const (
	CtxKeyIP        CtxKey = "ip_address"
	CtxKeyUserAgent CtxKey = "user_agent"
)

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

func getStringFromContext(ctx context.Context, key CtxKey) string {
	val := ctx.Value(key)
	if s, ok := val.(string); ok {
		return s
	}
	return ""
}

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
