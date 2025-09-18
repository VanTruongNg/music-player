package jwt

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// JWKS represents JSON Web Key Set structure
type JWKS struct {
	Keys []JWK `json:"keys"`
}

// JWK represents a single JSON Web Key
type JWK struct {
	Kty string `json:"kty"` // Key Type
	Crv string `json:"crv"` // Curve (for EdDSA)
	Use string `json:"use"` // Usage
	Kid string `json:"kid"` // Key ID
	X   string `json:"x"`   // X coordinate (base64url encoded)
}

// JWKSClient handles fetching and caching JWKS from auth service
type JWKSClient struct {
	authServiceURL string
	httpClient     *http.Client
	cache          *jwksCache
}

type jwksCache struct {
	data      *JWKS
	lastFetch time.Time
	ttl       time.Duration
	mutex     sync.RWMutex
}

// NewJWKSClient creates a new JWKS client
func NewJWKSClient(authServiceURL string) *JWKSClient {
	return &JWKSClient{
		authServiceURL: authServiceURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		cache: &jwksCache{
			ttl: 1 * time.Hour, // Cache for 1 hour
		},
	}
}

// GetJWKS fetches JWKS from auth service with caching
func (c *JWKSClient) GetJWKS() (*JWKS, error) {
	// Check cache first
	c.cache.mutex.RLock()
	if c.cache.data != nil && time.Since(c.cache.lastFetch) < c.cache.ttl {
		cached := c.cache.data
		c.cache.mutex.RUnlock()
		return cached, nil
	}
	c.cache.mutex.RUnlock()

	// Cache miss or expired, fetch from auth service
	jwks, err := c.fetchJWKS()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrJWKSFetchFailed, err)
	}

	// Update cache
	c.cache.mutex.Lock()
	c.cache.data = jwks
	c.cache.lastFetch = time.Now()
	c.cache.mutex.Unlock()

	return jwks, nil
}

// fetchJWKS performs the actual HTTP request to get JWKS
func (c *JWKSClient) fetchJWKS() (*JWKS, error) {
	url := fmt.Sprintf("%s/.well-known/jwks.json", c.authServiceURL)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var jwks JWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, fmt.Errorf("failed to decode JWKS response: %w", err)
	}

	return &jwks, nil
}

// GetKeyByKID finds a specific key by its Key ID
func (c *JWKSClient) GetKeyByKID(kid string) (*JWK, error) {
	jwks, err := c.GetJWKS()
	if err != nil {
		return nil, err
	}

	for _, key := range jwks.Keys {
		if key.Kid == kid {
			return &key, nil
		}
	}

	return nil, fmt.Errorf("%w: kid=%s", ErrKeyNotFound, kid)
}

// InvalidateCache forces a fresh fetch on next request
func (c *JWKSClient) InvalidateCache() {
	c.cache.mutex.Lock()
	c.cache.data = nil
	c.cache.mutex.Unlock()
}
