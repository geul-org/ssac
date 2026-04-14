//ff:type feature=pkg-crypto type=model
//ff:what TOTP OTP 생성 요청 모델
package crypto

type GenerateOTPRequest struct {
	Issuer      string
	AccountName string
}
