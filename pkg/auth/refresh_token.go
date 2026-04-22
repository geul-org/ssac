//ff:func feature=pkg-auth type=util control=sequence topic=auth-jwt
//ff:what HS256 서명된 JWT 리프레시 토큰을 발급한다
package auth

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// RefreshToken signs a new HS256 refresh token using the claims supplied in
// req.Claims. The `exp` claim is set to now + RefreshTTL (as configured via
// Configure). The signing secret is read from os.Getenv(Config.SecretEnv)
// on every call.
func RefreshToken(req RefreshTokenRequest) (RefreshTokenResponse, error) {
	cfg := currentConfig()
	if cfg.SecretEnv == "" {
		return RefreshTokenResponse{}, errors.New("auth: SecretEnv not configured")
	}
	secret := os.Getenv(cfg.SecretEnv)
	if secret == "" {
		return RefreshTokenResponse{}, fmt.Errorf("auth: %s not set", cfg.SecretEnv)
	}
	if cfg.RefreshTTL <= 0 {
		return RefreshTokenResponse{}, errors.New("auth: RefreshTTL not configured")
	}

	claims, err := claimsToMap(req.Claims)
	if err != nil {
		return RefreshTokenResponse{}, fmt.Errorf("auth: marshal claims: %w", err)
	}
	exp := time.Now().Add(cfg.RefreshTTL)
	claims["exp"] = exp.Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return RefreshTokenResponse{}, err
	}
	return RefreshTokenResponse{RefreshToken: signed, ExpiresAt: exp}, nil
}
