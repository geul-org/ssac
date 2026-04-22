//ff:type feature=pkg-auth type=model topic=auth-jwt
//ff:what JWT 액세스 토큰 발급 요청 모델
package auth

// IssueTokenRequest holds the inputs for IssueToken.
//
// Claims is passed through as the JWT payload. Callers should use a struct
// with JSON tags whose keys match the claim names expected by their
// application (e.g. `json:"user_id"`, `json:"email"`). The concrete type is
// opaque to this package — Claims is marshaled via encoding/json and then
// remapped into jwt.MapClaims before signing. A nil Claims is normalized to
// an empty object so the issued token carries only the injected `exp` claim.
type IssueTokenRequest struct {
	Claims any
}
