package jwt

import "github.com/golang-jwt/jwt/v5"

type AccessClaims struct {
	SID string `json:"sid"`
	AV  uint64 `json:"av"`
	jwt.RegisteredClaims
}