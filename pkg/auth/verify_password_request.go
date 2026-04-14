//ff:type feature=pkg-auth type=model
//ff:what 비밀번호 검증 요청 모델
package auth

type VerifyPasswordRequest struct {
	PasswordHash string
	Password     string
}
