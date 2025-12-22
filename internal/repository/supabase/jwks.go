package supabase

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"sync"
	"time"
)

type jwksResponse struct {
	Keys []jwkKey `json:"keys"`
}

type jwkKey struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type JwksCache struct {
	jwksURL string
	client  *http.Client

	mu        sync.RWMutex
	publicKey map[string]*rsa.PublicKey
	expiresAt time.Time
	ttl       time.Duration
}

func NewJwksCache(jwksURL string) *JwksCache {
	return &JwksCache{
		jwksURL:   jwksURL,
		client:    &http.Client{Timeout: 10 * time.Second},
		publicKey: map[string]*rsa.PublicKey{},
		ttl:       10 * time.Minute,
	}
}

func (c *JwksCache) GetPublicKey(kid string) (*rsa.PublicKey, error) {
	if kid == "" {
		return nil, errors.New("kid not provided")
	}

	c.mu.RLock()
	if time.Now().Before(c.expiresAt) {
		if key, ok := c.publicKey[kid]; ok {
			c.mu.RUnlock()
			return key, nil
		}
	}
	c.mu.RUnlock()

	if err := c.refresh(); err != nil {
		return nil, err
	}

	c.mu.RLock()
	defer c.mu.RUnlock()
	key, ok := c.publicKey[kid]
	if !ok {
		return nil, errors.New("kid not found in jwks")
	}
	return key, nil
}

func (c *JwksCache) refresh() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if time.Now().Before(c.expiresAt) {
		return nil
	}

	resp, err := c.client.Get(c.jwksURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("failed to fetch jwks")
	}

	var parsed jwksResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return err
	}

	newKeys := make(map[string]*rsa.PublicKey, len(parsed.Keys))
	for _, key := range parsed.Keys {
		if key.Kty != "RSA" || key.N == "" || key.E == "" || key.Kid == "" {
			continue
		}

		pub, err := rsaFromJwk(key.N, key.E)
		if err != nil {
			continue
		}
		newKeys[key.Kid] = pub
	}

	if len(newKeys) == 0 {
		return errors.New("jwks has no valid keys")
	}

	c.publicKey = newKeys
	c.expiresAt = time.Now().Add(c.ttl)
	return nil
}

func rsaFromJwk(nB64 string, eB64 string) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(nB64)
	if err != nil {
		return nil, err
	}
	eBytes, err := base64.RawURLEncoding.DecodeString(eB64)
	if err != nil {
		return nil, err
	}

	n := new(big.Int).SetBytes(nBytes)
	e := new(big.Int).SetBytes(eBytes)
	if !e.IsInt64() {
		return nil, errors.New("invalid exponent")
	}

	return &rsa.PublicKey{N: n, E: int(e.Int64())}, nil
}
