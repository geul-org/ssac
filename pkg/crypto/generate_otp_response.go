//ff:type feature=pkg-crypto type=model
//ff:what TOTP OTP 생성 응답 모델
package crypto

type GenerateOTPResponse struct {
	Secret string
	URL    string // otpauth:// URL (QR 코드용)
}
