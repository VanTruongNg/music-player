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
	UserID      string    `json:"user_id"`
	Status      string    `json:"status"`
	AV          uint64    `json:"av"`
	IP          string    `json:"ip"`
	UserAgent   string    `json:"user_agent"`
	CreatedAt   time.Time `json:"created_at"`
	RTCurrent   string    `json:"rt_current"`
	RTPrev      string    `json:"rt_prev"`
	RTRotatedAt time.Time `json:"rt_rotated_at"`
}

type TokenManager interface {
	IssueInitialTokens(ctx context.Context, userID string) (string, string, error)
	RefreshToken(ctx context.Context, claims *jwt.RefreshClaims) (string, string, error)
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
	jti := ulid.Make().String()
	const avInit uint64 = 1

	accessToken, _, err := tm.jwtService.SignAccessToken(userID, sid, avInit)
	if err != nil {
		return "", "", err
	}

	refreshToken, _, err := tm.jwtService.SignRefreshToken(userID, sid, jti)
	if err != nil {
		return "", "", err
	}

	session := SessionInfo{
		UserID:      userID,
		AV:          avInit,
		IP:          ip,
		UserAgent:   userAgent,
		CreatedAt:   time.Now().UTC(),
		Status:      "active",
		RTCurrent:   jti,
		RTPrev:      "",
		RTRotatedAt: time.Now().UTC(),
	}

	key := "auth:session:" + sid
	refreshTTL := tm.jwtService.GetRefreshTTL()
	err = tm.redisUtil.SetJSON(ctx, key, session, refreshTTL)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func (tm *tokenManager) RefreshToken(ctx context.Context, claims *jwt.RefreshClaims) (string, string, error) {
	ip := getStringFromContext(ctx, CtxKeyIP)
	userAgent := getStringFromContext(ctx, CtxKeyUserAgent)

	key := "auth:session:" + claims.SID
	var sess SessionInfo
	if err := tm.redisUtil.GetJSON(ctx, key, &sess); err != nil {
		return "", "", jwt.ErrSessionNotFound
	}

	if sess.Status != "active" {
		return "", "", jwt.ErrSessionRevoked
	}

	newJTI := ulid.Make().String()
	sess.AV++

	accessToken, _, err := tm.jwtService.SignAccessToken(claims.UserID, claims.SID, sess.AV)
	if err != nil {
		return "", "", err
	}

	refreshToken, _, err := tm.jwtService.SignRefreshToken(claims.UserID, claims.SID, newJTI)
	if err != nil {
		return "", "", err
	}

	prev := sess.RTCurrent
	sess.RTPrev = prev
	sess.RTCurrent = newJTI
	sess.RTRotatedAt = time.Now().UTC()
	sess.IP = ip
	sess.UserAgent = userAgent

	if err := tm.redisUtil.SetJSON(ctx, key, sess, tm.jwtService.GetRefreshTTL()); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
