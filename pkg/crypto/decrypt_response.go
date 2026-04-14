//ff:type feature=pkg-crypto type=model
//ff:what AES-256-GCM 복호화 응답 모델
package crypto

type DecryptResponse struct {
	Plaintext string
}
