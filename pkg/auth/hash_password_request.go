//ff:type feature=pkg-auth type=model
//ff:what 비밀번호 해싱 요청 모델
package auth

type HashPasswordRequest struct {
	Password string
}
