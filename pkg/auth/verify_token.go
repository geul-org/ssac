//ff:func feature=pkg-auth type=util control=sequence topic=auth-jwt
//ff:what HS256 JWT 토큰을 검증하고 MapClaims를 반환한다
package auth

import (
	"errors"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// VerifyToken parses req.Token as an HS256-signed JWT and returns the claim
// set as jwt.MapClaims. The signing secret is read from
// os.Getenv(Config.SecretEnv) on every call. Only HS256 is accepted;
// alg-confusion attacks are rejected via jwt.WithValidMethods.
func VerifyToken(req VerifyTokenRequest) (VerifyTokenResponse, error) {
	cfg := currentConfig()
	if cfg.SecretEnv == "" {
		return VerifyTokenResponse{}, errors.New("auth: SecretEnv not configured")
	}
	secret := os.Getenv(cfg.SecretEnv)
	if secret == "" {
		return VerifyTokenResponse{}, fmt.Errorf("auth: %s not set", cfg.SecretEnv)
	}

	parsed, err := jwt.Parse(req.Token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil {
		return VerifyTokenResponse{}, err
	}
	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok || !parsed.Valid {
		return VerifyTokenResponse{}, errors.New("auth: invalid token")
	}
	return VerifyTokenResponse{Claims: claims}, nil
}
