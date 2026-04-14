//ff:type feature=pkg-storage type=model
//ff:what 서명된 URL 생성 요청 모델
package storage

type PresignURLRequest struct {
	Bucket    string
	Key       string
	ExpiresIn int // 초 단위 (기본 3600)
	Endpoint  string
	Region    string
}
