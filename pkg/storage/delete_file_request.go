//ff:type feature=pkg-storage type=model
//ff:what S3 파일 삭제 요청 모델
package storage

type DeleteFileRequest struct {
	Bucket   string
	Key      string
	Endpoint string
	Region   string
}
