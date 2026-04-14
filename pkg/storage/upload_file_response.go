//ff:type feature=pkg-storage type=model
//ff:what S3 파일 업로드 응답 모델
package storage

type UploadFileResponse struct {
	URL string
}
