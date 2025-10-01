package jwt

import "github.com/golang-jwt/jwt/v5"

type AccessClaims struct {
	SID string `json:"sid"`
	SV  uint64 `json:"sv"`
	jwt.RegisteredClaims
}