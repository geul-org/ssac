//ff:type feature=pkg-auth type=model topic=auth-jwt
//ff:what JWT 검증 응답 모델
package auth

import "github.com/golang-jwt/jwt/v5"

// VerifyTokenResponse returns the verified claim set as jwt.MapClaims.
// Callers re-marshal via encoding/json into their own typed Claim struct
// when they need field-level access.
type VerifyTokenResponse struct {
	Claims jwt.MapClaims
}
