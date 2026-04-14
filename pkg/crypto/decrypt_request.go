//ff:type feature=pkg-crypto type=model
//ff:what AES-256-GCM 복호화 요청 모델
package crypto

type DecryptRequest struct {
	Ciphertext string // base64 인코딩
	Key        string // 32바이트 hex
}
