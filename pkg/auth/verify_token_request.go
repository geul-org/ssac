//ff:type feature=pkg-auth type=model topic=auth-jwt
//ff:what JWT 검증 요청 모델
package auth

// VerifyTokenRequest holds the inputs for VerifyToken.
type VerifyTokenRequest struct {
	Token string
}
