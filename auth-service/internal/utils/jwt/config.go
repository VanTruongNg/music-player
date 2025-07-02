package jwt

import (
	"time"

	"github.com/spf13/viper"
)

type JWTConfig struct {
	AccessSecret  string
	RefreshSecret string
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
}

func LoadJWTConfig() (*JWTConfig, error) {
	accessSecret := viper.GetString("JWT_ACCESS_SECRET")
	refreshSecret := viper.GetString("JWT_REFRESH_SECRET")
	if accessSecret == "" || refreshSecret == "" {
		return nil, ErrInvalidJWTConfig
	}

	accessTTL := time.Duration(viper.GetInt("JWT_ACCESS_TTL")) * time.Second
	refreshTTL := time.Duration(viper.GetInt("JWT_REFRESH_TTL")) * time.Second
	if accessTTL == 0 || refreshTTL == 0 {
		return nil, ErrInvalidJWTConfig
	}

	return &JWTConfig{
		AccessSecret:  accessSecret,
		RefreshSecret: refreshSecret,
		AccessTTL:     accessTTL,
		RefreshTTL:    refreshTTL,
	}, nil
}
