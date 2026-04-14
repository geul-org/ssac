//ff:func feature=pkg-file type=loader control=sequence
//ff:what S3 파일 저장소 생성 — 클라이언트와 버킷 이름으로 인스턴스 반환
package file

import "github.com/aws/aws-sdk-go-v2/service/s3"

// NewS3File creates a FileModel backed by AWS S3.
func NewS3File(client *s3.Client, bucket string) FileModel {
	return &s3File{client: client, bucket: bucket}
}
