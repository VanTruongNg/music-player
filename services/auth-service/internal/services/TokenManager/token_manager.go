package tokenmanager

import (
	"auth-service/internal/utils/jwt"
	redisutil "auth-service/internal/utils/redis"
	"context"
	"time"

	"github.com/oklog/ulid/v2"
)

type CtxKey string

const (
	CtxKeyIP        CtxKey = "ip_address"
	CtxKeyUserAgent CtxKey = "user_agent"
)

type SessionInfo struct {
	UserID    string    `json:"user_id"`
	SV        uint64    `json:"sv"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
}

type TokenManager interface {
	IssueInitialTokens(ctx context.Context, userID string) (string, string, error)
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

func (tm *tokenManager) IssueInitialTokens(ctx context.Context, userID string) (string, string, error) {
	ip := getStringFromContext(ctx, CtxKeyIP)
	userAgent := getStringFromContext(ctx, CtxKeyUserAgent)

	sid := ulid.Make().String()
	const svInit uint64 = 1

	accessToken, _, err := tm.jwtService.SignAccessToken(userID, sid, svInit)
	if err != nil {
		return "", "", err
	}

	refreshToken, _, err := tm.jwtService.SignRefreshToken(userID)
	if err != nil {
		return "", "", err
	}

	session := SessionInfo{
		UserID:    userID,
		SV:        svInit,
		IP:        ip,
		UserAgent: userAgent,
		CreatedAt: time.Now().UTC(),
	}

	key := "auth:session:" + sid
	accessTTL := tm.jwtService.GetAccessTTL()
	err = tm.redisUtil.SetJSON(ctx, key, session, accessTTL)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}
