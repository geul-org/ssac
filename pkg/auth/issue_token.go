//ff:func feature=pkg-auth type=util control=sequence topic=auth-jwt
//ff:what HS256 서명된 JWT 액세스 토큰을 발급한다
package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// IssueToken signs a new HS256 access token using the claims supplied in
// req.Claims. The `exp` claim is set to now + AccessTTL (as configured via
// Configure). The signing secret is read from os.Getenv(Config.SecretEnv) on
// every call so secrets may be rotated without re-configuring; an empty
// secret returns an error.
func IssueToken(req IssueTokenRequest) (IssueTokenResponse, error) {
	cfg := currentConfig()
	if cfg.SecretEnv == "" {
		return IssueTokenResponse{}, errors.New("auth: SecretEnv not configured")
	}
	secret := os.Getenv(cfg.SecretEnv)
	if secret == "" {
		return IssueTokenResponse{}, fmt.Errorf("auth: %s not set", cfg.SecretEnv)
	}
	if cfg.AccessTTL <= 0 {
		return IssueTokenResponse{}, errors.New("auth: AccessTTL not configured")
	}

	claims, err := claimsToMap(req.Claims)
	if err != nil {
		return IssueTokenResponse{}, fmt.Errorf("auth: marshal claims: %w", err)
	}
	exp := time.Now().Add(cfg.AccessTTL)
	claims["exp"] = exp.Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return IssueTokenResponse{}, err
	}
	return IssueTokenResponse{AccessToken: signed, ExpiresAt: exp}, nil
}

// claimsToMap converts the passthrough Claims (any) to jwt.MapClaims via
// JSON round-trip. A nil Claims produces an empty map. Shared by IssueToken
// and RefreshToken.
func claimsToMap(in any) (jwt.MapClaims, error) {
	if in == nil {
		return jwt.MapClaims{}, nil
	}
	raw, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}
	out := jwt.MapClaims{}
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, err
	}
	return out, nil
}
