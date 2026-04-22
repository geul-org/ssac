//ff:type feature=pkg-auth type=model topic=auth-jwt
//ff:what JWT 리프레시 토큰 발급 요청 모델
package auth

// RefreshTokenRequest holds the inputs for RefreshToken.
//
// Claims is opaque passthrough; see IssueTokenRequest for semantics.
type RefreshTokenRequest struct {
	Claims any
}
