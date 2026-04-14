//ff:type feature=pkg-crypto type=model
//ff:what AES-256-GCM 암호화 요청 모델
package crypto

type EncryptRequest struct {
	Plaintext string
	Key       string // 32바이트 hex
}
