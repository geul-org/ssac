//ff:type feature=pkg-storage type=model
//ff:what S3 파일 업로드 요청 모델
package storage

type UploadFileRequest struct {
	Bucket      string
	Key         string
	Data        []byte
	ContentType string
	Endpoint    string // MinIO 등 커스텀 엔드포인트 (빈 문자열이면 AWS 기본)
	Region      string
}
