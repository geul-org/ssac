//ff:type feature=pkg-crypto type=model
//ff:what AES-256-GCM 암호화 응답 모델
package crypto

type EncryptResponse struct {
	Ciphertext string // base64 인코딩
}
