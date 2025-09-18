package jwt

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type JWTConfig struct {
	AccessPrivateKey ed25519.PrivateKey
	AccessKID        string
	AccessTTL        time.Duration
	RefreshSecret string
	RefreshTTL    time.Duration

	// JWKS for verification (test purposes)
	JWKSFile string
}

type JWKS struct {
	Keys []JWK `json:"keys"`
}

type JWK struct {
	Kty string `json:"kty"`
	Crv string `json:"crv"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	X   string `json:"x"`
}

func LoadJWTConfig() (*JWTConfig, error) {
	refreshSecret := viper.GetString("JWT_REFRESH_SECRET")
	if refreshSecret == "" {
		return nil, ErrInvalidJWTConfig
	}

	// Load TTL values
	accessTTL := time.Duration(viper.GetInt("JWT_ACCESS_TTL")) * time.Second
	refreshTTL := time.Duration(viper.GetInt("JWT_REFRESH_TTL")) * time.Second
	if accessTTL == 0 || refreshTTL == 0 {
		return nil, ErrInvalidJWTConfig
	}

	// Load access token EdDSA
	accessPrivateKeyFile := viper.GetString("JWT_ACCESS_PRIVATE_KEY_FILE")
	accessKID := viper.GetString("JWT_ACCESS_KID")
	jwksFile := viper.GetString("JWT_JWKS_FILE")

	if accessPrivateKeyFile == "" || accessKID == "" {
		return nil, fmt.Errorf("JWT_ACCESS_PRIVATE_KEY_FILE and JWT_ACCESS_KID must be set")
	}

	// Load EdDSA private key for access token signing
	accessPrivateKey, err := loadEd25519PrivateKey(accessPrivateKeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load EdDSA private key from %s: %w", accessPrivateKeyFile, err)
	}

	return &JWTConfig{
		AccessPrivateKey: accessPrivateKey,
		AccessKID:        accessKID,
		AccessTTL:        accessTTL,
		RefreshSecret:    refreshSecret,
		RefreshTTL:       refreshTTL,
		JWKSFile:         jwksFile,
	}, nil
}

// loadEd25519PrivateKey loads an Ed25519 private key from PEM file
func loadEd25519PrivateKey(filepath string) (ed25519.PrivateKey, error) {
	keyData, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	ed25519Key, ok := privateKey.(ed25519.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("key is not an Ed25519 private key")
	}

	return ed25519Key, nil
}

func (cfg *JWTConfig) GetPublicKeyFromJWKS(kid string) (ed25519.PublicKey, error) {
	if cfg.JWKSFile == "" {
		return nil, fmt.Errorf("JWKS file not configured")
	}

	jwksData, err := os.ReadFile(cfg.JWKSFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read JWKS file: %w", err)
	}

	var jwks JWKS
	if err := json.Unmarshal(jwksData, &jwks); err != nil {
		return nil, fmt.Errorf("failed to parse JWKS: %w", err)
	}

	for _, key := range jwks.Keys {
		if key.Kid == kid && key.Kty == "OKP" && key.Crv == "Ed25519" {
			pubKeyBytes, err := decodeBase64URL(key.X)
			if err != nil {
				return nil, fmt.Errorf("failed to decode public key: %w", err)
			}

			if len(pubKeyBytes) != ed25519.PublicKeySize {
				return nil, fmt.Errorf("invalid Ed25519 public key size")
			}

			return ed25519.PublicKey(pubKeyBytes), nil
		}
	}

	return nil, fmt.Errorf("key with KID %s not found in JWKS", kid)
}

// decodeBase64URL decodes base64url without padding
func decodeBase64URL(s string) ([]byte, error) {
	// Add padding if needed
	s = strings.TrimRight(s, "=")
	switch len(s) % 4 {
	case 2:
		s += "=="
	case 3:
		s += "="
	}

	// Replace URL-safe characters
	s = strings.Replace(s, "-", "+", -1)
	s = strings.Replace(s, "_", "/", -1)

	return base64.StdEncoding.DecodeString(s)
}
