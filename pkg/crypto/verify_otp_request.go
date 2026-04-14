//ff:type feature=pkg-crypto type=model
//ff:what TOTP OTP 검증 요청 모델
package crypto

type VerifyOTPRequest struct {
	Code   string
	Secret string
}
