//ff:type feature=pkg-auth type=model
//ff:what 비밀번호 해싱 응답 모델
package auth

type HashPasswordResponse struct {
	HashedPassword string
}
