package supabase

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/l-fraga2811/back-sable/internal/config"
)

type TokenValidator struct {
	jwks      *JwksCache
	jwtSecret []byte
}

func NewTokenValidator(cfg *config.Config) *TokenValidator {
	return &TokenValidator{
		jwks:      NewJwksCache(cfg.JwksURL),
		jwtSecret: []byte(cfg.JwtSecret),
	}
}

type AccessTokenClaims struct {
	jwt.RegisteredClaims
	Email        string                 `json:"email"`
	Role         string                 `json:"role"`
	UserMetadata map[string]interface{} `json:"user_metadata"`
}

func (c AccessTokenClaims) Username() string {
	if c.UserMetadata == nil {
		return ""
	}
	value, ok := c.UserMetadata["username"]
	if !ok {
		return ""
	}
	username, ok := value.(string)
	if !ok {
		return ""
	}
	return username
}

func (v *TokenValidator) Validate(tokenString string) (AccessTokenClaims, error) {
	alg, kid, err := tokenHeader(tokenString)
	if err != nil {
		return AccessTokenClaims{}, err
	}

	claims := AccessTokenClaims{}
	var parsed *jwt.Token

	if alg == "HS256" {
		if len(v.jwtSecret) == 0 {
			return AccessTokenClaims{}, errors.New("SUPABASE_JWT_SECRET not defined")
		}
		parsed, err = jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			return v.jwtSecret, nil
		}, jwt.WithValidMethods([]string{"HS256"}))
	} else {
		parsed, err = jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			key, err := v.jwks.GetPublicKey(kid)
			if err != nil {
				return nil, err
			}
			return key, nil
		}, jwt.WithValidMethods([]string{"RS256"}))
	}

	if err != nil {
		return AccessTokenClaims{}, err
	}
	if parsed == nil || !parsed.Valid {
		return AccessTokenClaims{}, errors.New("invalid token")
	}

	if claims.ExpiresAt != nil && time.Now().After(claims.ExpiresAt.Time) {
		return AccessTokenClaims{}, errors.New("token expired")
	}

	if claims.Subject == "" {
		return AccessTokenClaims{}, errors.New("sub missing")
	}

	return claims, nil
}

func tokenHeader(tokenString string) (string, string, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return "", "", errors.New("malformed token")
	}

	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return "", "", err
	}

	var header struct {
		Kid string `json:"kid"`
		Alg string `json:"alg"`
	}
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return "", "", err
	}
	if header.Alg == "" {
		return "", "", errors.New("alg missing")
	}
	if header.Alg != "HS256" && header.Kid == "" {
		return "", "", errors.New("kid missing")
	}

	return header.Alg, header.Kid, nil
}
