//ff:type feature=pkg-auth type=model topic=auth-jwt
//ff:what JWT 액세스 토큰 발급 응답 모델
package auth

import "time"

// IssueTokenResponse is the result of IssueToken.
//
// AccessToken is the HS256-signed JWT string. ExpiresAt mirrors the `exp`
// claim embedded in the token so callers can track expiry without re-parsing.
type IssueTokenResponse struct {
	AccessToken string
	ExpiresAt   time.Time
}
