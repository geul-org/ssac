//ff:type feature=pkg-file type=model
//ff:what 파일 업로드 요청 모델
package file

type UploadRequest struct {
	Key  string
	Body string
}
