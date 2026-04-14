//ff:type feature=pkg-auth type=model
//ff:what 비밀번호 리셋 토큰 생성 응답 모델
package auth

type GenerateResetTokenResponse struct {
	Token string
}
