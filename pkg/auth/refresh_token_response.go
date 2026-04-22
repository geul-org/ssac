//ff:type feature=pkg-auth type=model topic=auth-jwt
//ff:what JWT 리프레시 토큰 발급 응답 모델
package auth

import "time"

// RefreshTokenResponse is the result of RefreshToken.
type RefreshTokenResponse struct {
	RefreshToken string
	ExpiresAt    time.Time
}
